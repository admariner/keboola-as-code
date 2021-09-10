package json

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// ReadFile reads JSON file to target data.
func ReadFile(dir string, relPath string, target interface{}, errPrefix string) error {
	path := filepath.Join(dir, relPath)

	// Read file
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("missing %s file \"%s\"", errPrefix, relPath)
		}
		return fmt.Errorf("cannot read %s file \"%s\"", errPrefix, relPath)
	}

	// Decode file
	err = Decode(content, target)
	if err != nil {
		return fmt.Errorf("%s file \"%s\" is invalid:\n\t- %w", errPrefix, relPath, err)
	}
	return nil
}

// WriteFile writes JSON file from source data.
func WriteFile(dir string, relPath string, source interface{}, errPrefix string) error {
	path := filepath.Join(dir, relPath)
	data, err := Encode(source, true)
	if err != nil {
		return fmt.Errorf("cannot write %s file \"%s\"", errPrefix, relPath)
	}
	return os.WriteFile(path, data, 0644)
}

func Encode(v interface{}, pretty bool) ([]byte, error) {
	var data []byte
	var err error
	if pretty {
		data, err = json.MarshalIndent(v, "", "  ")
		data = append(data, '\n')
	} else {
		data, err = json.Marshal(v)
	}
	if err != nil {
		return nil, processJsonError(err)
	}
	return data, nil
}

func MustEncode(v interface{}, pretty bool) []byte {
	data, err := Encode(v, pretty)
	if err != nil {
		panic(err)
	}
	return data
}

func EncodeString(v interface{}, pretty bool) (string, error) {
	data, err := Encode(v, pretty)
	return string(data), err
}

func MustEncodeString(v interface{}, pretty bool) string {
	data, err := EncodeString(v, pretty)
	if err != nil {
		panic(err)
	}
	return data
}

func Decode(data []byte, m interface{}) error {
	if err := json.Unmarshal(data, m); err != nil {
		return processJsonError(err)
	}
	return nil
}

func MustDecode(data []byte, m interface{}) {
	if err := Decode(data, m); err != nil {
		panic(err)
	}
}

func DecodeString(data string, m interface{}) error {
	return Decode([]byte(data), m)
}

func MustDecodeString(data string, m interface{}) {
	if err := DecodeString(data, m); err != nil {
		panic(err)
	}
}

func processJsonError(err error) error {
	var typeError *json.UnmarshalTypeError
	var syntaxError *json.SyntaxError

	switch {
	// Custom error message
	case errors.As(err, &typeError):
		return fmt.Errorf("key \"%s\" has invalid type \"%s\"", typeError.Field, typeError.Value)
	case errors.As(err, &syntaxError):
		return fmt.Errorf("%w, offset: %d", err, syntaxError.Offset)
	default:
		return err
	}
}
