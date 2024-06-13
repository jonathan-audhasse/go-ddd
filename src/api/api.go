package api

import (
	"goddd/src/infra/repositories/memory"
	"goddd/src/services"
	"log"

	"github.com/gin-gonic/gin"
)

func Run() {
	// init all repositories
	repo := memory.NewMemoryRepositories()
	// init services
	_, err := services.NewServices(repo.CustRepo)
	if err != nil {
		log.Panic(err)
	}

	// init apis
	// custApi, err := NewCustomerApi(services.CustService)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	r := gin.Default()
	r.Use(Cors()) //For CORS

	// routes
	// r.POST("/service", wordsApi.AddWord)
	// r.GET("/service", wordsApi.GetWordFromPrefix)

	r.GET("/health", Health)

	// //Starting the application
	log.Fatal(r.Run(":" + apiPort))
	log.Println("Hello Jonathan")
}
