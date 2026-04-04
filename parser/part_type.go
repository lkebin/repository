package parser

import "strings"

var (
	Between                = PartType{"Between", 2, []string{"IsBetween", "Between"}}
	IsNotNull              = PartType{"IsNotNull", 0, []string{"IsNotNull", "NotNull"}}
	IsNull                 = PartType{"IsNull", 0, []string{"IsNull", "Null"}}
	LessThan               = PartType{"LessThan", 1, []string{"IsLessThan", "LessThan"}}
	LessThanEqual          = PartType{"LessThanEqual", 1, []string{"IsLessThanEqual", "LessThanEqual"}}
	GreaterThan            = PartType{"GreaterThan", 1, []string{"IsGreaterThan", "GreaterThan"}}
	GreaterThanEqual       = PartType{"GreaterThanEqual", 1, []string{"IsGreaterThanEqual", "GreaterThanEqual"}}
	Before                 = PartType{"Before", 1, []string{"IsBefore", "Before"}}
	After                  = PartType{"After", 1, []string{"IsAfter", "After"}}
	NotLike                = PartType{"NotLike", 1, []string{"IsNotLike", "NotLike"}}
	Like                   = PartType{"Like", 1, []string{"IsLike", "Like"}}
	StartingWith           = PartType{"StartingWith", 1, []string{"IsStartingWith", "StartingWith"}}
	EndingWith             = PartType{"EndingWith", 1, []string{"IsEndingWith", "EndingWith"}}
	IsNotEmpty             = PartType{"IsNotEmpty", 0, []string{"IsNotEmpty", "NotEmpty"}}
	IsEmpty                = PartType{"IsEmpty", 0, []string{"IsEmpty", "Empty"}}
	NotContaining          = PartType{"NotContaining", 1, []string{"IsNotContaining", "NotContaining"}}
	Containing             = PartType{"Containing", 1, []string{"IsContaining", "Containing"}}
	NotIn                  = PartType{"NotIn", 1, []string{"IsNotIn", "NotIn"}}
	In                     = PartType{"In", 1, []string{"IsIn", "In"}}
	Near                   = PartType{"Near", 1, []string{"IsNear", "Near"}}
	WithIn                 = PartType{"WithIn", 1, []string{"IsWithIn", "WithIn"}}
	Regex                  = PartType{"Regex", 1, []string{"IsRegex", "Regex"}}
	Exists                 = PartType{"Exists", 0, []string{"Exists"}}
	True                   = PartType{"True", 0, []string{"IsTrue", "True"}}
	False                  = PartType{"False", 0, []string{"IsFalse", "False"}}
	NegatingSimpleProperty = PartType{"NegatingSimpleProperty", 1, []string{"IsNot", "Not"}}
	SimpleProperty         = PartType{"SimpleProperty", 1, []string{"Is", "Equals"}}
)

var (
	// All contains all supported operators for method name parsing.
	// Near, WithIn, Regex, and Exists are intentionally excluded because they
	// require database-specific features or subquery support not available in
	// the generator.
	All = []PartType{
		IsNotNull,
		IsNull,
		Between,
		LessThan,
		LessThanEqual,
		GreaterThan,
		GreaterThanEqual,
		Before,
		After,
		NotLike,
		Like,
		StartingWith,
		EndingWith,
		IsNotEmpty,
		IsEmpty,
		NotContaining,
		Containing,
		NotIn,
		In,
		True,
		False,
		NegatingSimpleProperty,
		SimpleProperty,
	}

	AllKeywords = func() []string {
		var keywords []string
		for _, partType := range All {
			keywords = append(keywords, partType.Keywords...)
		}
		return keywords
	}()
)

type PartType struct {
	Name              string
	NumberOfArguments int
	Keywords          []string
}

func NewPartTypeFromProperty(property string) PartType {
	for _, partType := range All {
		if partType.supports(property) {
			return partType
		}
	}
	return SimpleProperty
}

func (p *PartType) supports(rawProperty string) bool {
	for _, keyword := range p.Keywords {
		if strings.HasSuffix(rawProperty, keyword) {
			return true
		}
	}
	return false
}

func (p *PartType) extractProperty(part string) string {
	for _, keyword := range p.Keywords {
		if strings.HasSuffix(part, keyword) {
			return strings.TrimSuffix(part, keyword)
		}
	}
	return part
}
