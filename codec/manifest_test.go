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
