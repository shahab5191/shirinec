package validators

import (
	"regexp"
	"strconv"

	"github.com/go-playground/validator/v10"
)

func alphaNumericSpaceValidator(fl validator.FieldLevel) bool {
	regexAlphaNumericSpace := regexp.MustCompile("^[a-zA-Z0-9 ]*$")
	return regexAlphaNumericSpace.MatchString(fl.Field().String())
}

func intLenValidator(fl validator.FieldLevel) bool {
	param := fl.Param()
	expectedLen, err := strconv.Atoi(param)
	if err != nil {
		return false
	}

	num := fl.Field().Int()
	actualLen := len(strconv.Itoa(int(num)))

	return actualLen == expectedLen
}
