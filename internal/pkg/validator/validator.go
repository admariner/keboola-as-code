package validator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/keboola/keboola-as-code/internal/pkg/utils"
)

func Validate(value interface{}) error {
	// Setup
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		if fld.Anonymous {
			return "__nested__"
		}

		// Use JSON field name in error messages
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return fld.Name
		}
		return name
	})

	// Do
	if err := validate.Struct(value); err != nil {
		var validationErrs validator.ValidationErrors
		switch {
		case errors.As(err, &validationErrs):
			return processValidateError(validationErrs)
		default:
			panic(err)
		}
	}

	return nil
}

func processValidateError(err validator.ValidationErrors) error {
	result := utils.NewMultiError()
	for _, e := range err {
		// Remove struct name, first part
		namespace := regexp.MustCompile(`^([^.]+\.)?(.*)$`).ReplaceAllString(e.Namespace(), `$2`)
		// Hide nested fields
		namespace = strings.ReplaceAll(namespace, `__nested__.`, ``)
		result.Append(fmt.Errorf(
			"key=\"%s\", value=\"%v\", failed \"%s\" validation",
			namespace,
			e.Value(),
			e.ActualTag(),
		))
	}

	return result.ErrorOrNil()
}