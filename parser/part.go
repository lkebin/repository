package parser

import "regexp"

type IgnoreCaseType int

const (
	Never IgnoreCaseType = iota
	Always
	WhenPossible
)

var (
	ignoreCase = regexp.MustCompile("Ignor(ing|e)Case")
)

type Part struct {
	property   string
	typ        PartType
	ignoreCase IgnoreCaseType
}

func NewPart(source string, isAlwaysIgnoreCase bool) *Part {
	part := &Part{}

	partToUse := part.detectAndSetIgnoreCase(source)

	if isAlwaysIgnoreCase && part.ignoreCase != Always {
		part.ignoreCase = WhenPossible
	}

	part.typ = NewPartTypeFromProperty(partToUse)
	part.property = part.typ.extractProperty(partToUse)

	return part
}

func (p *Part) detectAndSetIgnoreCase(part string) string {
	result := part
	match := ignoreCase.FindStringIndex(part)

	if match != nil {
		p.ignoreCase = Always
		result = part[:match[0]] + part[match[1]:]
	}

	return result
}
