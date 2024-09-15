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
	log.Println("加载配置文件...")
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}
	log.Printf("配置加载成功: %+v", cfg)

	// 初始化数据库
	log.Println("初始化数据库...")
	var db database.Database
	switch cfg.Database.Type {
	case "sqlite":
		db = sqlite.NewSQLiteDB()
	default:
		log.Fatalf("不支持的数据库类型: %s", cfg.Database.Type)
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
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	defer func() {
		log.Println("关闭数据库连接...")
		db.Close()
	}()

	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("数据库迁移失败: %w", err)
	}
	log.Println("数据库初始化完成")

	// 初始化 Qzone API 客户端
	log.Println("初始化 Qzone API 客户端...")
	qzoneAPIClient := qzone_api.NewQzoneAPIClient(cfg)
	log.Println("Qzone API 客户端初始化成功")

	// 初始化仓库
	log.Println("初始化仓库...")
	userRepo := persistence.NewUserRepository(db)
	momentRepo := persistence.NewMomentRepository(db)
	activityRepo := persistence.NewActivityRepository(db)
	boardMessageRepo := persistence.NewBoardMessageRepository(db)
	friendRepo := persistence.NewFriendRepository(db)
	log.Println("仓库初始化完成")

	// 初始化用例
	log.Println("初始化用例...")
	authUseCase := usecase.NewAuthUseCase(qzoneAPIClient, userRepo)
	momentUseCase := usecase.NewMomentUseCase(momentRepo)
	boardMessageUseCase := usecase.NewBoardMessageUseCase(boardMessageRepo)
	friendUseCase := usecase.NewFriendUseCase(friendRepo)
	exportUseCase := usecase.NewExportUseCase(momentRepo, boardMessageRepo, friendRepo)
	activityUseCase := usecase.NewActivityUseCase(qzoneAPIClient, activityRepo)
	reconstructionUseCase := usecase.NewReconstructionUseCase(activityRepo, momentRepo, boardMessageRepo)
	log.Println("用例初始化完成")

	// 启动应用程序
	log.Println("启动应用程序...")
	newApp := app.NewApp(authUseCase, momentUseCase, boardMessageUseCase, friendUseCase, exportUseCase, activityUseCase, reconstructionUseCase)
	if err := newApp.Run(ctx); err != nil {
		log.Fatalf("应用程序运行错误: %v", err)
	}
	log.Println("应用程序成功完成")
}
