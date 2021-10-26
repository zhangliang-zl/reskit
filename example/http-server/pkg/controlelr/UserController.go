package controlelr

import (
	"github.com/zhangliang-zl/reskit/example/http-server/pkg/service"
	"github.com/zhangliang-zl/reskit/web"
)

type UserController struct {
}

func (*UserController) Info(ctx *web.Context) {
	id := ctx.DefaultQuery("id", "0")
	userSvc := service.NewUserService()
	info, err := userSvc.Info(id)
	if err != nil {
		ctx.BadRequest(err)
		return
	}
	ctx.Success(info)
}
