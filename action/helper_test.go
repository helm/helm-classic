package action

import "github.com/helm/helm-classic/log"

func init() {
	// Turn on debug output, convert os.Exit(1) to panic()
	log.IsDebugging = true
}
