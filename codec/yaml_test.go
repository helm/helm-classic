package codec

import (
	"bytes"
	"io/ioutil"
	"path"
	"testing"
)

const testdata = "../testdata"

func TestYamlDecoderOne(t *testing.T) {
	d, err := ioutil.ReadFile(path.Join(testdata, "pod.yaml"))
	if err != nil {
		t.Error(err)
	}

	m, err := YAML.Decode(d).One()
	if err != nil {
		t.Error(err)
	}

	ref, err := m.Ref()
	if err != nil {
		t.Errorf("Could not get reference: %s", err)
	}
	if ref.Kind != "Pod" {
		t.Errorf("Expected a pod, got a %s", ref.Kind)
	}
	if ref.APIVersion != "v1" {
		t.Errorf("Expected v1, got %s", ref.APIVersion)
	}
}

func TestYamlDecoderAll(t *testing.T) {
	d, err := ioutil.ReadFile(path.Join(testdata, "three-pods-and-three-services.yaml"))
	if err != nil {
		t.Error(err)
	}

	ms, err := YAML.Decode(d).All()
	if err != nil {
		t.Error(err)
	}

	if len(ms) != 6 {
		t.Errorf("Expected 6 parts, got %d", len(ms))
	}

	for i := 0; i < 3; i++ {
		ref, err := ms[i*2].Ref()
		if err != nil {
			t.Errorf("Expected a reference for pod[%d]: %s", i, err)
		}

		if ref.Kind != "Pod" {
			t.Errorf("Expected Pod, got %s", ref.Kind)
		}

		ref, err = ms[i*2+1].Ref()
		if err != nil {
			t.Errorf("Expected a reference for service[%d]: %s", i, err)
		}

		if ref.Kind != "Service" {
			t.Errorf("Expected Service, got %s", ref.Kind)
		}
	}
}

func TestYamlEncoderAll(t *testing.T) {
	f1 := map[string]string{"one": "hello"}
	f2 := map[string]string{"two": "world"}

	var b bytes.Buffer
	if err := YAML.Encode(&b).All(f1, f2); err != nil {
		t.Errorf("Failed to encode: %s", err)
	}

	// This is a little fragile, since whitespace in YAML is not defined.
	expect := "one: hello\n\n---\ntwo: world\n"
	actual := b.String()
	if actual != expect {
		t.Errorf("Expected [%s]\nGot [%s]", expect, actual)
	}
}
