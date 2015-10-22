package model

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/deis/helm/helm/log"
	"github.com/deis/helm/helm/manifest"
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

	sortManifests(c, ms)

	return c, nil
}

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
			chart.Pods = append(chart.Pods, vo.(*v1.Pod))
		case "ReplicationController":
			chart.ReplicationControllers = append(chart.ReplicationControllers, vo.(*v1.ReplicationController))
		case "Service":
			chart.Services = append(chart.Services, vo.(*v1.Service))
		case "Secret":
			chart.Secrets = append(chart.Secrets, vo.(*v1.Secret))
		case "PersistentVolume":
			chart.PersistentVolumes = append(chart.PersistentVolumes, vo.(*v1.PersistentVolume))
		case "Namespace":
			chart.Namespaces = append(chart.Namespaces, vo.(*v1.Namespace))
		}
	}
}
