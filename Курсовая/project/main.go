package main

import (
	"bufio"
	"log"
	"os"
	"strings"

	"example/project/backend"
	"example/project/backend/handler"
	"example/project/backend/middleware"
	"example/project/backend/repository"
	"example/project/backend/service"
	_ "example/project/docs" // автогенерированные файлы

	swaggerFiles "github.com/swaggo/files"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const filename = ".env"

// @title           Ветеринарная клиника API
// @version         1.0
// @description     REST API для системы управления ветеринарной клиникой
// @host            localhost:8080
// @BasePath        /
func main() {
	file, err := os.Open(filename)
	if err != nil {
		log.Println("Предупреждение: Файл .env не найден, используем системные переменные")
	} else {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				os.Setenv(parts[0], parts[1])
			}
		}
		file.Close() // Закрываем сразу после чтения
		log.Println("Конфигурация из .env загружена")
	}

	// (теперь ConnectDB увидит переменные, если они там используются)
	db, err := backend.ConnectDB()
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	log.Println("База подключена и миграция выполнена успешно!")

	// 2. Инициализация Репозиториев (Слой DB)
	userRepo := repository.NewUserRepository(db)
	petRepo := repository.NewPetRepository(db)
	doctorRepo := repository.NewDoctorRepository(db)
	appRepo := repository.NewAppointmentRepository(db)
	visitRepo := repository.NewVisitRepository(db)
	invRepo := repository.NewInventoryRepository(db)

	// 3. Инициализация Сервисов (Слой логики)
	authSrv := service.NewAuthService(userRepo, doctorRepo)
	petSrv := service.NewPetService(petRepo, userRepo)
	appSrv := service.NewAppointmentService(appRepo, petRepo, userRepo, doctorRepo)
	healthSrv := service.NewHealthJournalService(visitRepo, doctorRepo)
	doctorSrv := service.NewDoctorService(doctorRepo)
	anaSrv := service.NewAnalyticsService(visitRepo, doctorRepo, invRepo)
	invSrv := service.NewInventoryService(invRepo)
	dashSrv := service.NewDashboardService(userRepo, petRepo, appRepo, visitRepo)
	userSrv := service.NewUsersService(userRepo)
	searchSrv := service.NewSearchService(petRepo, userRepo)
	// 4. Инициализация Хендлеров (Слой API)
	authHandler := handler.NewAuthHandler(authSrv)
	petHandler := handler.NewPetHandler(petSrv)
	visitHandler := handler.NewVisitHandler(healthSrv)
	appHandler := handler.NewAppointmentHandler(appSrv)
	doctorHandler := handler.NewDoctorHandler(doctorSrv, appSrv)
	adminHandler := handler.NewAdminHandler(anaSrv, invSrv, userSrv, doctorSrv)
	weightHandler := handler.NewWeightHandler(petSrv)
	dashHandler := handler.NewDashboardHandler(dashSrv)
	userHandler := handler.NewUserHandler(userSrv)
	searchHandler := handler.NewSearchHandler(searchSrv)

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

	// 6. Роутинг
	api := r.Group("/api")
	{
		// --- ПУБЛИЧНЫЕ РОУТЫ ---
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/register", authHandler.Register)
		}

		// --- ЗАЩИЩЕННЫЕ РОУТЫ (Требуется любой валидный токен) ---
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware()) // <--- Включаем проверку токена
		{
			// Метод logout теперь тут, так как разлогиниться может только тот, кто вошел
			protected.POST("/auth/logout", authHandler.Logout)

			// PetController
			pets := protected.Group("/pets")
			{
				pets.GET("/owner/:ownerId", petHandler.GetByOwner)
				pets.POST("", petHandler.AddPet)
				pets.PUT("/:petId", petHandler.UpdatePet)
				pets.DELETE("/:petId", petHandler.DeletePet)
			}

			// VisitController
			visits := protected.Group("/visits")
			{
				// Можно добавить проверку RoleMiddleware("doctor", "admin") для POST,
				// так как клиент не должен сам себе писать медкарту
				visits.POST("", visitHandler.SaveVisitCard)
				visits.GET("/pet/:petId", visitHandler.GetJournal)
				visits.GET("/:id", visitHandler.GetById)
			}

			// AppointmentController
			apps := protected.Group("/appointments")
			{
				apps.POST("", appHandler.Create)
				apps.GET("/owner/:ownerId", appHandler.GetByOwner)
				apps.PUT("/:id/status", appHandler.UpdateStatus)
				apps.DELETE("/:id", appHandler.Cancel)
				apps.GET("/busy-slots", appHandler.GetBusySlots)
			}

			// DoctorController
			docs := protected.Group("/doctors")
			{
				docs.GET("", doctorHandler.GetBySpecialty)
				docs.GET("/:id/schedule", doctorHandler.GetSchedule)
			}

			// WeightController
			weight := protected.Group("/weight")
			{
				weight.GET("/pet/:petId", weightHandler.GetHistory)
				weight.POST("/pet/:petId", weightHandler.AddRecord)
			}

			users := protected.Group("/users")
			{
				users.PUT("/:id", userHandler.UpdateProfile)
				users.PUT("/:id/password", userHandler.ChangePassword)
			}

			// Dashboard (Обзор)
			protected.GET("/dashboard/:ownerId", dashHandler.GetData)
			protected.GET("/search", searchHandler.Search)
			// --- ТОЛЬКО ДЛЯ АДМИНИСТРАТОРОВ ---
			admin := protected.Group("/admin")
			admin.Use(middleware.RoleMiddleware("admin"))
			{
				admin.GET("/stats", adminHandler.GetStats)
				admin.GET("/revenue", adminHandler.GetRevenue)
				admin.GET("/users", adminHandler.GetAllUsers)

				admin.PUT("/users/:id/role", adminHandler.UpdateUserRole)

				admin.POST("/doctors", adminHandler.CreateDoctor)

				admin.GET("/services", adminHandler.GetAllServices)
				admin.DELETE("/services/:id", adminHandler.DeleteService)
				admin.POST("/services", adminHandler.CreateSrv)
				admin.PUT("/services/:id", adminHandler.UpdateService)

				admin.GET("/meds", adminHandler.GetMeds)
				admin.POST("/meds", adminHandler.CreateMed)
				admin.PUT("/meds/:id", adminHandler.UpdateMed)
				admin.DELETE("/meds/:id", adminHandler.DeleteMed)

				admin.DELETE("/users/:id", adminHandler.DeleteUser)
				admin.POST("/doctors/create-full", adminHandler.RegisterDoctorFull)
			}
		}

		// Swagger остается доступным всем для удобства разработки
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// Запуск сервера
	log.Println("Сервер запущен на :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Ошибка запуска сервера: ", err)
	}
}
