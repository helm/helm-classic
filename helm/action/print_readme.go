package action

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/deis/helm/helm/log"
)

// PrintREADME prints the README file (if it exists) to the console.
func PrintREADME(chart, home string) {
	p := filepath.Join(home, WorkspaceChartPath, chart, "README.*")
	files, err := filepath.Glob(p)
	if err != nil || len(files) == 0 {
		// No README. Skip.
		log.Debug("No readme in %s", p)
		return
	}

	f, err := os.Open(files[0])
	if err != nil {
		log.Warn("Could not read README: %s", err)
		return
	}
	log.Msg(strings.Repeat("=", 40))
	io.Copy(log.Stdout, f)
	log.Msg(strings.Repeat("=", 40))
	f.Close()

}
