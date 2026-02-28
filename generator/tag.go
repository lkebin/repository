package generator

import "strings"

type TagOptions string

func ParseTag(tag string) (string, TagOptions) {
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx], TagOptions(tag[idx+1:])
	}
	return tag, TagOptions("")
}

func (t TagOptions) Get(optionName string) string {
	if len(t) == 0 {
		return ""
	}

	s := string(t)
	for s != "" {
		var next string
		i := strings.Index(s, ",")
		if i >= 0 {
			s, next = s[:i], s[i+1:]
		}

		j := strings.Index(s, "=")
		if j >= 0 {
			k, v := s[:j], s[j+1:]
			if k == optionName {
				return v
			}
		}

		s = next
	}

	return ""
}

func (t TagOptions) Contains(optionName string) bool {
	if len(t) == 0 {
		return false
	}
	s := string(t)
	for s != "" {
		var next string
		i := strings.Index(s, ",")
		if i >= 0 {
			s, next = s[:i], s[i+1:]
		}
		if s == optionName {
			return true
		}
		s = next
	}
	return false
}
