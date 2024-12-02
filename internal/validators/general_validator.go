package validators

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func alphaNumericSpaceValidator(fl validator.FieldLevel) bool{
    regexAlphaNumericSpace := regexp.MustCompile("^[a-zA-Z0-9 ]*$")
    return regexAlphaNumericSpace.MatchString(fl.Field().String())
}
