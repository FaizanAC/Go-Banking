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
	bankGroup := r.Group("/bank")
	{
		bankGroup.POST("/new-account", middleware.AuthorizeRequest, s.handleNewAccount)
		bankGroup.GET("/account/:id", middleware.AuthorizeRequest, s.handleGetAccount)
		bankGroup.POST("/deposit", middleware.AuthorizeRequest, s.handleDeposit)
		bankGroup.POST("/withdraw", middleware.AuthorizeRequest, s.handleWithdraw)
		bankGroup.GET("/activity-feed", middleware.AuthorizeRequest, s.handleActivityFeed)

		transferGroup := bankGroup.Group("/transfer")
		{
			transferGroup.POST("/send", middleware.AuthorizeRequest, s.handleSendTransfer)
			transferGroup.POST("/accept", middleware.AuthorizeRequest, s.handleAcceptTransfer)
		}
	}

	fmt.Println("Server is running on port", s.Port)
	r.Run()
}

func NewServer(db *gorm.DB, port string) *Server {
	return &Server{
		db:   db,
		Port: port,
	}
}
