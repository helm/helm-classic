package kubectl

import (
	"testing"
)

func TestGet(t *testing.T) {
	Client = TestRunner{
		out: []byte("running the get command"),
	}

	expects := "running the get command"
	out, _ := Client.Get([]byte{}, "", false)
	if string(out) != expects {
		t.Errorf("%s != %s", string(out), expects)
	}
}
