package validation

import (
	"errors"
	"fmt"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entrans "github.com/go-playground/validator/v10/translations/en"
	zhtrans "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
	"strings"
)

const (
	LangZH = "zh"
	LangEN = "en"
	TagEN  = "valid_en"
	TagZH  = "valid_zh"
)

type Validator struct {
	*validator.Validate
	lang    string
	tagName string

	translator   ut.Translator
	zhTranslator ut.Translator
	enTranslator ut.Translator
}

func (v *Validator) init() {
	v.registerTagName(v.tagName)
	zhtrans.RegisterDefaultTranslations(v.Validate, v.zhTranslator)
	entrans.RegisterDefaultTranslations(v.Validate, v.enTranslator)

	v.registerCustomFunctions()
}

func (v *Validator) registerTagName(tagName string) {
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get(tagName), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

func (v *Validator) RegisterCustomerFunc(tag string, zhTip, enTip string, fn validator.Func) {
	v.RegisterValidation(tag, fn)
	v.registerTranslate(v.zhTranslator, tag, zhTip)
	v.registerTranslate(v.enTranslator, tag, enTip)
}

func (v *Validator) registerTranslate(translator ut.Translator, tag string, text string) {
	v.RegisterTranslation(tag, translator, func(ut ut.Translator) error {
		return ut.Add(tag, text, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(tag, fe.Field())
		return t
	})
}

func (v *Validator) Struct(val interface{}) Result {
	err := v.Validate.Struct(val)
	return v.parseResult(err)
}

func (v *Validator) ValidateStruct(val interface{}) error {
	err := v.Validate.Struct(val)
	result := v.parseResult(err)
	if result.hasError {
		return errors.New(result.firstError)
	}

	return nil
}

func (v *Validator) Engine() interface{} {
	return v
}

func (v *Validator) Var(field interface{}, tag string) Result {
	err := v.Validate.Var(field, tag)
	return v.parseResult(err)
}

func (v *Validator) parseResult(err error) Result {
	if err != nil {
		// invalid error
		invalid, ok := err.(*validator.InvalidValidationError)
		if ok {
			return Result{
				paramsIsValid: false,
				hasError:      true,
				errors: map[string]string{
					"param error:": invalid.Error(),
				},
				firstError: fmt.Sprintf("param error: %s", invalid.Error()),
			}
		}

		// validation errors
		validErrors, ok := err.(validator.ValidationErrors)
		if ok {
			errs := validErrors.Translate(v.translator)
			firstError := ""
			if len(validErrors) > 0 {
				firstError = validErrors[0].Translate(v.translator)
			}

			return Result{
				paramsIsValid: true,
				hasError:      true,
				errors:        errs,
				firstError:    firstError,
			}
		}
	}

	return Result{
		paramsIsValid: true,
	}
}

func New(lang string) *Validator {

	var translator ut.Translator
	var tagName string

	uni := ut.New(zh.New(), en.New())
	zhTranslator, _ := uni.GetTranslator("zh")
	enTranslator, _ := uni.GetTranslator("en")

	if lang == LangZH {
		tagName = TagZH
		translator = zhTranslator
	} else {
		lang = LangEN
		tagName = TagEN
		translator = enTranslator
	}

	v := &Validator{
		Validate:     validator.New(),
		translator:   translator,
		zhTranslator: zhTranslator,
		enTranslator: enTranslator,
		tagName:      tagName,
	}

	v.init()
	return v
}
