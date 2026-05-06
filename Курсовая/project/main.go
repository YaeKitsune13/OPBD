package main

import (
	"log"

	"example/project/backend"
	"example/project/backend/handler"
	"example/project/backend/repository"
	"example/project/backend/service"
	_ "example/project/docs" // автогенерированные файлы

	swaggerFiles "github.com/swaggo/files"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Ветеринарная клиника API
// @version         1.0
// @description     REST API для системы управления ветеринарной клиникой
// @host            localhost:8080
// @BasePath        /
func main() {
	// 1. Подключение к базе данных через твой пакет backend
	db, err := backend.ConnectDB()
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	log.Println("База подключена и миграция выполнена успешно!")

	// 2. Инициализация Репозиториев (Слой DB)
	ownerRepo := repository.NewOwnerRepository(db)
	petRepo := repository.NewPetRepository(db)
	doctorRepo := repository.NewDoctorRepository(db)
	appRepo := repository.NewAppointmentRepository(db)
	visitRepo := repository.NewVisitRepository(db)
	invRepo := repository.NewInventoryRepository(db)

	// 3. Инициализация Сервисов (Слой логики)
	authSrv := service.NewAuthService(ownerRepo, doctorRepo)
	petSrv := service.NewPetService(petRepo, ownerRepo)
	appSrv := service.NewAppointmentService(appRepo, petRepo, ownerRepo, doctorRepo)
	healthSrv := service.NewHealthJournalService(visitRepo, doctorRepo)
	doctorSrv := service.NewDoctorService(doctorRepo)
	anaSrv := service.NewAnalyticsService(visitRepo, doctorRepo, invRepo)
	invSrv := service.NewInventoryService(invRepo)
	dashSrv := service.NewDashboardService(ownerRepo, petRepo, appRepo, visitRepo)

	// 4. Инициализация Хендлеров (Слой API)
	authHandler := handler.NewAuthHandler(authSrv)
	petHandler := handler.NewPetHandler(petSrv)
	visitHandler := handler.NewVisitHandler(healthSrv)
	appHandler := handler.NewAppointmentHandler(appSrv)
	doctorHandler := handler.NewDoctorHandler(doctorSrv, appSrv)
	adminHandler := handler.NewAdminHandler(anaSrv, invSrv)
	weightHandler := handler.NewWeightHandler(petSrv)
	dashHandler := handler.NewDashboardHandler(dashSrv)

	// 5. Настройка Gin
	r := gin.Default()

	// Настройка CORS (чтобы Vue/Vite мог слать запросы)
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 6. Роутинг (согласно таблице 2.1)
	api := r.Group("/api")
	{
		// AuthController
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/register", authHandler.Register)
			auth.POST("/logout", authHandler.Logout)
		}

		// PetController
		pets := api.Group("/pets")
		{
			pets.GET("/owner/:ownerId", petHandler.GetByOwner)
			pets.POST("", petHandler.AddPet)
			pets.PUT("/:petId", petHandler.UpdatePet)
			pets.DELETE("/:petId", petHandler.DeletePet)
		}

		// VisitController
		visits := api.Group("/visits")
		{
			visits.POST("", visitHandler.SaveVisitCard)
			visits.GET("/pet/:petId", visitHandler.GetJournal)
			visits.GET("/:id", visitHandler.GetById)
		}

		// AppointmentController
		apps := api.Group("/appointments")
		{
			apps.POST("", appHandler.Create)
			apps.GET("/owner/:ownerId", appHandler.GetByOwner)
			apps.PUT("/:id/status", appHandler.UpdateStatus)
			apps.DELETE("/:id", appHandler.Cancel)
		}

		// DoctorController
		docs := api.Group("/doctors")
		{
			docs.GET("", doctorHandler.GetBySpecialty)
			docs.GET("/:id/schedule", doctorHandler.GetSchedule)
		}

		// AdminController
		admin := api.Group("/admin")
		{
			admin.GET("/stats", adminHandler.GetStats)
			admin.GET("/revenue", adminHandler.GetRevenue)
			admin.DELETE("/services/:id", adminHandler.DeleteService)
			admin.POST("/meds", adminHandler.CreateMed)
		}

		// WeightController
		weight := api.Group("/weight")
		{
			weight.GET("/pet/:petId", weightHandler.GetHistory)
			weight.POST("/pet/:petId", weightHandler.AddRecord)
		}

		// Dashboard (Обзор)
		api.GET("/dashboard/:ownerId", dashHandler.GetData)
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// Запуск сервера
	log.Println("Сервер запущен на :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Ошибка запуска сервера: ", err)
	}
}
