package helper

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

type Validator struct {
	validator *validator.Validate
	trans     ut.Translator
}

func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

func NewValidator() *Validator {
	validate := validator.New()

	// Setup English translator
	english := en.New()
	uni := ut.New(english, english)
	trans, _ := uni.GetTranslator("en")
	_ = enTranslations.RegisterDefaultTranslations(validate, trans)

	return &Validator{
		validator: validate,
		trans:     trans,
	}
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (v *Validator) TranslateError(err error) (errs []FieldError) {
	if err == nil {
		return nil
	}

	validatorErrs := err.(validator.ValidationErrors)
	for _, e := range validatorErrs {
		translatedErr := e.Translate(v.trans)

		// Append a detailed error for each field
		errs = append(errs, FieldError{
			Field:   e.Field(),     // Field name
			Message: translatedErr, // Translated error message
		})
	}
	return errs
}
