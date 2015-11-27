package action

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/helm/helm/chart"
	"github.com/helm/helm/codec"
	"github.com/helm/helm/dependency"
	"github.com/helm/helm/log"
	"github.com/helm/helm/manifest"
	"github.com/helm/helm/parameters"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/v1"
	"k8s.io/kubernetes/pkg/api/errors"
	"k8s.io/kubernetes/pkg/runtime"
	utilerr "k8s.io/kubernetes/pkg/util/errors"

	"github.com/openshift/origin/pkg/template"
	"github.com/openshift/origin/pkg/template/generator"
	tapi "github.com/openshift/origin/pkg/template/api"
	templatevalidation "github.com/openshift/origin/pkg/template/api/validation"
	utilerr "k8s.io/kubernetes/pkg/util/errors"

)

// InstallOrder defines the order in which manifests should be installed, by Kind.
//
// Anything not on the list will be installed after the last listed item, in
// an indeterminate order.
var InstallOrder = []string{"Namespace", "Secret", "PersistentVolume", "ServiceAccount", "Service", "Pod", "ReplicationController", "DaemonSet", "Ingress", "Job"}

// UninstallOrder defines the order in which manifests are uninstalled.
//
// Unknown manifest types (those not explicitly referenced in this list) will
// be uninstalled before any of these, since we know that none of the core
// types depend on non-core types.
var UninstallOrder = []string{"Service", "Pod", "ReplicationController", "DaemonSet", "Secret", "PersistentVolume", "ServiceAccount", "Ingress", "Job", "Namespace"}

// Install loads a chart into Kubernetes.
//
// If the chart is not found in the workspace, it is fetched and then installed.
//
// During install, manifests are sent to Kubernetes in the ordered specified by InstallOrder.
func Install(chartName, home, namespace string, mode string, force bool, dryRun bool, valueFlag string, paramFolder string, printImportFolders bool, writeGeneratedKeys bool, generateSecretsData bool) {
	secretFlags := &secretSettings{PrintImportFolders: printImportFolders, WriteGeneratedKeys: writeGeneratedKeys, GenerateSecretsData: generateSecretsData}
	ochart := chartName
	r := mustConfig(home).Repos
	table, chartName := r.RepoChart(chartName)

	if !chartFetched(chartName, home) {
		log.Info("No chart named %q in your workspace. Fetching now.", ochart)
		fetch(chartName, chartName, home, table)
	}

	cd := filepath.Join(home, WorkspaceChartPath, chartName)
	cTemplates, err := chart.Load(cd)
	if err != nil {
		log.Die("Failed to load chart: %s", err)
	}
	c, err := processTemplates(cTemplates, valueFlag, paramFolder)

	// Give user the option to bale if dependencies are not satisfied.
	nope, err := dependency.Resolve(c.Chartfile, filepath.Join(home, WorkspaceChartPath))
	if err != nil {
		log.Warn("Failed to check dependencies: %s", err)
		if !force {
			log.Die("Re-run with --force to install anyway.")
		}
	} else if len(nope) > 0 {
		log.Warn("Unsatisfied dependencies:")
		for _, d := range nope {
			log.Msg("\t%s %s", d.Name, d.Version)
		}
		if !force {
			log.Die("Stopping install. Re-run with --force to install anyway.")
		}
	}

	CheckKubePrereqs()

	msg := "Running `kubectl create -f` ..."
	if dryRun {
		msg = "Performing a dry run of `kubectl create -f` ..."
	}
	log.Info(msg)
	if err := uploadManifests(c, namespace, mode, dryRun, secretFlags); err != nil {
		log.Die("Failed to upload manifests: %s", err)
	}
	log.Info("Done")

	PrintREADME(chartName, home)
}

func isSamePath(src, dst string) (bool, error) {
	a, err := filepath.Abs(dst)
	if err != nil {
		return false, err
	}
	b, err := filepath.Abs(src)
	if err != nil {
		return false, err
	}
	return a == b, nil
}

