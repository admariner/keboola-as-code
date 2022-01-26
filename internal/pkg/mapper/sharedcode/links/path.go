package links

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/keboola/keboola-as-code/internal/pkg/model"
	"github.com/keboola/keboola-as-code/internal/pkg/naming"
)

const (
	PathFormat = `{{:<PATH>}}` // link to shared code used locally
	PathRegexp = `[^:{}]+`
)

type pathUtils struct {
	re *regexp.Regexp
}

func newPathUtils() *pathUtils {
	re := regexp.MustCompile(
		strings.ReplaceAll(
			`^`+regexp.QuoteMeta(PathFormat)+`$`,
			`<PATH>`,
			`(`+PathRegexp+`)`,
		),
	)
	return &pathUtils{re: re}
}

func (v *pathUtils) match(script string, componentId model.ComponentId) string {
	comment := naming.CodeFileComment(naming.CodeFileExt(componentId))
	script = strings.TrimSpace(script)
	script = strings.TrimPrefix(script, comment)
	script = strings.TrimSpace(script)
	match := v.re.FindStringSubmatch(script)
	if len(match) > 0 {
		return match[1]
	}
	return ""
}

func (v *pathUtils) format(path string, componentId model.ComponentId) string {
	placeholder := strings.ReplaceAll(PathFormat, `<PATH>`, path)
	if ok := v.re.MatchString(placeholder); !ok {
		panic(fmt.Errorf(`shared code path "%s" is invalid`, path))
	}
	comment := naming.CodeFileComment(naming.CodeFileExt(componentId))
	return comment + ` ` + placeholder
}