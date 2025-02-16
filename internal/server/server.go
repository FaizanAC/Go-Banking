package server

import (
	"fmt"

	"github.com/FaizanAC/Go-Banking/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	db   *gorm.DB
	Port string
}

func (s *Server) Start() {
	r := gin.Default()

	// Health
	r.GET("/ping", s.handleHealthCheck)

	// User
	r.GET("/user/:id", middleware.AuthorizeRequest, s.handleGetUser)
	r.POST("/user", s.handleUserCreation)

	// Login
	r.POST("/login", s.handleLogin)

	// Bank
	r.POST("/new-account", middleware.AuthorizeRequest, s.handleNewAccount)
	r.GET("/account/:id", middleware.AuthorizeRequest, s.handleGetAccount)
	r.POST("/deposit", middleware.AuthorizeRequest, s.handleDeposit)
	r.POST("/withdraw", middleware.AuthorizeRequest, s.handleWithdraw)
	r.POST("/transfer", middleware.AuthorizeRequest, s.handleTransfer)

	fmt.Println("Server is running on port", s.Port)
	r.Run()
}

func NewServer(db *gorm.DB, port string) *Server {
	return &Server{
		db:   db,
		Port: port,
	}
}
