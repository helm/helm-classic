package chart

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/deis/helm/log"
	"github.com/deis/helm/manifest"
	//"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/v1"
)

// Chart represents a complete chart.
//
// A chart consists of the following parts:
//
// 	- Chart.yaml: In code, we refer to this as the Chartfile
// 	- manifests/*.yaml: The Kubernetes manifests
//
// On the Chart object, the manifests are sorted by type into a handful of
// recognized Kubernetes API v1 objects.
//
// TODO: Investigate treating these as unversioned.
type Chart struct {
	Chartfile              *Chartfile
	Pods                   []*v1.Pod
	ReplicationControllers []*v1.ReplicationController
	Services               []*v1.Service
	Namespaces             []*v1.Namespace
	Secrets                []*v1.Secret
	PersistentVolumes      []*v1.PersistentVolume

	Manifests []*manifest.Manifest
}

// Load loads an entire chart.
//
// This includes the Chart.yaml (*Chartfile) and all of the manifests.
//
// If you are just reading the Chart.yaml file, it is substantially more
// performant to use LoadChartfile.
func Load(chart string) (*Chart, error) {
	if fi, err := os.Stat(chart); err != nil {
		return nil, err
	} else if !fi.IsDir() {
		return nil, fmt.Errorf("Chart %s is not a directory.", chart)
	}

	cf, err := LoadChartfile(filepath.Join(chart, "Chart.yaml"))
	if err != nil {
		return nil, err
	}

	c := &Chart{
		Chartfile: cf,
	}

	ms, err := manifest.ParseDir(chart)
	if err != nil {
		return c, err
	}

	c.Manifests = ms
	sortManifests(c, ms)

	return c, nil
}

// Save writes an entire chart to disk.
//
// It will overwrite any files that it finds in the way.
//
// This writes a `Chart.yaml`, and then writes manifests into a `manifests`
// directory, creating the directory if it needs to.
/*
func (c *Chart) Save(dir string) error {
	if fi, err := os.Stat(dir); err != nil {
		return fmt.Errorf("Could not save Chart.yaml: %s", err)
	} else if !fi.IsDir() {
		return fmt.Errorf("Not a directory: %s", dir)
	}

	if err := c.Chartfile.Save(filepath.Join(dir, "Chart.yaml")); err != nil {
		return err
	}

	mdir := filepath.Join(dir, "manifests")
}
*/

const OriginFile = "HelmOriginFile"

// sortManifests sorts manifests into their respective categories, adding to the Chart.
func sortManifests(chart *Chart, manifests []*manifest.Manifest) {
	for _, m := range manifests {
		vo := m.VersionedObject

		if m.Version != "v1" {
			log.Warn("Unsupported version %q", m.Version)
			continue
		}

		switch m.Kind {
		default:
			log.Warn("No support for kind %s. Ignoring.", m.Kind)
		case "Pod":
			o := vo.(*v1.Pod)
			o.Annotations = setOriginFile(o.Annotations, m.Source)
			chart.Pods = append(chart.Pods, o)
		case "ReplicationController":
			o := vo.(*v1.ReplicationController)
			o.Annotations = setOriginFile(o.Annotations, m.Source)
			chart.ReplicationControllers = append(chart.ReplicationControllers, o)
		case "Service":
			o := vo.(*v1.Service)
			o.Annotations = setOriginFile(o.Annotations, m.Source)
			chart.Services = append(chart.Services, o)
		case "Secret":
			o := vo.(*v1.Secret)
			o.Annotations = setOriginFile(o.Annotations, m.Source)
			chart.Secrets = append(chart.Secrets, o)
		case "PersistentVolume":
			o := vo.(*v1.PersistentVolume)
			o.Annotations = setOriginFile(o.Annotations, m.Source)
			chart.PersistentVolumes = append(chart.PersistentVolumes, o)
		case "Namespace":
			o := vo.(*v1.Namespace)
			o.Annotations = setOriginFile(o.Annotations, m.Source)
			chart.Namespaces = append(chart.Namespaces, o)
		}
	}
}

func setOriginFile(ann map[string]string, source string) map[string]string {
	if len(ann) == 0 {
		return map[string]string{OriginFile: source}
	}
	ann[OriginFile] = source
	return ann
}
