package validators

import (
	"github.com/go-playground/validator/v10"
)

func RegisterValidators(validatorObject *validator.Validate) {
	validatorObject.RegisterValidation("mediaUploadBind", mediaUploadBindValidator)
}