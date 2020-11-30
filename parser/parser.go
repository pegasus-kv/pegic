package parser

import (
	"fmt"
	"strconv"
	"unicode"
)

type ParseResult struct {
	Output interface{}
	Err    error
}

func parseOk(output interface{}) *ParseResult {
	return &ParseResult{Output: output}
}

func parseErr(err error) *ParseResult {
	return &ParseResult{Err: err}
}

type Parser func(string) (*ParseResult, string)

func TagNoCase(t string) Parser {
	return func(s string) (*ParseResult, string) {
		if !hasPrefixNoCase(s, t) {
			return parseErr(fmt.Errorf("expected `%s`, found `%s`", t, s)), s
		}
		s = TrimPrefixNoCase(s, t)
		return parseOk(t), s
	}
}

func Alt(parsers ...Parser) Parser {
	return func(s string) (*ParseResult, string) {
		var last error
		for _, parser := range parsers {
			res, input := parser(s)
			if res.Err == nil {
				return res, input
			}
			last = res.Err
		}
		return parseErr(last), s
	}
}

func Array(parsers ...Parser) Parser {
	return func(s string) (*ParseResult, string) {
		var arr []interface{}
		for _, parser := range parsers {
			res, input := parser(s)
			if res.Err != nil {
				return res, input
			}
			s = input
			arr = append(arr, res.Output)
		}
		return parseOk(arr), s
	}
}

func ArrayWhiteSpace(parsers ...Parser) Parser {
	return func(s string) (*ParseResult, string) {
		var arr []interface{}
		lastSpace := true
		var ps string
		for i, parser := range parsers {
			res, input := parser(s)
			if res.Err != nil {
				return res, input
			}
			if res.Output != nil && !lastSpace {
				return parseErr(fmt.Errorf("expected whitespace, found `%s`", s)), input
			}
			s = input
			if res.Output != nil {
				ps = s
			}
			arr = append(arr, res.Output)
			if res.Output != nil && i != len(parsers)-1 {
				res, input = MultiSpace1(s)
				if res.Err != nil {
					lastSpace = false
				}
				s = input
			}
		}
		return parseOk(arr), ps
	}
}

func TakeWhile1(predicate func(rune) bool) Parser {
	return func(s string) (*ParseResult, string) {
		var first rune = 0
		for _, c := range s {
			first = c
			break
		}
		if first == 0 || !predicate(first) {
			return parseErr(fmt.Errorf("expected predicate, found `%s`", s)), s
		}
		chars, input := TakeWhile(predicate)(s)
		return parseOk(chars.Output), input
	}
}

func TakeWhile(predicate func(rune) bool) Parser {
	return func(s string) (*ParseResult, string) {
		idx := 0
		for _, ch := range s {
			if predicate(ch) {
				idx++
			} else {
				break
			}
		}
		return parseOk(s[0:idx]), s[idx:]
	}
}

func String(s string) (*ParseResult, string) {
	res, input := Array(TagNoCase(`"`), TakeWhile(func(r rune) bool { return r != '"' }), TagNoCase(`"`))(s)
	if res.Err != nil {
		return parseErr(fmt.Errorf("expected string, found `%s`", s)), s
	}
	output := res.Output.([]interface{})
	return parseOk(output[1]), input
}

func UInt(s string) (*ParseResult, string) {
	res, input := TakeWhile1(unicode.IsDigit)(s)
	if res.Err != nil {
		return parseErr(fmt.Errorf("expected uint, found `%s`", s)), s
	}
	n, err := strconv.ParseUint(res.Output.(string), 10, 32)
	if err != nil {
		return parseErr(fmt.Errorf("cannot parse uint from `%s`: %s", s, err)), s
	}
	return parseOk(uint(n)), input
}

func MultiSpace1(s string) (*ParseResult, string) {
	res, input := TakeWhile1(unicode.IsSpace)(s)
	if res.Err != nil {
		return parseErr(fmt.Errorf("expected whitespace, found `%s`", s)), s
	}
	return res, input
}

func MultiSpace0(s string) (*ParseResult, string) {
	res, input := TakeWhile(unicode.IsSpace)(s)
	if res.Err != nil {
		return parseErr(fmt.Errorf("expected whitespace, found `%s`", s)), s
	}
	return res, input
}

func Maybe(parsers ...Parser) Parser {
	return func(s string) (*ParseResult, string) {
		var arr []interface{}
		for i, parser := range parsers {
			res, input := parser(s)
			if res.Err != nil {
				if i != 0 {
					return res, s
				}
				return parseOk(nil), s
			}
			arr = append(arr, res.Output)
			s = input
		}
		return parseOk(arr), s
	}
}

func Opt(parser Parser) Parser {
	return func(s string) (*ParseResult, string) {
		res, input := parser(s)
		if res.Err != nil {
			return parseOk(nil), s
		}
		return res, input
	}
}
