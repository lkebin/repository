package parser

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	distinct             = "Distinct"
	countByTemplate      = regexp.MustCompile(`^` + countPattern + `(\p{Lu}.*?)??By`)
	existsByTemplate     = regexp.MustCompile(`^` + existsPattern + `(\p{Lu}.*?)??By`)
	deleteByTemplate     = regexp.MustCompile(`^` + deletePattern + `(\p{Lu}.*?)??By`)
	limitingQueryPattern = `(First|Top)(\d*)?`
	limitedQueryTemplate = regexp.MustCompile(`^(` + queryPattern + `)(` + distinct + `)?` + limitingQueryPattern + `(\p{Lu}.*?)??By`)
)

type Subject struct {
	IsDistinct bool
	IsCount    bool
	IsExists   bool
	IsDelete   bool
	IsLimiting bool
	MaxResults int
}

func NewSubject(subject string) *Subject {
	s := &Subject{}

	s.IsDistinct = strings.Contains(subject, distinct)
	s.IsCount = matches(subject, countByTemplate)
	s.IsExists = matches(subject, existsByTemplate)
	s.IsDelete = matches(subject, deleteByTemplate)
	s.MaxResults = returnMaxResultsIfFirstKSubject(subject)
	s.IsLimiting = s.MaxResults > 0

	return s
}

func matches(s string, r *regexp.Regexp) bool {
	return r.Match([]byte(s))
}

func returnMaxResultsIfFirstKSubject(subject string) int {
	matches := limitedQueryTemplate.FindStringSubmatch(subject)
	if len(matches) == 0 {
		return 0
	}

	if matches[4] == "" {
		return 1
	}

	mr, err := strconv.Atoi(matches[4])
	if err != nil {
		panic(err)
	}
	return mr
}
