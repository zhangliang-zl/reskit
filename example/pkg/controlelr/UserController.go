package controlelr

import (
	"github.com/zhangliang-zl/reskit/example/pkg/service"
	"github.com/zhangliang-zl/reskit/server"
)

type UserController struct {
}

func (*UserController) Info(ctx *server.Context) {
	id := ctx.DefaultQuery("id", "0")
	userSvc := service.NewUserService()
	info, err := userSvc.Info(id)
	if err != nil {
		ctx.BadRequest(err)
		return
	}
	ctx.Success(info)
}
