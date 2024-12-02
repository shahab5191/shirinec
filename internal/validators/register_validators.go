package validators

import (
	"log"

	"github.com/go-playground/validator/v10"
)

func RegisterValidators(validatorObject *validator.Validate) {
	err := validatorObject.RegisterValidation("mediaUploadBind", mediaUploadBindValidator)
	if err != nil {
		log.Fatalf("[Panic] - RegisterValidators - registering mediaUploadBind")
	}

	err = validatorObject.RegisterValidation("categoryCreateType", categoryCreateTypeValidator)
	if err != nil {
		log.Fatalf("[Panic] - RegisterValidators - registering categoryCreateType")
	}

	err = validatorObject.RegisterValidation("alphaNumericSpace", alphaNumericSpaceValidator)
	if err != nil {
		log.Fatalf("[Panic] - RegisterValidators - registering alphaNumericSpace")
	}
    
    err = validatorObject.RegisterValidation("intLen", intLenValidator)
    if err != nil {
		log.Fatalf("[Panic] - RegisterValidators - registering intLen")
	}
}
