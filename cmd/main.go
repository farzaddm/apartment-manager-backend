package main

import (
	"apartment-manager-backend/config"
	"apartment-manager-backend/initsamples"
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

	// --- MIGRATIONS ---
	//goose.SetBaseFS(migrations.FS)
	//if err := goose.SetDialect("postgres"); err != nil {
	//	log.Fatalf("failed to set goose dialect: %v", err)
	//}
	//
	//log.Println("Running database migrations...")
	//if err := goose.Up(sqlDB, "."); err != nil {
	//	log.Fatalf("failed to run migrations: %v", err)
	//}
	//log.Println("Migrations completed successfully!")

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
	ruleRepo := postgres.NewRuleRepository(db)
	pollRepo := postgres.NewPollRepository(db)

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
		unitRepo,
	)
	userService := service.NewUserService(userRepo, passwordHasher)
	apartmentSrv := service.NewApartmentService(apartmentRepo)
	ticketSrv := service.NewTicketService(ticketRepo, tagRepo)
	inviteCodeService := service.NewInviteCodeService(inviteCodeRepo, unitRepo, apartmentRepo)
	tagService := service.NewTagService(tagRepo)
	announcementService := service.NewAnnouncementService(announcementRepo, tagRepo)
	commentService := service.NewCommentService(commentRepo, ticketRepo)
	ruleService := service.NewRuleService(ruleRepo)
	pollService := service.NewPollService(pollRepo)
	unitService := service.NewUnitService(unitRepo, userRepo)

	// --- CONTROLLER ---
	authController := controller.NewAuthController(authService)
	userController := controller.NewUserController(userService)
	apartmentController := controller.NewApartmentController(apartmentSrv)
	tickerController := controller.NewTicketController(ticketSrv, commentService)
	inviteCodeController := controller.NewInviteCodeController(inviteCodeService)
	tagController := controller.NewTagController(tagService)
	announcementController := controller.NewAnnouncementController(announcementService)
	commentController := controller.NewCommentController(commentService)
	ruleController := controller.NewRuleController(ruleService)
	pollController := controller.NewPollController(pollService)
	unitController := controller.NewUnitController(unitService)

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
		Rule:         ruleController,
		Poll:         pollController,
		Unit:         unitController,
	}

	r := routes.SetUpRouter(controllers, jwtRepo)

	//TODO : THIS IMP IS JUST FOR TEST BUT IN RELEASE VER WE NEED TO IN MORE PROPER IMP
	initsamples.CreateOrOverWriteManagersAndAdminAndResident(db, passwordHasher)

	log.Printf("server running on :%s\n", cfg.App.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cfg.App.Port), r))
}
