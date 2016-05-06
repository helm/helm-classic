// Package search provides search features for Helm Classic.
package search

import (
	"errors"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/helm/helm-classic/chart"
	"github.com/helm/helm-classic/config"
	"github.com/helm/helm-classic/log"
)

// Result is a search result.
//
// Score indicates how close it is to match. The higher the score, the longer
// the distance.
type Result struct {
	Name  string
	Score int
}

// Index is a searchable index of chart information.
type Index struct {
	lines  map[string]string
	charts map[string]*chart.Chartfile
}

const sep = "\v"

// NewIndex creats a new Index.
//
// NewIndex indexes all of the chart tables configured in the config.yaml file.
// For that reason, it may cause substantial overhead on a large set of repos.
func NewIndex(cfg *config.Configfile, cachedir string) *Index {
	lines := map[string]string{}
	charts := map[string]*chart.Chartfile{}
	for _, table := range cfg.Repos.Tables {
		def := cfg.Repos.Default == table.Name

		base := filepath.Join(cachedir, table.Name, "*/")
		dirs, err := filepath.Glob(base)
		if err != nil {
			log.Err("Failed to read table %s: %s", table.Name, err)
		}

		for _, dir := range dirs {
			bname := filepath.Base(dir)
			c, err := chart.LoadChartfile(filepath.Join(dir, "Chart.yaml"))
			if err != nil {
				// This is not a chart. Skip it.
				continue
			}
			name := table.Name + "/" + c.Name
			if def {
				name = c.Name
			}
			line := c.Name + sep + table.Name + "/" + bname + sep + c.Description + sep + c.Details
			lines[name] = strings.ToLower(line)
			charts[name] = c
		}
	}
	return &Index{lines: lines, charts: charts}
}

// Search searches an index for the given term.
//
// Threshold indicates the maximum score a term may have before being marked
// irrelevant. (Low score means higher relevance. Golf, not bowling.)
//
// If regexp is true, the term is treated as a regular expression. Otherwise,
// term is treated as a literal string.
func (i *Index) Search(term string, threshold int, regexp bool) ([]*Result, error) {
	if regexp == true {
		return i.SearchRegexp(term, threshold)
	}
	return i.SearchLiteral(term, threshold), nil
}

// calcScore calculates a score for a match.
func (i *Index) calcScore(index int, matchline string) int {

	// This is currently tied to the fact that sep is a single char.
	splits := []int{}
	s := rune(sep[0])
	for i, ch := range matchline {
		if ch == s {
			splits = append(splits, i)
		}
	}

	for i, pos := range splits {
		if index > pos {
			continue
		}
		return i
	}
	return len(splits)
}

// SearchLiteral does a literal string search (no regexp).
func (i *Index) SearchLiteral(term string, threshold int) []*Result {
	term = strings.ToLower(term)
	buf := []*Result{}
	for k, v := range i.lines {
		res := strings.Index(v, term)
		if score := i.calcScore(res, v); res != -1 && score < threshold {
			buf = append(buf, &Result{Name: k, Score: score})
		}
	}
	return buf
}

// SearchRegexp searches using a regular expression.
func (i *Index) SearchRegexp(re string, threshold int) ([]*Result, error) {
	matcher, err := regexp.Compile(re)
	if err != nil {
		return []*Result{}, err
	}
	buf := []*Result{}
	for k, v := range i.lines {
		ind := matcher.FindStringIndex(v)
		if len(ind) == 0 {
			continue
		}
		if score := i.calcScore(ind[0], v); ind[0] >= 0 && score < threshold {
			buf = append(buf, &Result{Name: k, Score: score})
		}
	}
	return buf, nil
}

// ErrNoChart indicates that a chart is not in the chart cache.
var ErrNoChart = errors.New("no such chart")

// Chart gets the *Chartfile for a chart that was found during search.
//
// This is a convenience method for retrieving a cached chart that was located
// during search indexing.
func (i *Index) Chart(name string) (*chart.Chartfile, error) {
	c, ok := i.charts[name]
	if !ok {
		return nil, errors.New("no such chart")
	}
	return c, nil
}

// SortScore does an in-place sort of the results.
//
// Lowest scores are highest on the list. Matching scores are subsorted alphabetically.
func SortScore(r []*Result) {
	sort.Sort(scoreSorter(r))
}

// scoreSorter sorts results by score, and subsorts by alpha Name.
type scoreSorter []*Result

// Len returns the length of this scoreSorter.
func (s scoreSorter) Len() int { return len(s) }

// Swap performs an in-place swap.
func (s scoreSorter) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// Less compares a to b, and returns true if a is less than b.
func (s scoreSorter) Less(a, b int) bool {
	first := s[a]
	second := s[b]

	if first.Score > second.Score {
		return false
	}
	if first.Score < second.Score {
		return true
	}
	return first.Name < second.Name
}
