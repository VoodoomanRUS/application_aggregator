package main

import (
	_ "application_aggregator/docs"
	"application_aggregator/internal/handler"
	"application_aggregator/internal/repository/postgres"
	"application_aggregator/internal/service"
	"database/sql"
	"os/signal"
	"syscall"

	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Aggregator API
// @version 1.0
// @description API for managing organizations.
// @host localhost:8080
// @BasePath /api/v1
// @schemes http
func main() {
	//1. Подключение к БД
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_DSN"))
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal("Failed to close DB:", err)
		}
	}(db)

	//2. Репозиторий
	orgRepo := postgres.NewOrganizationRepository(db)

	//3. Сервис
	orgService := service.NewOrganizationService(orgRepo)

	//4. Хендлер
	orgHandler := handler.NewOrganizationHandler(orgService)

	//5. Роутер
	r := gin.Default()
	api := r.Group("/api/v1")
	api.POST("organizations", orgHandler.CreateOrganization)

	// ЭТА СТРОКА ДОЛЖНА БЫТЬ ПОСЛЕ ВСЕХ ДРУГИХ МАРШРУТОВ!
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//6. Запуск
	log.Println("Server started on port 8080")
	err = r.Run(":8080")
	if err != nil {
		return
	}

	//7. Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

}
