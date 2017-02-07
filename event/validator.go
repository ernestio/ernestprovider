package event

import (
	"errors"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
)

// Validator : Event validator
type Validator struct {
	Val   *validator.Validate
	Uni   *ut.UniversalTranslator
	Trans ut.Translator
}

// NewValidator ...
func NewValidator() *Validator {
	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")
	validate := validator.New()
	_ = en_translations.RegisterDefaultTranslations(validate, trans)

	return &Validator{Val: validate, Uni: uni, Trans: trans}
}

// Validate : ...
func (v *Validator) Validate(ev interface{}) error {
	err := v.Val.Struct(ev)
	if err != nil {
		msg := ""
		errs := err.(validator.ValidationErrors)
		for _, e := range errs {
			msg += e.Translate(v.Trans) + "\n"
		}
		return errors.New(msg)
	}
	return nil
}
