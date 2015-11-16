package action

import (
	"bytes"
	"os"
	"testing"

	pretty "github.com/deis/pkg/prettyprint"
	"github.com/helm/helm/log"
)

func TestRemove(t *testing.T) {
	tmpHome := createTmpHome()
	fakeUpdate(tmpHome)

	Fetch("kitchensink", "", tmpHome)

	var output bytes.Buffer
	log.Stderr = &output

	Remove("kitchensink", tmpHome)

	expected := pretty.Colorize("{{.Green}}--->{{.Default}} ") + "All clear! You have successfully removed kitchensink from your workspace.\n"

	expect(t, output.String(), expected)

	// reset log
	log.Stdout = os.Stdout
}
