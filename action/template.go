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
	"github.com/helm/helm/log"
	"gopkg.in/yaml.v2"
)

// Template renders a template to an output file.
func Template(out, in, data string) {
	var dest io.Writer
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

	var vals interface{}
	if data != "" {
		var err error
		vals, err = openValues(data)
		if err != nil {
			log.Die("Error opening value file: %s", err)
		}
	}

	tpl, err := ioutil.ReadFile(in)
	if err != nil {
		log.Die("Failed to read template file: %s", err)
	}

	if err := renderTemplate(dest, string(tpl), vals); err != nil {
		log.Die("Failed: %s", err)
	}
}

// openValues opens a values file and tries to parse it with the right parser.
//
// It returns an interface{} containing data, if found. Any error opening or
// parsing the file will be passed back.
func openValues(filename string) (interface{}, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
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

	if err := t.ExecuteTemplate(out, "helmTpl", vals); err != nil {
		return err
	}
	return nil
}
