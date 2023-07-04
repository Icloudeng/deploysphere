package lib

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func ResourcesRefValidation(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	match, _ := regexp.MatchString("^[a-z]+([0-9a-z]+(?:-[0-9a-z]+)?)*$", value)

	return match
}
