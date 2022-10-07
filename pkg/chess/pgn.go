package chess

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type PGNResult struct {
	Start Position
	Moves []Move
	Tags  map[string]string
}

func parseTag(r io.RuneScanner) (key string, value string, err error) {
	var (
		krs []rune
		vrs []rune
	)

	const (
		KEY = iota
		SPACE
		VALUE
	)
	state := KEY

	escape := false
	Loop:
	for {
		c, _, err := r.ReadRune()
		if err != nil {
			return key, value, err
		}
		switch state {
		case KEY:
			if unicode.IsSpace(c) {
				state = SPACE
			} else {
				krs = append(krs, c)
			}
		case SPACE:
			if c == '"' {
				state = VALUE
			} else if !unicode.IsSpace(c) {
				return key, value, fmt.Errorf("unexpected %U", c)
			}
		case VALUE:
			if escape {
				if c != '\\' && c != '"' {
					return key, value, fmt.Errorf("invalid escape sequence \"\\%c\"", c)
				}
				vrs = append(vrs, c)
				escape = false
				continue
			}
			if c == '\\' {
				escape = true
				continue
			}
			if c == '"' {
				break Loop
			}
			vrs = append(vrs, c)
		}
	}

	c, _, err := r.ReadRune()
	if err != nil {
		return key, value, err
	} else if c != ']' {
		return key, value, fmt.Errorf("expected \"]\" (%U), got %U", ']', c)
	}

	return string(krs), string(vrs), nil
}

func parsePGNTags(r io.RuneScanner) (map[string]string, error) {
	tags := make(map[string]string)
	for {
		c, _, err := r.ReadRune()
		if err != nil {
			return tags, err
		}
		if unicode.IsSpace(c) {
			continue
		}
		if c == '[' {
			key, value, err := parseTag(r)
			if err != nil {
				return tags, err
			}
			tags[key] = value
		} else {
			return tags, r.UnreadRune()
		}
	}
}

func ParsePGN(r io.Reader) (PGNResult, error) {
	res := PGNResult{}

	var rs io.RuneScanner
	if rrs, ok := r.(io.RuneScanner); ok {
		rs = rrs
	} else {
		rs = bufio.NewReader(r)
	}

	tags, err := parsePGNTags(rs)
	if err != nil {
		return res, err
	}
	res.Tags = tags

	if notation, exists := tags["FEN"]; exists {
		pos, err := FEN().Parse(strings.NewReader(notation))
		if err != nil {
			return res, err
		}
		res.Start = pos
	} else {
		res.Start = StartingPosition()
	}

	movs, err := Algebraic().Parse(res.Start, rs)
	if err != nil {
		return res, nil
	}
	res.Moves = movs

	return res, nil
}