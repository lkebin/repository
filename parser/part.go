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
	Property   string
	Type       PartType
	IgnoreCase IgnoreCaseType
}

func NewPart(source string, isAlwaysIgnoreCase bool) *Part {
	part := &Part{}

	partToUse := part.detectAndSetIgnoreCase(source)

	if isAlwaysIgnoreCase && part.IgnoreCase != Always {
		part.IgnoreCase = WhenPossible
	}

	part.Type = NewPartTypeFromProperty(partToUse)
	part.Property = part.Type.extractProperty(partToUse)

	return part
}

func (p *Part) detectAndSetIgnoreCase(part string) string {
	result := part
	match := ignoreCase.FindStringIndex(part)

	if match != nil {
		p.IgnoreCase = Always
		result = part[:match[0]] + part[match[1]:]
	}

	return result
}
