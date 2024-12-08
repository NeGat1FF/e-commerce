package main

import (
	"github.com/NeGat1FF/e-commerce/user-service/internal/config"
	"github.com/NeGat1FF/e-commerce/user-service/internal/db"
	"github.com/NeGat1FF/e-commerce/user-service/internal/handlers"
	"github.com/NeGat1FF/e-commerce/user-service/internal/repository"
	"github.com/NeGat1FF/e-commerce/user-service/internal/service"
	"github.com/NeGat1FF/e-commerce/user-service/pkg/logger"
	"github.com/NeGat1FF/e-commerce/user-service/proto"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	godotenv.Load()
	logger.Init("info")

	config := config.LoadConfig()

	db, err := db.InitDB("postgres://user:password@localhost:5432/e-commerce")
	if err != nil {
		logger.Logger.Fatal("Failed to connect to the database", zap.Error(err))
	}

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient(config.NOTIFICATION_URL, opts...)
	if err != nil {
		logger.Logger.Fatal("Failed to connect to the notification service", zap.Error(err))
	}

	notifcationService := proto.NewMailServiceClient(conn)

	repo := repository.NewUserRepository(db)
	service := service.NewUserService(repo, notifcationService, config.JWTSecret, config.Addr, config.Port)
	handler := handlers.NewUserHandler(service)

	ginServer := gin.Default()

	route := ginServer.Group("/api/v1/user")

	route.POST("/register", handler.Register)
	route.POST("/login", handler.Login)
	route.POST("/refresh_token", handler.RefreshToken)
	route.POST("/resend_verification_email", handler.ResendVerificationEmail)
	route.POST("/verify_email", handler.VerifyEmail)
	route.POST("/forgot_password", handler.ForgotPassword)
	route.POST("/reset_password", handler.ResetPassword)

	route.PATCH("/update", handler.UpdateUser)
	route.DELETE("/delete", handler.DeleteUser)

	err = ginServer.Run(":8050")
	if err != nil {
		logger.Logger.Fatal("Failed to start the server", zap.Error(err))
	}
}
