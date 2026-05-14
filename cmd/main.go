package main

import (
	"apartment-manager-backend/config"
	"apartment-manager-backend/internal/application/service"
	"apartment-manager-backend/internal/infrastructure/database"
	"net/http"

	//"apartment-manager-backend/internal/infrastructure/repository/postgres"
	"apartment-manager-backend/internal/infrastructure/repository/redis"
	"apartment-manager-backend/internal/infrastructure/sms"
	"apartment-manager-backend/internal/presentation/controller"
	"apartment-manager-backend/internal/presentation/routes"
	"log"
)

func main() {

	// --- CONFIG ---
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// --- POSTGRES ---
	db := database.NewPostgresDB(cfg.Postgres.DSN())
	sqlDB, _ := db.DB()
	err = sqlDB.Ping()
	if err != nil {
		log.Fatal("database not reachable")
	}

	// --- REDIS ---
	redisClient := database.NewRedisClient(cfg.Redis)
	err = redisClient.Set(database.Ctx, "test", "ok", 0).Err()
	if err != nil {
		log.Fatal(err)
	}

	// --- REPOSITORIES ---
	//userRepo := postgres.NewUserRepository(db)
	otpRepo := redis.NewOTPRepository(redisClient)

	// ------ SMS ------
	smsService := sms.NewFileSMS("otp_log.txt")

	// --- SERVICES ---
	SendOtpService := service.NewSendOtpService(otpRepo, smsService)

	// --- CONTROLLER ---
	SendOtpController := controller.NewAuthHandler(SendOtpService)

	// --- ROUTE ---
	controllers := &routes.Controllers{
		SendOTP: SendOtpController,
	}

	r := routes.SetUpRouter(controllers)

	log.Println("server running on :8080")
	http.ListenAndServe(":8080", r)
}