// Processes any OpenShift templates inside the chart and
// removes a new chart without any OpenShift templates
func processTemplates(c *chart.Chart, valueFlag string, paramFolder string) (*chart.Chart, error) {
	if len(c.Templates) == 0 {
		return c, nil
	}
	nc := &chart.Chart{
		Chartfile: c.Chartfile,
	}

	ms := []*manifest.Manifest{}
	for _, t := range c.Templates {
		log.Debug("Replacing templates in %s with %d objects", t.Name, len(t.Objects))
		tpl := &tapi.Template{}
		scheme := runtime.NewScheme()
		scheme.Convert(t, tpl)

		if len(t.Parameters) != len(tpl.Parameters) {
			for _, p := range t.Parameters {
				p2 := tapi.Parameter{Name: p.Name, Value: p.Value, Generate: p.Generate, From: p.From, DisplayName: p.DisplayName, Description: p.Description}
				tpl.Parameters = append(tpl.Parameters, p2)
			}
		}
		if len(t.Parameters) != len(tpl.Parameters) {
			log.Die("Failed to convert template %s with %d parameters as has %d runtime parameters", tpl.Name, len(t.Parameters), len(tpl.Parameters))
		}

		chartName := c.Chartfile.Name
		customParams, err := parameters.LoadChartParameters(paramFolder, chartName)
		if err != nil {
			log.Die("Failed to load previous chart parameter values %s\n", err)
		}
		customized := false
		if len(valueFlag) > 0 {
			values := strings.Split(valueFlag, ",")
			for _, keypair := range values {
				p := strings.SplitN(keypair, "=", 2)
				if len(p) != 2 {
					log.Die("invalid parameter assignment in %q: %q\n", t.Name, keypair)
					continue
				}
				customized = true;
				customParams.Values[p[0]] = p[1]
			}
		}

		for key, value := range customParams.Values {
			if v := template.GetParameterByName(tpl, key); v != nil {
				v.Value = value
				v.Generate = ""
				template.AddParameter(tpl, *v)
			} else {
				log.Die("unknown parameter name %q\n", key)
			}
		}

		kubeCodec := runtime.CodecFor(api.Scheme, t.APIVersion)
		for _, o := range t.Objects {
			o2, err := kubeCodec.Decode(o.RawJSON)
			if err != nil {
				log.Die("Failed to unmarshal JSON with error: %s", err)
			}
			tpl.Objects = append(tpl.Objects, o2)
		}

		if errs := templatevalidation.ValidateProcessedTemplate(tpl); len(errs) > 0 {
			err := errors.NewInvalid("template", tpl.Name, errs)
			log.Die("Failed to validate template: %s", err)
			return nil, err
		}

		generators := map[string]generator.Generator{
			"expression": generator.NewExpressionValueGenerator(rand.New(rand.NewSource(time.Now().UnixNano()))),
		}
		processor := template.NewProcessor(generators)
		if errs := processor.Process(tpl); len(errs) > 0 {
			log.Info("Errors in processor")
			log.Die(utilerr.NewAggregate(errs).Error())
			return nil, errors.NewInvalid("template", tpl.Name, errs)
		}

		for _, tobject := range tpl.Objects {
			buffer := new(bytes.Buffer)
			if err := kubeCodec.EncodeToStream(tobject, buffer); err != nil {
				log.Die("Failed to encode codec: %s", err)
			}
			json := buffer.String()
			doc, err := codec.YAML.Decode(buffer.Bytes()).One()
			if err != nil {
				log.Die("Failed parse RC: %s", err)
			}
			ref, err := doc.Ref()
			if err != nil {
				log.Die("Failed parsing Ref of template object: %s", err)
			} else {
				m := &manifest.Manifest{Version: ref.APIVersion, Kind: ref.Kind, VersionedObject: doc, Source: json}
				ms = append(ms, m)
			}
		}

		if customized {
			err := parameters.SaveChartParameters(paramFolder, chartName, customParams)
			if err != nil {
				log.Die("Failed to save chart parameters: %s", err)
			}
		}
	}
	chart.SortManifests(nc, ms)
	return nc, nil
}

