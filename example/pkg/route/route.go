package route

import (
	"github.com/zhangliang-zl/reskit/example/pkg/controlelr"
	"github.com/zhangliang-zl/reskit/server"
)

func Init(srv *server.Server) {
	userCtrl := &controlelr.UserController{}
	srv.Route("GET", "/user/info", userCtrl.Info)
}
