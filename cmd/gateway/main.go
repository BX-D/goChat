package main

import (
	"fmt"
	"log"

	"github.com/boxuanduan/gochat/config"
	"github.com/boxuanduan/gochat/internal/chat"
	"github.com/boxuanduan/gochat/internal/gateway"
	"github.com/boxuanduan/gochat/internal/gateway/handler"
	"github.com/boxuanduan/gochat/internal/repo/mysql"
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

	if err != nil {
		log.Fatalf("init mysql: %v", err)
	}

	// Create repo instances
	repo := mysql.NewUserRepo(init)

	// Create handler instances
	userService := chat.NewUserService(repo, cfg.JWT)

	// Setup Gin router
	r := gin.Default()
	userHandler := handler.NewUserHandler(userService)

	// Register routes
	userHandler.RegisterRoutes(r)

	// Create and run hub
	hub := gateway.NewHub()
	go hub.Run()

	// Register WebSocket route
	wsHandler := handler.NewWSHandler(hub, cfg.JWT)
	r.GET("/ws", wsHandler.ServeWS)

	// Start server
	if err := r.Run(fmt.Sprintf(":%d", cfg.Gateway.HTTPPort)); err != nil {
		log.Fatalf("run server: %v", err)
	}
}
