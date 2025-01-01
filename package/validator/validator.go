package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validate is a struct that holds the validator implementation.
type Validate struct {
	*validator.Validate
}

// New initializes a new validator.
func New() *Validate {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fl reflect.StructField) string {
		name := strings.SplitN(fl.Tag.Get("json"), ",", 2)
		return name[0]
	})
	return &Validate{validate}
}

// ValidatorErrors func for show validation errors for each invalid fields.
func ValidatorErrors(err error) map[string]string {
	// Define variable for error fields.
	errFields := map[string]string{}

	// Make error message for each invalid field.
	for _, err := range err.(validator.ValidationErrors) {
		// Get name of the field's struct.
		structName := strings.Split(err.Namespace(), ".")[0]
		// --> first (0) element is the founded name

		// Append error message to the map, where key is a field name,
		// and value is an error description.
		errFields[err.Field()] = fmt.Sprintf(
			"failed '%s' tag check (value '%s' is not valid for %s struct)",
			err.Tag(), err.Value(), structName,
		)
	}

	return errFields
}
