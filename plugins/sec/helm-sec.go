package main

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/aokoli/goutils"
	"github.com/codegangsta/cli"
	"github.com/helm/helm-classic/log"
	"github.com/pborman/uuid"
	"golang.org/x/crypto/nacl/box"
)

var version = "dev"

func main() {
	app := cli.NewApp()
	app.Name = "helm-sec"
	app.Version = version
	app.Usage = `Manage secrets.`
	app.ArgsUsage = "[SecretName] [SecretValue]"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "file,f",
			Usage: "The name of the file to write to.",
		},
		cli.StringFlag{
			Name:  "input, i",
			Usage: "A file whose contents will be used as the secret value.",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "The name, as placed in the metadata name field.",
		},
		cli.BoolFlag{
			Name:  "password,p",
			Usage: "Generate a password to be used in the value. Use any printable ASCII character.",
		},
		cli.BoolFlag{
			Name:  "alphanum",
			Usage: "Generate an alphanumeric password to be used in the value.",
		},
		cli.BoolFlag{
			Name:  "alpha",
			Usage: "Generate a letters-only password to be used in the value.",
		},
		cli.BoolFlag{
			Name:  "numeric",
			Usage: "Generate an numeric password to be used in the value.",
		},
		cli.BoolFlag{
			Name:  "uuid",
			Usage: "Generate a UUID v4 random identifier.",
		},
		cli.BoolFlag{
			Name:  "box",
			Usage: "Generate a public and private key for NaCl box.",
		},
		cli.BoolFlag{
			Name:  "quiet, q",
			Usage: "Hush optional information. e.g. don't print generated passwords.",
		},
		cli.IntFlag{
			Name:  "len,l",
			Value: 16,
			Usage: "The length of the generated field.",
		},
	}

	app.Action = sec
	app.Run(os.Args)
}

func sec(c *cli.Context) {
	a := c.Args()

	f, err := fileOrStdout(c)
	if err != nil {
		log.Die("Could not open file for writing: %s", err)
	}
	defer f.Close()

	if len(a) < 1 {
		log.Die("At least one argument (name) is required.")
	}
	name := a[0]
	mdname := c.String("name")
	if mdname == "" {
		mdname = name
	}

	// These take multiple values
	if c.Bool("box") {
		vals, err := boxKeys(c, name)
		if err != nil {
			log.Die("Could not generate keys: %s", err)
		}
		if err := renderSecret(f, c, mdname, vals); err != nil {
			log.Die("Failed to generate keys: %s", err)
		}
		return
	}

	// Default is to handle just one value.
	v, err := resolveValue(a, c)
	if err != nil {
		log.Die("Failed to get the secret: %s", err)
	}

	if err := renderSecret(f, c, mdname, map[string]interface{}{name: v}); err != nil {
		log.Die("Failed to generate secret: %s", err)
	}
}

func boxKeys(c *cli.Context, name string) (map[string]interface{}, error) {
	pub, priv, err := box.GenerateKey(rand.Reader)
	if err != nil {
		return map[string]interface{}{}, err
	}
	return map[string]interface{}{
		name:          base64.StdEncoding.EncodeToString(priv[:]),
		name + ".pub": base64.StdEncoding.EncodeToString(pub[:]),
	}, nil
}

func resolveValue(a []string, c *cli.Context) (string, error) {
	val := ""
	if len(a) > 1 {
		val = a[1]
	}

	if c.Bool("password") {
		val = genPW(c, goutils.RandomAscii)
	} else if c.Bool("alpha") {
		val = genPW(c, goutils.RandomAlphabetic)
	} else if c.Bool("alphanum") {
		val = genPW(c, goutils.RandomAlphaNumeric)
	} else if c.Bool("numeric") {
		val = genPW(c, goutils.RandomNumeric)
	} else if c.Bool("uuid") {
		val = uuid.New()
		printOrQuiet(c, "UUID: %s", val)
	} else if c.Bool("box") {

	} else if i := c.String("input"); len(i) > 0 {
		data, err := ioutil.ReadFile(i)
		if err != nil {
			return val, err
		}
		val = string(data)
	}

	val = base64.StdEncoding.EncodeToString([]byte(val))

	return val, nil
}

func printOrQuiet(c *cli.Context, msg string, v ...interface{}) {
	if !c.Bool("quiet") {
		log.Info(msg, v...)
	}
}

type lenPWFunc func(l int) (string, error)

func genPW(c *cli.Context, f lenPWFunc) string {
	d, _ := f(c.Int("len"))
	printOrQuiet(c, "Password: %q", d)
	return d
}

// fileOrStdout returns a file if one is specified on the commnad line, or Stdout.
//
// If a file is specified, but cannot be opened for writing, this generates an error
// and returns Stdout.
func fileOrStdout(c *cli.Context) (*os.File, error) {
	f := c.String("file")
	if strings.TrimSpace(f) == "" {
		return os.Stdout, nil
	}
	file, err := os.Create(f)
	if err != nil {
		return os.Stdout, err
	}
	return file, nil
}

const tpl = `kind: Secret
apiVersion: v1
metadata:
  name: {{.MDName}}
  data:{{range $k, $v := .Values}}
    {{$k}}: {{$v}}{{end}}
`

func renderSecret(f io.Writer, c *cli.Context, name string, values map[string]interface{}) error {
	vals := map[string]interface{}{
		"Values": values,
		"MDName": name,
	}

	t := template.Must(template.New("secret").Parse(tpl))
	return t.Execute(f, vals)
}
