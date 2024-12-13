package validators

import (
	"log"

	"github.com/go-playground/validator/v10"
)

func RegisterValidators(validatorObject *validator.Validate) {
	err := validatorObject.RegisterValidation("categoryCreateType", categoryCreateTypeValidator)
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

	err = validatorObject.RegisterValidation("financialRole", financialRoleValidator)
	if err != nil {
		log.Fatalf("[Panic] - RegisterValidators - registering financialRoleValidator")
	}
}
