package parser

import "strings"

var (
	Between                = PartType{2, []string{"IsBetween", "Between"}}
	IsNotNull              = PartType{0, []string{"IsNotNull", "NotNull"}}
	IsNull                 = PartType{0, []string{"IsNull", "Null"}}
	LessThan               = PartType{1, []string{"IsLessThan", "LessThan"}}
	LessThanEqual          = PartType{1, []string{"IsLessThanEqual", "LessThanEqual"}}
	GreaterThan            = PartType{1, []string{"IsGreaterThan", "GreaterThan"}}
	GreaterThanEqual       = PartType{1, []string{"IsGreaterThanEqual", "GreaterThanEqual"}}
	Before                 = PartType{1, []string{"IsBefore", "Before"}}
	After                  = PartType{1, []string{"IsAfter", "After"}}
	NotLike                = PartType{1, []string{"IsNotLike", "NotLike"}}
	Like                   = PartType{1, []string{"IsLike", "Like"}}
	StartingWith           = PartType{1, []string{"IsStartingWith", "StartingWith"}}
	EndingWith             = PartType{1, []string{"IsEndingWith", "EndingWith"}}
	IsNotEmpty             = PartType{0, []string{"IsNotEmpty", "NotEmpty"}}
	IsEmpty                = PartType{0, []string{"IsEmpty", "Empty"}}
	NotContaining          = PartType{1, []string{"IsNotContaining", "NotContaining"}}
	Containing             = PartType{1, []string{"IsContaining", "Containing"}}
	NotIn                  = PartType{1, []string{"IsNotIn", "NotIn"}}
	In                     = PartType{1, []string{"IsIn", "In"}}
	Near                   = PartType{1, []string{"IsNear", "Near"}}
	WithIn                 = PartType{1, []string{"IsWithIn", "WithIn"}}
	Regex                  = PartType{1, []string{"IsRegex", "Regex"}}
	Exists                 = PartType{0, []string{"Exists"}}
	True                   = PartType{0, []string{"IsTrue", "True"}}
	False                  = PartType{0, []string{"IsFalse", "False"}}
	NegatingSimpleProperty = PartType{1, []string{"IsNot", "Not"}}
	SimpleProperty         = PartType{1, []string{"Is", "Equals"}}
)

var (
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
		Near,
		WithIn,
		Regex,
		Exists,
		True,
		False,
		NegatingSimpleProperty,
		SimpleProperty,
	}

	AllKeywords = func() []string {
		var keywords []string
		for _, partType := range All {
			keywords = append(keywords, partType.keywords...)
		}
		return keywords
	}()
)

type PartType struct {
	numberOfArguments int
	keywords          []string
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
	for _, keyword := range p.keywords {
		if strings.HasSuffix(rawProperty, keyword) {
			return true
		}
	}
	return false
}

func (p *PartType) extractProperty(part string) string {
	for _, keyword := range p.keywords {
		if strings.HasSuffix(part, keyword) {
			return strings.TrimSuffix(part, keyword)
		}
	}
	return part
}
