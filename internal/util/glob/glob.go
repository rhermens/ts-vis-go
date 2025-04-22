package glob

import "slices"

import g "github.com/gobwas/glob"

type Glob g.Glob;

func AnyMatches(filters []Glob, input ...string) bool {
	for _, filter := range filters {
		if slices.ContainsFunc(input, filter.Match) {
			return true
		}
	}

	return false
}

func CompileFilters(input []string) []Glob {
	var filters []Glob
	for _, filter := range input {
		filters = append(filters, g.MustCompile(filter))
	}

	return filters
}

