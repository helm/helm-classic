package repo

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/vcs"
	"github.com/deis/helm/log"
	"gopkg.in/yaml.v2"
)

const DefaultRepofile = `default: charts
tables:
  - name: charts
    repo: https://github.com/deis/charts
`

var NotFound = errors.New("No local repository")

// Repofile describes a configuration file for remote repositories.
type Repofile struct {
	filename string

	// Dir points to the directory where the Git repositories are stored.
	Dir string `yaml:"-"`
	// Default is the local name of the default repository.
	Default string `yaml:"default"`
	// Tables is a list of table items.
	Tables []*Table `yaml:tables`
}

// Table describes a single table entry.
type Table struct {
	// Name is the local name of the repository.
	Name string `yaml:"name"`
	// Repo is the remote Git URL to the repository.
	Repo string `yaml:"repo"`
}

// LoadRepofile loads a Repofile from a filename.
//
// The directory where the git repos live is determined by analyzing the path
// to the filename. It can be overridden by explicitly setting Repofile.Dir.
func LoadRepofile(filename string) (*Repofile, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	abs, err := filepath.Abs(filename)
	if err != nil {
		return nil, err
	}

	rf, err := ParseRepofile(b)
	if err != nil {
		return rf, err
	}
	rf.filename = abs
	rf.Dir = filepath.Dir(abs)

	return rf, nil
}

// ParseRepofile parses a byte slice into a *Repofile.
//
// Note that the Repofile.Dir is set to "."
func ParseRepofile(data []byte) (*Repofile, error) {
	r := &Repofile{
		filename: "Remote.yaml",
		Dir:      ".",
	}
	return r, yaml.Unmarshal(data, r)
}

// Save writes the Repofile as YAML into the named file.
func (r *Repofile) Save(filename string) error {
	b, err := yaml.Marshal(r)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, b, 0755)
}

// RepoChart takes a fully qualified name and returns a repo name and a chart name.
//
// If no repo name is present in the fully qualified name, the default repo
// is returned.
func (r *Repofile) RepoChart(name string) (string, string) {
	res := strings.SplitN(name, "/", 2)
	if len(res) == 1 {
		return r.Default, name
	}
	return res[0], res[1]
}

// Add adds the named remote and then fetches it.
func (r *Repofile) Add(name, repo string) error {
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

	if err := r.Save(r.filename); err != nil {
		r.deleteRepo(name)
		return err
	}

	return nil
}

// Update performs an update of the local copy.
//
// This does a Git fast-forward pull from the remote repo.
func (r *Repofile) Update(name string) error {
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
	return NotFound
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

// Update all remotes.
//
// This does a Git fast-forward pull from each remote repo.
func (r *Repofile) UpdateAll() error {
	for _, table := range r.Tables {
		log.Info("Updating %s", table.Name)
		rpath := filepath.Join(r.Dir, table.Name)
		g, err := ensureRepo(table.Repo, rpath)
		if err != nil {
			return err
		}
		if err := g.Update(); err != nil {
			return err
		}
	}
	return nil
}

// Delete removes a local copy of a remote.
//
// This destroys the on-disk cache and removes the entry from the YAML file.
func (r *Repofile) Delete(name string) error {
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
	if err := r.Save(r.filename); err != nil {
		return err
	}

	return r.deleteRepo(name)
}

func (r *Repofile) deleteRepo(name string) error {
	rpath := filepath.Join(r.Dir, name)
	if fi, err := os.Stat(rpath); err != nil || !fi.IsDir() {
		log.Info("Deleted nothing. No repo named %s", name)
		return nil
	}

	log.Debug("Deleting %s", rpath)
	return os.RemoveAll(rpath)
}
