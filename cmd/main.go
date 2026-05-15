package main

import (
	"apartment-manager-backend/config"
	"apartment-manager-backend/internal/application/service"
	"apartment-manager-backend/internal/infrastructure/database"
	"apartment-manager-backend/internal/infrastructure/jwt"
	"apartment-manager-backend/internal/infrastructure/repository/postgres"
	"apartment-manager-backend/pkg/hasher"
	"net/http"

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
	userRepo := postgres.NewUserRepository(db)
	otpRepo := redis.NewOTPRepository(redisClient)
	refreshRepo := redis.NewRefreshTokenRepository(redisClient)
	jwtRepo := jwt.NewTokenService(cfg.JWT)

	// ------ SMS ------
	smsService := sms.NewFileSMS("otp_log.txt")

	// ----- Hasher -----
	passwordHasher := hasher.NewBcryptHasher()

	// --- SERVICES ---
	SendOtpService := service.NewSendOtpService(otpRepo, smsService)
	VerifyOtpService := service.NewVerifyOTPService(otpRepo, userRepo, refreshRepo, jwtRepo)
	RegisterService := service.NewRegisterService(userRepo, otpRepo, refreshRepo, jwtRepo, passwordHasher)
	LoginService := service.NewLoginService(userRepo, refreshRepo, jwtRepo, passwordHasher)
	RefreshService := service.NewRefreshTokenService(refreshRepo, jwtRepo)
	LogoutService := service.NewLogoutService(refreshRepo)
	profileService := service.NewProfileService(userRepo)

	// --- CONTROLLER ---
	SendOtpController := controller.NewAuthHandler(SendOtpService)
	VerifyOtpController := controller.NewVerifyController(VerifyOtpService)
	registerController := controller.NewRegisterController(RegisterService)
	loginController := controller.NewLoginController(LoginService)
	refreshController := controller.NewRefreshController(RefreshService)
	logoutController := controller.NewLogoutController(LogoutService)
	profileController := controller.NewProfileController(profileService)

	// --- ROUTE ---
	controllers := &routes.Controllers{
		SendOTP:   SendOtpController,
		VerifyOTP: VerifyOtpController,
		Register:  registerController,
		Login:     loginController,
		Refresh:   refreshController,
		Logout:    logoutController,
		Profile:   profileController,
	}

	r := routes.SetUpRouter(controllers, jwtRepo)

	log.Println("server running on :8080")
	http.ListenAndServe(":8080", r)
}
