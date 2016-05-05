package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Masterminds/vcs"
	"github.com/helm/helm-classic/log"
	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/yaml.v2"
)

// ErrNotFound indicates no local repository could be found.
var ErrNotFound = errors.New("No local repository")

// Configfile is the top-level conifguration object for Helm Classic.
type Configfile struct {
	// filename may contain a reference back to the file that was read into
	// this object.
	filename string

	// Repos points to the repository configuration
	Repos     *Repos     `yaml:"repos"`
	Workspace *Workspace `yaml:"workspace"`
}

// Repos describes a collection of repository (table) mappings.
type Repos struct {
	// Dir points to the directory where the Git repositories are stored.
	// Currently we do not expose this as a changeable option, but we may
	// choose to in the future.
	Dir string `yaml:"-"`
	// Default is the local name of the default repository.
	Default string `yaml:"default"`
	// Tables is a list of table items.
	Tables []*Table `yaml:"tables"`
}

// Workspace describes a workspace location and configuration.
type Workspace struct {
	// Dir indicates where the workspace is.
	// Currently, this is not exposed via configuration.
	Dir string `yaml:"-"`
}

// Table describes a single table entry.
type Table struct {
	// Name is the local name of the repository.
	Name string `yaml:"name"`
	// Repo is the remote Git URL to the repository.
	Repo string `yaml:"repo"`
}

// Load loads a configuration by filename.
func Load(filename string) (*Configfile, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	abs, err := filepath.Abs(filename)
	if err != nil {
		return nil, err
	}

	cfg, err := Parse(b)
	if err != nil {
		return cfg, err
	}
	cfg.filename = abs
	if cfg.Repos.Dir == "" {
		cfg.Repos.Dir = filepath.Join(filepath.Dir(abs), "cache")
	}

	if cfg.Workspace == nil {
		cfg.Workspace = &Workspace{}
	}

	if cfg.Workspace.Dir == "" {
		cfg.Workspace.Dir = filepath.Join(filepath.Dir(abs), "workspace")
	}

	return cfg, nil
}

// Parse parses a byte slice into a *Configfile.
func Parse(data []byte) (*Configfile, error) {
	r := &Configfile{
		filename: "config.yaml",
	}

	if err := yaml.Unmarshal(data, r); err != nil {
		return r, err
	}
	return r, nil
}

// Save writes the Configfile as YAML into the named file.
func (c *Configfile) Save(filename string) error {
	if filename == "" {
		filename = c.filename
	}
	b, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, b, 0755)
}

// RepoChart takes a fully qualified name and returns a repo name and a chart name.
//
// If no repo name is present in the fully qualified name, the default repo
// is returned.
func (r *Repos) RepoChart(name string) (string, string) {
	res := strings.SplitN(name, "/", 2)
	if len(res) == 1 {
		return r.Default, name
	}
	return res[0], res[1]
}

// Add adds the named remote and then fetches it.
func (r *Repos) Add(name, repo string) error {
	for _, r := range r.Tables {
		if r.Name == name {
			return fmt.Errorf("Remote %s already exists, and is pointed to %s", name, r.Repo)
		}
	}

	nt := &Table{
		Name: name,
		Repo: repo,
	}

	r.Tables = append(r.Tables, nt)
	if err := r.Update(name); err != nil {
		return err
	}

	return nil
}

// Exists checks if a repo exists by name
func (r *Repos) Exists(name string) bool {
	for _, r := range r.Tables {
		if r.Name == name {
			return true
		}
	}
	return false
}

// Update performs an update of the local copy.
//
// This does a Git fast-forward pull from the remote repo.
func (r *Repos) Update(name string) error {
	for _, t := range r.Tables {
		if t.Name == name {
			rpath := filepath.Join(r.Dir, name)
			g, err := ensureRepo(t.Repo, rpath)
			if err != nil {
				return err
			}
			return g.Update()
		}
	}
	return ErrNotFound
}

