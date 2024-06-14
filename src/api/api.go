package api

import (
	"goddd/src/infra/repositories/memory"
	"goddd/src/services"
	"log"

	"github.com/gin-gonic/gin"
)

type Api struct {
	router *gin.Engine
	cfg Config
}

func NewApi() Api {
	// init all repositories
	repo := memory.NewMemoryRepository()
	// init services
	services, err := services.NewServices(repo)
	if err != nil {
		log.Panic(err)
	}
	return Api{router: NewRouter(services), cfg: NewConfig()}
}

func (a Api) Run() {
	// //Starting the application
	log.Fatal(a.router.Run(":" + a.cfg.Port))
}
