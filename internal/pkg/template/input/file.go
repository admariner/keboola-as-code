package input

import (
	"github.com/keboola/keboola-as-code/internal/pkg/encoding/json"
	"github.com/keboola/keboola-as-code/internal/pkg/encoding/jsonnet"
	"github.com/keboola/keboola-as-code/internal/pkg/filesystem"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/errors"
)

const (
	FileName = "inputs.jsonnet"
)

func Path() string {
	return filesystem.Join("src", FileName)
}

type file struct {
	StepsGroups StepsGroups `json:"stepsGroups" validate:"dive"`
}

func newFile() *file {
	return &file{
		StepsGroups: make(StepsGroups, 0),
	}
}

func loadFile(fs filesystem.Fs, ctx *jsonnet.Context) (*file, error) {
	// Check if file exists
	path := Path()
	if !fs.IsFile(path) {
		return nil, errors.Errorf("file \"%s\" not found", path)
	}

	// Read file
	fileDef := filesystem.NewFileDef(path).SetDescription("inputs")
	content := newFile()
	if _, err := fs.FileLoader().WithJsonnetContext(ctx).ReadJsonnetFileTo(fileDef, content); err != nil {
		return nil, err
	}

	// Validate
	if err := content.validate(); err != nil {
		return nil, err
	}

	return content, nil
}

func saveFile(fs filesystem.Fs, content *file) error {
	// Validate
	if err := content.validate(); err != nil {
		return err
	}

	// Convert to Json
	jsonContent, err := json.EncodeString(content, true)
	if err != nil {
		return err
	}

	// Convert to Jsonnet
	jsonnetStr, err := jsonnet.Format(jsonContent)
	if err != nil {
		return err
	}

	// Write file
	f := filesystem.NewRawFile(Path(), jsonnetStr)
	if err := fs.WriteFile(f); err != nil {
		return err
	}

	return nil
}

func (f file) validate() error {
	return f.StepsGroups.ValidateDefinitions()
}
