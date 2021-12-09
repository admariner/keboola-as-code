package orderedmap

import (
	"fmt"
	"strings"

	"github.com/spf13/cast"
	"github.com/umisama/go-regexpcache"
)

// KeyPath - path to a value in the OrderedMap (JSON file).
type Key []Step

type Step interface {
	String() string
}

type MapStep string

type SliceStep int

func KeyFromStr(str string) Key {
	parts := strings.FieldsFunc(str, func(r rune) bool {
		return r == '.' || r == '['
	})

	out := make(Key, 0)
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}

		// Is slice step? eg. [123]
		matches := regexpcache.MustCompile(`^(\d+)\]$`).FindStringSubmatch(part)
		if matches != nil {
			out = append(out, SliceStep(cast.ToInt(matches[1])))
		} else {
			out = append(out, MapStep(part))
		}
	}

	return out
}

func (v Key) String() string {
	parts := make([]string, 0)
	for _, step := range v {
		parts = append(parts, step.String())
	}
	return strings.ReplaceAll(strings.Join(parts, "."), `.[`, `[`)
}

func (v Key) WithoutLast() Key {
	l := len(v)
	if l == 0 {
		return nil
	}
	return v[0 : l-1]
}

func (v Key) Last() Step {
	l := len(v)
	if l == 0 {
		return nil
	}
	return v[l-1]
}

func (v MapStep) String() string {
	return string(v)
}

func (v SliceStep) String() string {
	return fmt.Sprintf("[%v]", int(v))
}