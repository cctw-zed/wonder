package container

import (
	"database/sql"
	"github.com/cctw-zed/wonder/internal/application/service"
	"github.com/cctw-zed/wonder/internal/interfaces/http"
)

type Container struct {
	UserHandler *http.UserHandler
}

func NewContainer(db *sql.DB) *Container {
	// 基础设施层
	idGen := snowflake.NewGenerator(1)
	userRepo := repository.NewUserRepository(db)

	// 应用服务层
	userService := service.NewUserService(userRepo, idGen)

	// 接口层
	userHandler := http.NewUserHandler(userService)

	return &Container{
		UserHandler: userHandler,
	}
}
