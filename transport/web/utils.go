package web

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/zhangliang-zl/reskit/helper/validation"
)

func BindValidator(validator *validation.Validator) {
	binding.Validator = validator
}

type Map map[string]interface{}
