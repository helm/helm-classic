package action

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/BurntSushi/toml"
	"github.com/Masterminds/sprig"
	"github.com/helm/helm-classic/log"
	"gopkg.in/yaml.v2"
)

var err error

//GenerateTemplate evaluates a template and writes it to an io.Writer
func GenerateTemplate(out io.Writer, in io.Reader, vals interface{}) {
	tpl, err := ioutil.ReadAll(in)
	if err != nil {
		log.Die("Failed to read template file: %s", err)
	}

	if err := renderTemplate(out, string(tpl), vals); err != nil {
		log.Die("Template rendering failed: %s", err)
	}
}

// Template renders a template to an output file.
func Template(out, in, data string, force bool) error {
	var dest io.Writer
	_, err = os.Stat(out)
	if !(force || os.Getenv("HELM_FORCE_FLAG") == "true") && err == nil {
		return fmt.Errorf("File %s already exists. To overwrite it, please re-run this command with the --force/-f flag.", out)
	}
	if out != "" {
		f, err := os.Create(out)
		if err != nil {
			log.Die("Failed to open %s: %s", out, err)
		}
		defer func() {
			if err := f.Close(); err != nil {
				log.Err("Error closing file: %s", err)
			}
		}()
		dest = f
	} else {
		dest = log.Stdout
	}

	inReader, err := os.Open(in)
	if err != nil {
		log.Die("Failed to open template file: %s", err)
	}

	var vals interface{}
	if data != "" {
		var err error
		vals, err = openValues(data)
		if err != nil {
			log.Die("Error opening value file: %s", err)
		}
	}

	GenerateTemplate(dest, inReader, vals)
	return nil
}

// openValues opens a values file and tries to parse it with the right parser.
//
// It returns an interface{} containing data, if found. Any error opening or
// parsing the file will be passed back.
func openValues(filename string) (interface{}, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		// We generate a warning here, but do not require that a values
		// file exists.
		log.Warn("Skipped file %s: %s", filename, err)
		return map[string]interface{}{}, nil
	}

	ext := filepath.Ext(filename)
	var um func(p []byte, v interface{}) error
	switch ext {
	case ".yaml", ".yml":
		um = yaml.Unmarshal
	case ".toml":
		um = toml.Unmarshal
	case ".json":
		um = json.Unmarshal
	default:
		return nil, fmt.Errorf("Unsupported file type: %s", ext)
	}

	var res interface{}
	err = um(data, &res)
	return res, err
}

// renderTemplate renders a template and values into an output stream.
//
// tpl should be a string template.
func renderTemplate(out io.Writer, tpl string, vals interface{}) error {
	t, err := template.New("helmTpl").Funcs(sprig.TxtFuncMap()).Parse(tpl)
	if err != nil {
		return err
	}

	log.Debug("Vals: %#v", vals)

	if err = t.ExecuteTemplate(out, "helmTpl", vals); err != nil {
		return err
	}
	return nil
}