// uploadManifests sends manifests to Kubectl in a particular order.
func uploadManifests(c *chart.Chart, namespace string, mode string, dryRun bool, secretFlags *secretSettings) error {

	// Install known kinds in a predictable order.
	for _, k := range InstallOrder {
		for _, m := range c.Kind[k] {
			o := m.VersionedObject
			o.AddAnnotations(map[string]string{
				chart.AnnFile:         m.Source,
				chart.AnnChartVersion: c.Chartfile.Version,
				chart.AnnChartDesc:    c.Chartfile.Description,
				chart.AnnChartName:    c.Chartfile.Name,
			})
			var data []byte
			var err error
			if data, err = o.JSON(); err != nil {
				return err
			}
			log.Info("Data: %s", data)
			if err := kubectlCreate(data, namespace, dryRun); err != nil {
				return err
			}
		}
	}

	// Install unknown kinds afterward. Order here is not predictable.
	for _, k := range c.UnknownKinds(InstallOrder) {
		for _, o := range c.Kind[k] {
			o.VersionedObject.AddAnnotations(map[string]string{chart.AnnFile: o.Source})
			if err := marshalAndCreate(o.VersionedObject, namespace, dryRun); err != nil {
				return err
			}
		}
	}

	return nil
}

func marshalAndCreate(o interface{}, ns string, dry bool) error {
	var b bytes.Buffer
	if err := codec.JSON.Encode(&b).One(o); err != nil {
		return err
	}
	return kubectlCreate(b.Bytes(), ns, mode, dry)
}

// Check by chart directory name whether a chart is fetched into the workspace.
//
// This does NOT check the Chart.yaml file.
func chartFetched(chartName, home string) bool {
	p := filepath.Join(home, WorkspaceChartPath, chartName, "Chart.yaml")
	log.Debug("Looking for %q", p)
	if fi, err := os.Stat(p); err != nil || fi.IsDir() {
		log.Debug("No chart: %s", err)
		return false
	}
	return true
}

// kubectlCreate calls `kubectl create` and sends the data via Stdin.
//
// If dryRun is set to true, then we just output the command that was
// going to be run to os.Stdout and return nil.
func kubectlCreate(data []byte, ns string, dryRun bool) error {
	a := []string{"create", "-f", "-"}

	if ns != "" {
		a = append([]string{"--namespace=" + ns}, a...)
	}

	if dryRun {
		cmd := "kubectl"
		for _, arg := range a {
			cmd = fmt.Sprintf("%s %s", cmd, arg)
		}
		cmd = fmt.Sprintf("%s < %s", cmd, data)
		log.Info(cmd)
		return nil
	}

	c := exec.Command("kubectl", a...)
	in, err := c.StdinPipe()
	if err != nil {
		return err
	}

	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if err := c.Start(); err != nil {
		return err
	}

	log.Debug("File: %s", string(data))
	in.Write(data)
	in.Close()

	return c.Wait()
}

func kubeCtlGetResourceVersion(ns string, kind string, name string) (string, error) {
	b, err := kubeCtlGetJson(ns, kind, name)
	if err != nil {
		return "", err
	}
	kubeCodec := runtime.CodecFor(api.Scheme, defaultAPIVersion)
	o, err := kubeCodec.Decode(b)
	if err != nil {
		return "", err
	}
	objectMeta, err := api.ObjectMetaFor(o)
	if err != nil {
		return "", err
	}
	return objectMeta.ResourceVersion, nil
}

func kubeCtlGetJson(ns string, kind string, name string) ([]byte, error) {
	cmd := "kubectl"
	a := []string{}
	if ns != "" {
		a = append([]string{"--namespace=" + ns}, a...)
	}
	a = append([]string{"get", strings.ToLower(kind), name, "-ojson"}, a...)
	log.Info("About to invoke command %s %s", cmd, strings.Join(a, " "))
	return exec.Command(cmd, a...).Output()
}