func ensureRepo(repo, dir string) (*vcs.GitRepo, error) {
	if fi, err := os.Stat(dir); err != nil {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, err
		}
	} else if !fi.IsDir() {
		return nil, fmt.Errorf("File %s exists, but is not a directory.", dir)
	}
	git, err := vcs.NewGitRepo(repo, dir)
	if err != nil {
		return nil, err
	}

	git.Logger = log.New()

	if !git.CheckLocal() {
		if err := git.Get(); err != nil {
			return git, err
		}
	}
	return git, nil
}

// UpdateAll does a git fast-forward pull from each remote repo.
func (r *Repos) UpdateAll() error {
	for _, table := range r.Tables {
		log.Info("Checking repository %s", table.Name)
		rpath := filepath.Join(r.Dir, table.Name)
		g, err := ensureRepo(table.Repo, rpath)
		if err != nil {
			return err
		}

		if g.IsDirty() {
			return fmt.Errorf("Repository '%s' is dirty.  Commit changes before updating", table.Name)
		}

		initialVersion, err := g.Version()
		if err != nil {
			return fmt.Errorf("Could not get current sha of repository '%s'.", table.Name)
		}

		if err := g.Update(); err != nil {
			return err
		}
		diff, err := repoChartDiff(rpath, initialVersion)
		if err != nil {
			return err
		}
		printSummary(diff)
	}
	return nil
}

func repoChartDiff(rpath, initialVersion string) (string, error) {
	// build git diff-tree command
	cmd := exec.Command("git", "-C", rpath, "diff-tree", "--name-status", fmt.Sprintf("%s..HEAD", initialVersion))

	log.Debug("git diff cmd: %s", cmd.Args)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	// cleanup any trailing whitespace
	return strings.TrimSpace(string(out)), nil
}

type repoSummary map[string][]string

func (s repoSummary) add(status, chart string) {
	ignores := []string{"README.md", ".gitignore", ".travis.yml", "Makefile", "_test"}
	for _, i := range ignores {
		if chart == i {
			return
		}
	}
	s[status] = append(s[status], chart)
}

// printSummary prints a diff of charts after upate
func printSummary(diff string) {
	if len(diff) == 0 {
		log.Msg("Already up-to-date.")
		return
	}

	s := make(repoSummary)

	// parse git diff-tree
	for _, line := range strings.Split(diff, "\n") {
		kv := strings.Split(line, "\t")
		st, chart := kv[0], kv[1]

		s.add(st, chart)
	}

	// width of columns
	colwidth := 29

	// get console width
	maxwidth := 118
	if w, _, err := terminal.GetSize(int(os.Stdout.Fd())); err == nil {
		maxwidth = w
	}

	// print results
	for st, charts := range s {
		switch st {
		case "A":
			log.Msg("Added %d charts", len(charts))
		case "D":
			log.Msg("Sent %d charts to the depths", len(charts))
		case "M":
			log.Msg("Updated %d charts", len(charts))
		}

		line := ""
		for _, ch := range charts {
			// if adding this column passes the max
			// print and reset to zero
			if len(line)+colwidth > maxwidth {
				log.Msg(line)
				line = ""
			}
			// append to line with padding
			line = fmt.Sprintf("%s%-29s", line, ch)
		}
		log.Msg(line)
	}
}

// Delete removes a local copy of a remote.
//
// This destroys the on-disk cache and removes the entry from the YAML file.
func (r *Repos) Delete(name string) error {
	res := []*Table{}

	counter := 0
	for _, t := range r.Tables {
		if t.Name == name {
			counter++
			continue
		}
		res = append(res, t)
	}
	if counter == 0 {
		return fmt.Errorf("No repository named %s", name)
	}

	r.Tables = res
	return r.deleteRepo(name)
}

func (r *Repos) deleteRepo(name string) error {
	rpath := filepath.Join(r.Dir, name)
	if fi, err := os.Stat(rpath); err != nil || !fi.IsDir() {
		log.Info("Deleted nothing. No repo named %s", name)
		return nil
	}

	log.Debug("Deleting %s", rpath)
	return os.RemoveAll(rpath)
}
