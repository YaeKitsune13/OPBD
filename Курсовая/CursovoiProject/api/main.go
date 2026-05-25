package main

import (
	"api/internal/handler"
	"api/internal/middleware"
	"api/internal/repository"
	"api/internal/service"

	_ "api/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	database "api/utils"

	"github.com/gin-gonic/gin"
)

// @title Veterinary Clinic API
// @version 1.0
// @description API server for Veterinary Clinic management system.
// @host localhost:3000
// @BasePath /api
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	db := database.InitDB()

	userRepo := repository.NewUserRepository(db)
	authSvc := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authSvc)

	petRepo := repository.NewPetRepository(db)
	petSvc := service.NewPetService(petRepo)
	petHandler := handler.NewPetHandler(petSvc)

	dashRepo := repository.NewDashboardRepository(db)
	dashSvc := service.NewDashboardService(dashRepo, petRepo)
	dashHandler := handler.NewDashboardHandler(dashSvc)

	appRepo := repository.NewAppointmentRepository(db)
	bookingSvc := service.NewBookingService(appRepo, petRepo)
	bookingHandler := handler.NewBookingHandler(bookingSvc)

	shopRepo := repository.NewShopRepository(db)
	shopSvc := service.NewShopService(shopRepo)
	shopHandler := handler.NewShopHandler(shopSvc)

	statsRepo := repository.NewStatsRepository(db)
	statsSvc := service.NewStatsService(statsRepo)
	statsHandler := handler.NewStatsHandler(statsSvc, petSvc)

	docRepo := repository.NewDoctorRepository(db)
	docSvc := service.NewDoctorService(docRepo)
	docHandler := handler.NewDoctorHandler(docSvc)

	userSvc := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc)

	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")
	{
		api.POST("/auth/register", authHandler.Register)
		api.POST("/auth/login", authHandler.Login)

		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/users/:id", userHandler.GetProfile)
			protected.PUT("/users/:id", userHandler.UpdateProfile)
			protected.PUT("/users/:id/password", userHandler.ChangePassword)

			protected.GET("/pets/owner/:userId", petHandler.GetByOwner)
			protected.POST("/pets", petHandler.Create)
			protected.PUT("/pets/:petId", petHandler.Update)
			protected.DELETE("/pets/:petId", petHandler.Delete)

			protected.GET("/dashboard/:userId", dashHandler.GetDashboard)

			protected.GET("/book/init/:userId", bookingHandler.GetInit)
			protected.GET("/appointments/busy-slots", bookingHandler.GetBusySlots)
			protected.POST("/appointments", bookingHandler.Create)
			protected.GET("/appointments/client/:userId", bookingHandler.GetHistory)

			protected.GET("/doctor/schedule", docHandler.GetSchedule)
			protected.PATCH("/appointments/:id/complete", docHandler.CompleteVisit)
			protected.PATCH("/appointments/:id/status", docHandler.UpdateStatus)
			protected.GET("/doctor/patients", docHandler.GetPatients)
			protected.GET("/doctor/patients/:id/history", docHandler.GetHistory)

			protected.GET("/medications", shopHandler.GetProducts)
			protected.GET("/cart/:userId", shopHandler.GetCart)
			protected.POST("/cart/:userId", shopHandler.AddToCart)
			protected.PUT("/cart/:itemId", shopHandler.UpdateCart)
			protected.DELETE("/cart/:itemId", shopHandler.DeleteCart)
			protected.POST("/orders", shopHandler.Checkout)
			protected.GET("/orders/:userId", shopHandler.GetOrders)

			protected.GET("/stats/pets/:userId", statsHandler.GetUserPets)
			protected.GET("/stats/weight/:petId", statsHandler.GetWeightData)
		}
	}

	r.Run(":3000")
}
