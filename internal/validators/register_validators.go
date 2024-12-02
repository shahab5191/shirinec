package validators

import (
	"github.com/go-playground/validator/v10"
)

func RegisterValidators(validatorObject *validator.Validate) {
    validatorObject.RegisterValidation("mediaUploadBind", mediaUploadBindValidator)
    validatorObject.RegisterValidation("categoryCreateType", categoryCreateTypeValidator)
    validatorObject.RegisterValidation("alphaNumericSpace", alphaNumericSpaceValidator)
}
