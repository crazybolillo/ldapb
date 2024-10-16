package internal

import (
	"github.com/crazybolillo/eryth/pkg/model"
	"net/url"
	"regexp"
)

var attributeRegex = regexp.MustCompile(`\((\w+)=([*\w]+)\)`)

func parseFilter(filter string) model.ContactPageFilter {
	res := model.ContactPageFilter{}

	matches := attributeRegex.FindAllStringSubmatch(filter, -1)
	for _, match := range matches {
		attribute := match[1]
		switch attribute {
		case "cn":
			res.Name = match[2]
		case "telephoneNumber":
			res.Phone = match[2]
		}
	}

	// Can this abomination be called heuristics?
	if len(filter) > 2 && []rune(filter)[1] == '|' {
		res.Operator = "or"
	}

	return res
}

func filterToHttp(filter string) string {
	parsed := parseFilter(filter)

	values := url.Values{}
	values.Set("name", parsed.Name)
	values.Set("phone", parsed.Phone)
	values.Set("op", parsed.Operator)

	return values.Encode()
}
