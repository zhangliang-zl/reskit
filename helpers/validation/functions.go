package validation

import (
	"github.com/go-playground/validator/v10"
	"github.com/zhangliang-zl/reskit/helpers/idcard"
	"regexp"
)

func IsIDCard(id string) bool {
	return idcard.NewParser(id).Validate()
}

func IsPhone(phone string) bool {
	pattern := `^(\+?86-?)?(1)[0-9]{10}$`
	matched, _ := regexp.MatchString(pattern, phone)
	return matched
}

func IsZipCode(v string) bool {
	pattern := `^[0-9a-zA-Z]{4,20}$`
	matched, _ := regexp.MatchString(pattern, v)
	return matched
}

func IsNameSpelling(v string) bool {
	pattern := `^[A-Z]{0,100}$`
	matched, _ := regexp.MatchString(pattern, v)
	return matched
}

func IsTel(v string) bool {
	pattern := `^[0-9-]{4,20}$`
	matched, _ := regexp.MatchString(pattern, v)
	return matched
}

const (
	invalidZhTip = "{0}不合法"
	invalidEnTip = "{0} invalid"
)

func (v *Validator) registerCustomFunctions() {
	v.RegisterCustomerFunc("idcard", invalidZhTip, invalidEnTip, func(fl validator.FieldLevel) bool {
		return IsIDCard(fl.Field().String())
	})

	v.RegisterCustomerFunc("phone", invalidZhTip, invalidEnTip, func(fl validator.FieldLevel) bool {
		return IsPhone(fl.Field().String())
	})

	v.RegisterCustomerFunc("tel", invalidZhTip, invalidEnTip, func(fl validator.FieldLevel) bool {
		return IsTel(fl.Field().String())
	})

	v.RegisterCustomerFunc("namespelling", invalidZhTip, invalidEnTip, func(fl validator.FieldLevel) bool {
		return IsNameSpelling(fl.Field().String())
	})

	v.RegisterCustomerFunc("zipcode", invalidZhTip, invalidEnTip, func(fl validator.FieldLevel) bool {
		return IsZipCode(fl.Field().String())
	})
}
