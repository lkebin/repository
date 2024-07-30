package parser

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	distinct             = "Distinct"
	countByTemplate      = regexp.MustCompile("^" + countPattern + "(\\p{Lu}.*?)??By")
	existsByTemplate     = regexp.MustCompile("^" + existsPattern + "(\\p{Lu}.*?)??By")
	deleteByTemplate     = regexp.MustCompile("^" + deletePattern + "(\\p{Lu}.*?)??By")
	limitingQueryPattern = "(First|Top)(\\d*)?"
	limitedQueryTemplate = regexp.MustCompile("^(" + queryPattern + ")(" + distinct + ")?" + limitingQueryPattern + "(\\p{Lu}.*?)??By")
)

type Subject struct {
	isDistinct bool
	isCount    bool
	isExists   bool
	isDelete   bool
	isLimiting bool
	maxResults int
}

func NewSubject(subject string) *Subject {
	s := &Subject{}

	s.isDistinct = strings.Contains(subject, distinct)
	s.isCount = matches(subject, countByTemplate)
	s.isExists = matches(subject, existsByTemplate)
	s.isDelete = matches(subject, deleteByTemplate)
	s.maxResults = returnMaxResultsIfFirstKSubject(subject)
	s.isLimiting = s.maxResults > 0

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
