package main

import (
	"context"
	"log"
	"qzone-history/internal/delivery/app"
	"qzone-history/internal/infrastructure/config"
	"qzone-history/internal/infrastructure/persistence"
	"qzone-history/internal/infrastructure/qzone_api"
	"qzone-history/internal/usecase"
	"qzone-history/pkg/database"
	"qzone-history/pkg/database/sqlite"
)

func main() {
	ctx := context.Background()

	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	var db database.Database
	switch cfg.Database.Type {
	case "sqlite":
		db = sqlite.NewSQLiteDB()
	default:
		log.Fatalf("unsupported database type: %s", cfg.Database.Type)
	}
	dbConfig := &database.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
	}
	err = db.Connect(dbConfig)

	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("failed to auto migrate: %w", err)
	}
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	//
	qzoneAPIClient := qzone_api.NewQzoneAPIClient(cfg)

	// 初始化仓库
	userRepo := persistence.NewUserRepository(db)
	momentRepo := persistence.NewMomentRepository(db)
	activityRepo := persistence.NewActivityRepository(db)
	boardMessageRepo := persistence.NewBoardMessageRepository(db)
	friendRepo := persistence.NewFriendRepository(db)

	// 初始化用例
	authUseCase := usecase.NewAuthUseCase(qzoneAPIClient, userRepo)
	momentUseCase := usecase.NewMomentUseCase(momentRepo)
	boardMessageUseCase := usecase.NewBoardMessageUseCase(boardMessageRepo)
	friendUseCase := usecase.NewFriendUseCase(friendRepo)
	exportUseCase := usecase.NewExportUseCase(momentRepo, boardMessageRepo, friendRepo)
	activityUseCase := usecase.NewActivityUseCase(qzoneAPIClient, activityRepo)
	reconstructionUseCase := usecase.NewReconstructionUseCase(activityRepo, momentRepo, boardMessageRepo)

	// 启动应用程序
	newApp := app.NewApp(authUseCase, momentUseCase, boardMessageUseCase, friendUseCase, exportUseCase, activityUseCase, reconstructionUseCase)
	if err := newApp.Run(ctx); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}
