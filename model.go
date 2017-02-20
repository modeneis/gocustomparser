package gocustomparser

import "golang.org/x/text/encoding"

type Parser interface {
	Parse(resultType interface{})
}

type CustomParser struct {
	File            string
	Separator       rune
	SkipFirstLine   bool
	SkipEmptyValues bool
	PRNReader       func(raw string) (line []string, err error)
	CustomDecoder   *encoding.Decoder
}
