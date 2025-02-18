package server

import (
	"fmt"

	"github.com/FaizanAC/Go-Banking/internal/middleware"
	"github.com/FaizanAC/Go-Banking/internal/server/handlers"
	"github.com/FaizanAC/Go-Banking/internal/server/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	db   *gorm.DB
	Port string
}

func (s *Server) SetupRouter() *gin.Engine {
	r := gin.Default()

	healthHandler := handlers.HealthHandler{}

	userService := services.NewUserService(s.db)
	userHandler := handlers.NewUserHandler(userService)

	loginService := services.NewLoginService(s.db)
	loginHandler := handlers.NewLoginHandler(loginService)

	bankService := services.NewBankService(s.db)
	bankHandler := handlers.NewBankHandler(bankService)

	// Health
	r.GET("/ping", healthHandler.HandlePing)

	// User
	r.GET("/user/:id", middleware.AuthorizeRequest, userHandler.HandleGetUser)
	r.POST("/user", userHandler.HandleUserCreation)

	// Login
	r.POST("/login", loginHandler.HandleLogin)

	// Bank
	bankGroup := r.Group("/bank")
	{
		bankGroup.POST("/new-account", middleware.AuthorizeRequest, bankHandler.HandleNewAccount)
		bankGroup.GET("/accounts", middleware.AuthorizeRequest, bankHandler.HandleGetAccounts)
		bankGroup.POST("/deposit", middleware.AuthorizeRequest, bankHandler.HandleDeposit)
		bankGroup.POST("/withdraw", middleware.AuthorizeRequest, bankHandler.HandleWithdraw)
		bankGroup.GET("/activity-feed", middleware.AuthorizeRequest, bankHandler.HandleActivityFeed)

		transferGroup := bankGroup.Group("/transfer")
		{
			transferGroup.POST("/send", middleware.AuthorizeRequest, bankHandler.HandleSendTransfer)
			transferGroup.POST("/accept", middleware.AuthorizeRequest, bankHandler.HandleAcceptTransfer)
		}
	}

	return r
}

func (s *Server) Start() {
	r := s.SetupRouter()

	fmt.Println("Server is running on port", s.Port)
	r.Run()
}

func NewServer(db *gorm.DB, port string) *Server {
	return &Server{
		db:   db,
		Port: port,
	}
}
