package action

import (
	"path/filepath"

	"github.com/helm/helm/log"
	"github.com/helm/helm/search"
)

// Search looks for packages with 'term' in their name.
func Search(term, homedir string, regexp bool) {
	cfg := mustConfig(homedir)
	cdir := filepath.Join(homedir, CachePath)

	i := search.NewIndex(cfg, cdir)
	res, err := i.Search(term, 5, regexp)
	if err != nil {
		log.Die("Failed to search: %s", err)
	}

	if len(res) == 0 {
		log.Err("No results found. Try using '--regexp'.")
		return
	}

	search.SortScore(res)

	for _, r := range res {
		c, _ := i.Chart(r.Name)
		log.Msg("%s - %s", r.Name, c.Description)
	}
}
