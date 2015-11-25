package codec

import (
	"io/ioutil"
	"path"
	"strings"
	"testing"
)

func TestManifest(t *testing.T) {
	d, err := ioutil.ReadFile(path.Join(testdata, "pod.yaml"))
	if err != nil {
		t.Error(err)
	}

	m, err := YAML.Decode(d).One()
	if err != nil {
		t.Errorf("Failed parse: %s", err)
	}

	pod, err := m.Pod()
	if err != nil {
		t.Errorf("Failed to decode into pod: %s", err)
	}

	if pod.Name != "cassandra" {
		t.Errorf("Expected name 'cassandra', got %q", pod.Name)
	}
}

type kindFunc func(m *Object) error

func assertKind(t *testing.T, file string, kf kindFunc) {
	d, err := ioutil.ReadFile(path.Join(testdata, file))
	if err != nil {
		t.Error(err)
	}

	m, err := YAML.Decode(d).One()
	if err != nil {
		t.Errorf("Failed parse: %s", err)
	}

	if err = kf(m); err != nil {
		t.Errorf("Failed to decode %s into its kind: %s", file, err)
	}
}

func TestKnownKinds(t *testing.T) {
	kinds := map[string]kindFunc{
		"pod.yaml":                     func(m *Object) error { _, err := m.Pod(); return err },
		"rc.yaml":                      func(m *Object) error { _, err := m.RC(); return err },
		"daemonset.yaml":               func(m *Object) error { _, err := m.DaemonSet(); return err },
		"horizontalpodautoscaler.yaml": func(m *Object) error { _, err := m.HorizontalPodAutoscaler(); return err },
		"ingress.yaml":                 func(m *Object) error { _, err := m.Ingress(); return err },
		"job.yaml":                     func(m *Object) error { _, err := m.Job(); return err },
		"serviceaccount.yaml":          func(m *Object) error { _, err := m.ServiceAccount(); return err },
		"service.yaml":                 func(m *Object) error { _, err := m.Service(); return err },
		"namespace.yaml":               func(m *Object) error { _, err := m.Namespace(); return err },
	}

	for uptown, funk := range kinds {
		assertKind(t, uptown, funk /*gonna give it to ya*/)
		// Don't believe me? Just watch.
	}
}

func TestServiceAccount(t *testing.T) {
	d, err := ioutil.ReadFile(path.Join(testdata, "serviceaccount.yaml"))
	if err != nil {
		t.Error(err)
	}

	m, err := YAML.Decode(d).One()
	if err != nil {
		t.Errorf("Failed parse: %s", err)
	}

	_, err = m.ServiceAccount()
	if err != nil {
		t.Errorf("Failed to decode into pod: %s", err)
	}
}

func TestObjectYAML(t *testing.T) {
	d, err := ioutil.ReadFile(path.Join(testdata, "serviceaccount.yaml"))
	if err != nil {
		t.Error(err)
	}
	m, err := YAML.Decode(d).One()
	if err != nil {
		t.Errorf("Failed parse: %s", err)
	}

	if out, err := m.YAML(); err != nil {
		t.Errorf("Failed to write YAML: %s", err)
	} else if len(out) == 0 {
		t.Error("YAML len is 0")
	}
}
func TestObjectJSON(t *testing.T) {
	d, err := ioutil.ReadFile(path.Join(testdata, "serviceaccount.yaml"))
	if err != nil {
		t.Error(err)
	}
	m, err := YAML.Decode(d).One()
	if err != nil {
		t.Errorf("Failed parse: %s", err)
	}

	if out, err := m.JSON(); err != nil {
		t.Errorf("Failed to write JSON: %s", err)
	} else if len(out) == 0 {
		t.Error("JSON len is 0")
	}
}

func TestAddLabels(t *testing.T) {
	d, err := ioutil.ReadFile(path.Join(testdata, "pod.yaml"))
	if err != nil {
		t.Error(err)
	}

	m, err := YAML.Decode(d).One()
	if err != nil {
		t.Errorf("Failed parse: %s", err)
	}

	labels := map[string]string{
		"foo":   "bar",
		"drink": "slurm",
	}

	if err := m.AddLabels(labels); err != nil {
		t.Errorf("Failed to add labels: %s", err)
	}

	t.Log(string(m.data))
	if !strings.Contains(string(m.data), "drink: slurm") {
		t.Errorf("Could not find 'drink:slurm' in \n%s", string(m.data))
	}

	_, err = m.Pod()
	if err != nil {
		t.Errorf("Failed to decode into pod: %s", err)
	}
}

func TestTemplateManifest(t *testing.T) {
	d, err := ioutil.ReadFile(path.Join(testdata, "template.yaml"))
	if err != nil {
		t.Error(err)
	}

	m, err := YAML.Decode(d).One()
	if err != nil {
		t.Errorf("Failed parse: %s", err)
	}

	template, err := m.Template()
	if err != nil {
		t.Errorf("Failed to decode into template: %s", err)
	}

	if template.Name != "console" {
		t.Errorf("Expected name 'console', got %q", template.Name)
	}
	if len(template.Objects) != 4 {
		t.Errorf("Expected 4 template objects, got %d", len(template.Objects))
	} else {
		for _, json := range template.Objects {
			rcm, err := YAML.Decode(json.RawJSON).One()
			if err != nil {
				t.Errorf("Failed parse RC: %s", err)
			}
			ref, err := rcm.Ref()
			if err != nil {
				t.Errorf("Failed parsing Ref of template object: %s", err)
			} else if ref.Kind == "ReplicationController" {
				rc, err := rcm.RC()
				if err != nil {
					t.Errorf("Failed unmarshalling of RC: %s", err)
				}
				if rc.Kind != "ReplicationController" {
					t.Errorf("Expected kind 'ReplicationController' for template object 4, got %s", rc.Kind)
				}
				if rc.Name != "fabric8" {
					t.Errorf("Expected name 'fabric8' for template object 4, got %s", rc.Name)
				}
			}
		}
	}
}
