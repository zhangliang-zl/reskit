package facade

import (
	"github.com/zhangliang-zl/reskit/app"
	"github.com/zhangliang-zl/reskit/db"
	"gorm.io/gorm"
)

func RegisterDB(opts db.Options, id string) error {
	id = buildID(compDb, id)
	logger, err := Logger(compDb)
	if err != nil {
		return err
	}

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

	comp := app.MakeComponent(client, closeFunc)
	appInstance.RegisterComponent(comp, id)
	return nil
}

func DB(id string) *gorm.DB {
	id = buildID(compDb, id)
	res, ok := App().Component(id)
	if !ok {
		panic(compRedis + noRegister)
	}
	return res.Object().(*gorm.DB)
}
