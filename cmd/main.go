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
	"fmt"
	"log"
	"net/http"
)

const TEST_AND_USE_HASHER bool = true           //TODO : Remove This is Release Mode
const test_password_for_debug string = "123456" //TODO : Remove This is Release Mode

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
	unitRepo := postgres.NewUnitRepository(db)
	inviteCodeRepo := postgres.NewInviteCodeRepository(db)
	tagRepo := postgres.NewTagRepository(db)
	announcementRepo := postgres.NewAnnouncementRepository(db)
	commentRepo := postgres.NewCommentRepository(db)

	// ------ SMS ------
	smsService := sms.NewFileSMS("otp_log.txt")

	// ----- Hasher -----
	passwordHasher := hasher.NewBcryptHasher()

	//TODO : Remove This is Release Mode
	if TEST_AND_USE_HASHER {
		s, err := passwordHasher.Hash(test_password_for_debug)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("test-hash-password:%s*\n",s)
	}

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
	userService := service.NewUserService(userRepo, passwordHasher)
	apartmentSrv := service.NewApartmentService(apartmentRepo)
	ticketSrv := service.NewTicketService(ticketRepo)
	inviteCodeService := service.NewInviteCodeService(inviteCodeRepo, unitRepo, apartmentRepo)
	tagService := service.NewTagService(tagRepo)
	announcementService := service.NewAnnouncementService(announcementRepo, tagRepo)
	commentService := service.NewCommentService(commentRepo, ticketRepo)

	// --- CONTROLLER ---
	authController := controller.NewAuthController(authService)
	userController := controller.NewUserController(userService)
	apartmentController := controller.NewApartmentController(apartmentSrv)
	tickerController := controller.NewTicketController(ticketSrv, commentService)
	inviteCodeController := controller.NewInviteCodeController(inviteCodeService)
	tagController := controller.NewTagController(tagService)
	announcementController := controller.NewAnnouncementController(announcementService)
	commentController := controller.NewCommentController(commentService)

	// --- ROUTE ---
	controllers := &routes.Controllers{
		Auth:         authController,
		User:         userController,
		Apartment:    apartmentController,
		Ticket:       tickerController,
		InviteCode:   inviteCodeController,
		Tag:          tagController,
		Announcement: announcementController,
		Comment:      commentController,
	}

	r := routes.SetUpRouter(controllers, jwtRepo)

	log.Println("server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
