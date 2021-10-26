package persist

import "github.com/go-kratos/kratos/v2/log"

var LogHelper = log.NewHelper(log.With(log.DefaultLogger, "project", "test-1"))

var DBLogger log.Logger
