package main

import (
	"io"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/codegangsta/cli"
	"github.com/helm/helm/log"
	"gopkg.in/yaml.v2"
)

var version = "dev"

func main() {
	app := cli.NewApp()
	app.Name = "helm-tpl"
	app.Version = version
	app.Usage = `Run a file through Go's text template engine`
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "file,f",
			Usage: "The name of the file to write to.",
		},
		cli.StringFlag{
			Name:  "template, t",
			Usage: "The template file.",
		},
		cli.StringFlag{
			Name:  "values,d",
			Usage: "The path to a YAML file with values.",
		},
	}

	app.Action = generate
	app.Run(os.Args)
}

func generate(c *cli.Context) {
	dname := c.String("file")
	tname := c.String("template")
	vname := c.String("values")

	if dname == "" {
		log.Die("No destination is set. Use --file.")
	}
	if tname == "" {
		log.Die("No template is set. Use --template.")
	}

	var vals interface{}
	if vname == "" {
		log.Warn("No template values set. Using built-ins.")
	} else {
		var err error
		vals, err = readValues(vname)
		if err != nil {
			log.Die("Error reading %s", err)
		}
	}

	out, err := os.Create(dname)
	if err != nil {
		log.Die("Error creating %s: %s", dname, err)
	}

	if err := renderTemplate(tname, out, vals); err != nil {
		// Die will prevent a defer from running.
		out.Close()
		log.Die("Error rendering %s: %s", tname, err)
	}
	out.Close()
}

func readValues(filename string) (interface{}, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		log.Warn("Empty template values file %s", filename)
		return map[string]string{}, nil
	}

	dest := map[string]interface{}{}
	if err := yaml.Unmarshal(data, dest); err != nil {
		return nil, err
	}
	log.Info("Values: %v", dest)
	return dest, nil
}

func renderTemplate(tfile string, out io.Writer, vals interface{}) error {
	t, err := template.New("helmTpl").Funcs(sprig.TxtFuncMap()).ParseFiles(tfile)
	if err != nil {
		return err
	}

	if err := t.ExecuteTemplate(out, tfile, vals); err != nil {
		return err
	}
	return nil
}
