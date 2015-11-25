package codec

import (
	"io/ioutil"
	"path"
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
