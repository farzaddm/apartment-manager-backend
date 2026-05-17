package main

import (
	"apartment-manager-backend/config"
	"apartment-manager-backend/internal/application/service"
	"apartment-manager-backend/internal/infrastructure/database"
	"apartment-manager-backend/internal/infrastructure/jwt"
	"apartment-manager-backend/internal/infrastructure/repository/postgres"
	"apartment-manager-backend/internal/infrastructure/repository/redis"
	"apartment-manager-backend/internal/infrastructure/sms"
	"apartment-manager-backend/internal/presentation/controller"
	"apartment-manager-backend/internal/presentation/routes"
	"apartment-manager-backend/pkg/hasher"
	"log"
	"net/http"
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
	apartmentRepo := postgres.NewApartmentRepository(db)
	ticketRepo := postgres.NewTicketRepository(db)

	// ------ SMS ------
	smsService := sms.NewFileSMS("otp_log.txt")

	// ----- Hasher -----
	passwordHasher := hasher.NewBcryptHasher()

	// --- SERVICES ---
	authService := service.NewAuthService(
		userRepo,
		otpRepo,
		refreshRepo,
		jwtRepo,
		smsService,
		passwordHasher,
		cfg.JWT.RefreshExpireDays,
	)
	profileService := service.NewProfileService(userRepo)
	apartmentSrv := service.NewApartmentService(apartmentRepo)
	ticketSrv := service.NewTicketService(ticketRepo)

	// --- CONTROLLER ---
	authController := controller.NewAuthController(authService)
	profileController := controller.NewProfileController(profileService)
	apartmentController := controller.NewApartmentController(apartmentSrv)
	tickerController := controller.NewTicketController(ticketSrv)

	// --- ROUTE ---
	controllers := &routes.Controllers{
		Auth:      authController,
		Profile:   profileController,
		Apartment: apartmentController,
		Ticket:    tickerController,
	}

	r := routes.SetUpRouter(controllers, jwtRepo)

	log.Println("server running on :8080")
	http.ListenAndServe(":8080", r)
}
