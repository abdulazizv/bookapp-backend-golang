package main

import (
	"github.com/casbin/casbin/v2"
	defaultrolemanager "github.com/casbin/casbin/v2/rbac/default-role-manager"
	"github.com/casbin/casbin/v2/util"
	"github.com/gomodule/redigo/redis"
	"gitlab.com/bookapp/api"
	"gitlab.com/bookapp/config"
	"gitlab.com/bookapp/pkg/db"
	"gitlab.com/bookapp/pkg/logger"
	"gitlab.com/bookapp/storage"
	r "gitlab.com/bookapp/storage/redis"
)

func main() {
	cfg := config.Load()
	log := logger.New("debug", "bookapp")

	var (
		casbinEnforcer *casbin.Enforcer
	)

	casbinEnforcer, err := casbin.NewEnforcer(cfg.AuthFilePath, cfg.CsvFilePath)
	if err != nil {
		log.Error("casbin enforcer error: ", logger.Error(err))
		return
	}

	err = casbinEnforcer.LoadPolicy()
	if err != nil {
		log.Error("casbin error load policy", logger.Error(err))
		return
	}

	casbinEnforcer.GetRoleManager().(*defaultrolemanager.RoleManagerImpl).AddMatchingFunc("keyMatch", util.KeyMatch)
	casbinEnforcer.GetRoleManager().(*defaultrolemanager.RoleManagerImpl).AddMatchingFunc("keyMatch3", util.KeyMatch3)

	connDb, err := db.ConnectToDb(cfg)
	if err != nil {
		logger.Error(err)
	}
	strg := storage.NewStoragePg(connDb)

	pool := &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", cfg.RedisHost+":"+cfg.RedisPort)
		},
	}

	apiServer := api.New(&api.Options{
		Cfg:            cfg,
		Storage:        strg,
		Log:            log,
		Redis:          r.NewRedisRepo(pool),
		CasbinEnforcer: casbinEnforcer,
	})

	if err := apiServer.Run(cfg.HttpPort); err != nil {
		log.Fatal("failed to run server: %v", logger.Error(err))
	}

}
