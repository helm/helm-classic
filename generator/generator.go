package generator

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/helm/helm-classic/log"
)

// GeneratorKeyword is used to generate new charts
const GeneratorKeyword = "helm:generate "

// Walk walks a chart directory and executes generators as it finds them.
//
// Returns the number of generators executed.
//
// Walking will error out whenever a generator cannot be completely executed.
// This includes cases such as not finding the generator referenced, and
// cases where the generator itself exits with a non-zero exit code.
func Walk(dir string, exclude []string, force bool) (int, error) {

	excludes := make(map[string]bool, len(exclude))
	for i := 0; i < len(exclude); i++ {
		excludes[filepath.Join(dir, exclude[i])] = true
	}

	count := 0
	err := filepath.Walk(dir, func(path string, fi os.FileInfo, err error) error {

		// dive-bomb if we hit an error.
		if err != nil {
			return err
		}

		// Exclude anything explicitly excluded.
		if excludes[path] == true {
			if fi.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip directory entries. If the directory prefix is . or _, skip the
		// contents of the directory as well.
		if fi.IsDir() {
			return skip(path)
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		line, err := readGenerator(f)
		if err != nil {
			return err
		}
		if line == "" {
			return nil
		}
		// Run the generator.
		os.Setenv("HELM_GENERATE_COMMAND", line)
		os.Setenv("HELM_GENERATE_FILE", path)
		os.Setenv("HELM_GENERATE_DIR", dir)
		line = os.ExpandEnv(line)
		os.Setenv("HELM_GENERATE_COMMAND_EXPANDED", line)
		log.Debug("File: %s, Command: %s", path, line)
		count++

		// Execute the command in the file's directory to make relative
		// paths usable.
		origin, err := os.Getwd()
		if err != nil {
			log.Warn("Could not get PWD: %s", err)
		} else if err := os.Chdir(dir); err != nil {
			log.Warn("Could not change directory to %s: %s", dir, err)
		} else {
			origin = dir
			defer func() {
				if e := os.Chdir(origin); e != nil {
					log.Warn("Could not return to %s: %s", origin, e)
				}
			}()
		}
		err = execute(line, force)
		if err != nil {
			return fmt.Errorf("failed to execute %s (%s): %s", line, path, err)
		}
		return nil
	})

	return count, err
}

func execute(command string, force bool) error {
	args := strings.Fields(command)
	if len(args) == 0 {
		return errors.New("empty command")
	}
	name := args[0]
	if args[0] == "helm" && (args[1] == "template" || args[1] == "tpl") && force {
		args = append([]string{args[1], "-f"}, args[2:]...)
	} else {
		args = args[1:]
	}

	// Templates will often include a helm:generate header such as the following:
	// #helm:generate helm template ...
	// For backwards compatibility with older charts, this needs to be supported, even though the
	// name of Helm Classic binary has changed to "helmc". To accommodate this AND especially
	// ensure that a Helm Classic generator doesn't inadvertently invoke the new kubernetes/helm,
	// we detect generator commands starting with "helm" and replace them with "helmc".
	if name == "helm" {
		name = "helmc"
	}

	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

// skip indicates whether the directory's contents should be skipped.
//
// error is nil unless the directory passes the skip test, in which acse it is
// filepath.SkipDir
func skip(path string) error {
	base := filepath.Base(path)
	if base[0] == '.' || base[0] == '_' {
		return filepath.SkipDir
	}
	return nil
}

// Read the generator from a file.
//
// An error indicates that something went wrong.
//
// An empty string indicates that there was no generator.
//
// A string is to be treated as the value of the generator, without the
// `helm:generate` prefix.
func readGenerator(file *os.File) (string, error) {

	f := bufio.NewReader(file)

	// Look for leading `//`, `#`, or `/*`
	var b []byte
	var err error
	if b, err = f.Peek(3); err != nil {
		return "", nil
	}

	offset := 0
	suffix := ""
	if b[0] == '#' {
		offset++
		if b[1] == ' ' {
			offset++
		}
	} else if b[0] == '/' && (b[1] == '/' || b[1] == '*') {
		offset += 2
		if b[2] == ' ' {
			offset++
		}
		if b[1] == '*' {
			suffix = "*/"
		}
	} else {
		return "", nil
	}

	if _, err := f.Discard(offset); err != nil {
		return "", err
	}

	// If we get here, we have a comment header. Next, check if it's a helm:generate header.
	if b, err = f.Peek(len(GeneratorKeyword)); err != nil {
		return "", nil
	}

	slug := string(b)
	if slug != GeneratorKeyword {
		return "", nil
	}
	if _, err := f.Discard(len(GeneratorKeyword)); err != nil {
		return "", err
	}

	// At this point, we know that we have a helm:generate header. Read to EOL.
	line, err := f.ReadString('\n')
	if err != nil {
		return "", err
	}

	line = strings.TrimSpace(line)
	if len(suffix) > 0 {
		line = strings.TrimSpace(strings.TrimSuffix(line, suffix))
	}
	return line, err
}
