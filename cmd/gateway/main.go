package main

import (
	"fmt"
	"log"

	"github.com/boxuanduan/gochat/config"
	"github.com/boxuanduan/gochat/internal/chat"
	"github.com/boxuanduan/gochat/internal/gateway"
	"github.com/boxuanduan/gochat/internal/gateway/handler"
	"github.com/boxuanduan/gochat/internal/repo/mysql"
	"github.com/boxuanduan/gochat/internal/repo/redis"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	fmt.Printf("Gateway starting on :%d\n", cfg.Gateway.HTTPPort)
	fmt.Printf("MySQL DSN: %s\n", cfg.MySQL.DSN())
	fmt.Printf("RPC framework: %s\n", cfg.RPC.Framework)

	// Init mysql repo
	init, err := mysql.Init(&cfg.MySQL)
	// Init redis repo
	redisClient := redis.Init(&cfg.Redis)

	if err != nil {
		log.Fatalf("init mysql: %v", err)
	}

	// Create repo instances
	userRepo := mysql.NewUserRepo(init)
	msgRepo := mysql.NewMessageRepo(init)
	seqRepo := redis.NewSeqRepo(redisClient)
	groupRepo := mysql.NewGroupRepo(init)

	// Create handler instances
	userService := chat.NewUserService(userRepo, cfg.JWT)
	groupService := chat.NewGroupService(groupRepo)

	// Setup Gin router
	r := gin.Default()
	userHandler := handler.NewUserHandler(userService)
	msgSvc := chat.NewMessageService(msgRepo, seqRepo)
	msgHandler := handler.NewMsgHandler(msgSvc)
	groupHandler := handler.NewGroupHandler(groupService)

	// Register routes
	userHandler.RegisterRoutes(r)

	// Create and run hub
	hub := gateway.NewHub(msgSvc, groupService)
	go hub.Run()

	// Register WebSocket route
	wsHandler := handler.NewWSHandler(hub, cfg.JWT)
	r.GET("/ws", wsHandler.ServeWS)
	r.GET("/api/messages", msgHandler.PullMessages)
	r.POST("/api/messages/ack", msgHandler.Ack)

	// Register group routes
	r.POST("/api/group/create", groupHandler.CreateGroup)
	r.POST("/api/group/join", groupHandler.JoinGroup)
	r.GET("/api/group/:id/members", groupHandler.GetMembers)

	// Start server
	if err := r.Run(fmt.Sprintf(":%d", cfg.Gateway.HTTPPort)); err != nil {
		log.Fatalf("run server: %v", err)
	}
}
