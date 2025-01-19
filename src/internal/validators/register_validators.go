package validators

import (
	"log"

	"github.com/go-playground/validator/v10"
)

func RegisterValidators(validatorObject *validator.Validate) {
	if err := validatorObject.RegisterValidation("categoryCreateType", categoryCreateTypeValidator); err != nil {
		log.Fatalf("[Panic] - RegisterValidators - registering categoryCreateType")
	}

	if err := validatorObject.RegisterValidation("alphaNumericSpace", alphaNumericSpaceValidator); err != nil {
		log.Fatalf("[Panic] - RegisterValidators - registering alphaNumericSpace")
	}

	if err := validatorObject.RegisterValidation("intLen", intLenValidator); err != nil {
		log.Fatalf("[Panic] - RegisterValidators - registering intLen")
	}

	if err := validatorObject.RegisterValidation("financialRole", financialRoleValidator); err != nil {
		log.Fatalf("[Panic] - RegisterValidators - registering financialRoleValidator")
	}

	if err := validatorObject.RegisterValidation("accountType", accountTypeValidator); err != nil {
		log.Fatalf("[Panic] - RegisterValidators - registering accountTypeValidator")
	}
}
