package validation

import (
	"testing"
)

type Student struct {
	Phone  string `validate:"phone" valid_zh:"手机号"`
	Email  string `validate:"email" valid_zh:"电子邮件" valid_en:"eemail"`
	IDcard string `validate:"idcard" valid_zh:"身份证" `
	Age    int    `validate:"max=1,min=12" valid_zh:"年龄" `
}

func TestValidatorStruct(t *testing.T) {
	student := Student{
		Phone:  "",
		Email:  "testemal",
		Age:    40,
		IDcard: "a1123123",
	}

	v := New(LangZH)
	result := v.Struct(student)
	if !result.HasError() {
		t.Error("struct error")
	}
}

func TestValidatorVar(t *testing.T) {
	v := New(LangEN)
	result := v.Var("123@qq.com", "email")
	if result.HasError() {
		t.Error(result.FirstError())
	}
}
