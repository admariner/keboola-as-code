package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cast"

	"github.com/keboola/keboola-as-code/internal/pkg/sql"
)

type UsedSharedCodeRows []ConfigRowKey

type LinkToSharedCode struct {
	Config ConfigKey
	Rows   UsedSharedCodeRows
}

type Transformation struct {
	Blocks           []*Block
	LinkToSharedCode *LinkToSharedCode
}

// Block - transformation block.
type Block struct {
	BlockKey
	PathInProject `json:"-"`
	Name          string `json:"name" validate:"required" metaFile:"true"`
	Codes         Codes  `json:"codes" validate:"omitempty,dive"`
}

type Codes []*Code

// Code - transformation code.
type Code struct {
	CodeKey
	PathInProject `json:"-"`
	CodeFileName  string  `json:"-"` // eg. "code.sql", "code.py", ...
	Name          string  `json:"name" validate:"required" metaFile:"true"`
	Scripts       Scripts `json:"script"` // scripts, eg. SQL statements
}

type Scripts []Script

type Script interface {
	Content() string
	Clone() Script
}

// StaticScript is script defined by user (it is not link to shared code).
type StaticScript struct {
	Value string
}

func (v UsedSharedCodeRows) IdsSlice() []interface{} {
	var ids []interface{}
	for _, rowKey := range v {
		ids = append(ids, rowKey.Id.String())
	}
	return ids
}

func (v *Transformation) VisitCodes(callback func(code *Code)) {
	for _, block := range v.Blocks {
		for _, code := range block.Codes {
			callback(code)
		}
	}
}

func (v *Transformation) VisitScripts(callback func(code *Code, script Script)) {
	for _, block := range v.Blocks {
		for _, code := range block.Codes {
			for _, script := range code.Scripts {
				callback(code, script)
			}
		}
	}
}

func (v *Transformation) MapScripts(callback func(code *Code, script Script) Script) {
	for _, block := range v.Blocks {
		for _, code := range block.Codes {
			for index, script := range code.Scripts {
				code.Scripts[index] = callback(code, script)
			}
		}
	}
}

func (v *Transformation) Clone() *Transformation {
	if v == nil {
		return nil
	}
	clone := *v
	clone.LinkToSharedCode = v.LinkToSharedCode.Clone()
	clone.Blocks = make([]*Block, len(v.Blocks))
	for i, block := range v.Blocks {
		clone.Blocks[i] = block.Clone()
	}
	return &clone
}

func (v *LinkToSharedCode) Clone() *LinkToSharedCode {
	if v == nil {
		return nil
	}
	clone := *v
	clone.Rows = make([]ConfigRowKey, len(v.Rows))
	for i, row := range v.Rows {
		clone.Rows[i] = row
	}
	return &clone
}

func (k Kind) IsBlock() bool {
	return k.Name == BlockKind
}

func (k Kind) IsCode() bool {
	return k.Name == CodeKind
}

func (b *Block) Clone() *Block {
	clone := *b
	clone.Codes = b.Codes.Clone()
	return &clone
}

func (c *Code) Clone() *Code {
	clone := *c
	clone.Scripts = c.Scripts.Clone()
	return &clone
}

func (v Codes) Clone() Codes {
	if v == nil {
		return nil
	}

	out := make(Codes, len(v))
	for index, item := range v {
		out[index] = item.Clone()
	}
	return out
}

func (v Scripts) Clone() Scripts {
	if v == nil {
		return nil
	}

	out := make(Scripts, len(v))
	for index, item := range v {
		out[index] = item.Clone()
	}
	return out
}

func (b Block) String() string {
	buf := new(bytes.Buffer)
	_, _ = fmt.Fprintln(buf, `# `, b.Name)
	for _, code := range b.Codes {
		_, _ = fmt.Fprint(buf, code.String())
	}
	return buf.String()
}

func (c Code) String() string {
	return fmt.Sprintf("## %s\n%s", c.Name, c.Scripts.String(c.ComponentId))
}

func (v Scripts) Slice() []interface{} {
	var out []interface{}
	for _, script := range v {
		out = append(out, script.Content())
	}
	return out
}

func (v Scripts) String(componentId ComponentId) string {
	var items []string
	for _, script := range v {
		items = append(items, script.Content())
	}

	switch componentId.String() {
	case `keboola.snowflake-transformation`:
		fallthrough
	case `keboola.synapse-transformation`:
		fallthrough
	case `keboola.oracle-transformation`:
		return sql.Join(items) + "\n"
	default:
		return strings.Join(items, "\n") + "\n"
	}
}

// MarshalJSON converts Scripts to JSON []string.
func (v Scripts) MarshalJSON() ([]byte, error) {
	var scripts []string
	for _, script := range v {
		scripts = append(scripts, script.(StaticScript).Value)
	}
	return json.Marshal(scripts)
}

// UnmarshalJSON converts JSON string or []string to Scripts.
func (v *Scripts) UnmarshalJSON(data []byte) error {
	var script string
	var scripts []string
	if err := json.Unmarshal(data, &script); err == nil {
		scripts = append(scripts, script)
	} else if err := json.Unmarshal(data, &scripts); err != nil {
		return err
	}

	*v = make(Scripts, 0)
	for _, item := range scripts {
		*v = append(*v, StaticScript{Value: item})
	}
	return nil
}

func (v StaticScript) Content() string {
	return v.Value
}

func (v StaticScript) Clone() Script {
	return v
}

func NormalizeScript(script string) string {
	return strings.TrimRight(script, "\n\r\t ")
}

func ScriptsFromStr(content string, componentId ComponentId) Scripts {
	content = NormalizeScript(content)
	var items []string
	switch componentId.String() {
	case `keboola.snowflake-transformation`:
		fallthrough
	case `keboola.synapse-transformation`:
		fallthrough
	case `keboola.oracle-transformation`:
		items = sql.Split(content)
	default:
		items = []string{content}
	}

	scripts := make(Scripts, 0)
	for _, item := range items {
		scripts = append(scripts, StaticScript{Value: item})
	}
	return scripts
}

func ScriptsFromSlice(items []interface{}) Scripts {
	var scripts Scripts
	for _, item := range items {
		scripts = append(scripts, StaticScript{Value: cast.ToString(item)})
	}
	return scripts
}
