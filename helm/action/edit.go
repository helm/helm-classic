package action

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
)

import (
	"github.com/deis/helm/helm/log"
)

// delimeterRegexp captures relative path and file content
var delimeterRegexp = regexp.MustCompile(`--- # (.+)\n([\s\S]*?)--- # end\n`)

// Edit charts using the shell-defined $EDITOR
//
// - chartName being edited
// - homeDir is the helm home directory for the user
func Edit(chartName, homeDir string) {

	chartDir := path.Join(homeDir, "workspace", "charts", chartName)

	// enumerate chart files
	files, err := listChart(chartDir)
	if err != nil {
		log.Die("could not list chart: %v", err)
	}

	// join chart with YAML delimeters
	contents, err := joinChart(chartDir, files)
	if err != nil {
		log.Die("could not join chart data: %v", err)
	}

	// write chart to temporary file
	f, err := ioutil.TempFile(os.TempDir(), "helm-edit")
	if err != nil {
		log.Die("could not open tempfile: %v", err)
	}
	f.Write(contents.Bytes())
	f.Close()

	// NOTE: removing the tempfile causes issues with editors
	// that fork, so we let the OS remove them later

	openEditor(f.Name())
	saveChart(chartDir, f.Name())

}

// listChart enumerates all of the relevant files in a chart
func listChart(chartDir string) ([]string, error) {

	var files []string

	metadataFile := path.Join(chartDir, "Chart.yaml")
	manifestDir := path.Join(chartDir, "manifests")

	// check for existence of important files and directories
	for _, path := range []string{chartDir, metadataFile, manifestDir} {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return nil, err
		}
	}

	// add metadata file to front of list
	files = append(files, metadataFile)

	// add manifest files
	walker := func(fname string, fi os.FileInfo, e error) error {
		if e != nil {
			log.Warn("Encounter error walking %q: %s", fname, e)
			return nil
		}

		if filepath.Ext(fname) == ".yaml" {
			files = append(files, fname)
		}

		return nil
	}
	filepath.Walk(manifestDir, walker)

	return files, nil
}

// joinChart reads chart files and joins them with YAML delimiters
func joinChart(chartDir string, files []string) (bytes.Buffer, error) {

	var output bytes.Buffer

	for _, f := range files {
		contents, err := ioutil.ReadFile(f)
		if err != nil {
			return output, err
		}

		rf, err := filepath.Rel(chartDir, f)
		if err != nil {
			log.Warn("Could not find relative path: %s", err)
			return output, err
		}

		delimiter := fmt.Sprintf("--- # %s\n", rf)

		output.WriteString(delimiter)
		output.Write(contents)
		output.WriteString("--- # end\n")

	}

	return output, nil
}

// openEditor opens the given filename in an interactive editor
func openEditor(filename string) {
	var cmd *exec.Cmd

	editor := os.ExpandEnv("$EDITOR")
	if editor == "" {
		log.Die("must set shell $EDITOR")
	}

	args := []string{filename}
	cmd = exec.Command(editor, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// saveChart reads a delimited chart and write out its parts
// to the workspace directory
func saveChart(chartDir string, filename string) error {

	// read the serialized chart file
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	chartData := make(map[string][]byte)

	// use a regular expression to read file paths and content
	match := delimeterRegexp.FindAllSubmatch(contents, -1)
	for _, m := range match {
		chartData[string(m[1])] = m[2]
	}

	// save edited chart data to the workspace
	for k, v := range chartData {
		fp := path.Join(chartDir, k)
		if err := ioutil.WriteFile(fp, v, 0644); err != nil {
			log.Die("could not write chart file", err)
		}
	}
	return nil

}
