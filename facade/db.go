package facade

import (
	"github.com/zhangliang-zl/reskit"
	"github.com/zhangliang-zl/reskit/component"
	"github.com/zhangliang-zl/reskit/db"
	"gorm.io/gorm"
)

func RegisterDB(opts db.Options, id string) error {
	logger := Logger(compDb)
	client, err := db.New(opts, logger)
	if err != nil {
		return err
	}

	closeFunc := func() error {
		s, err := client.DB()
		if err != nil {
			return err
		}
		return s.Close()
	}

	instance := component.Make(client, nil, closeFunc)
	return reskit.App().SetComponent(compDb, id, instance)
}

func DB(id string) *gorm.DB {
	res, ok := reskit.App().Component(compDb, id)

	if !ok {
		panic(compDb + noRegister)
	}

	return res.(*gorm.DB)
}
