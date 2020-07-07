package controller

import (
	"math/rand"
	"net/http"
	"prototype2/domain"

	"github.com/gin-gonic/gin"
	"github.com/getsentry/sentry-go"
)

type userController struct {
	userService domain.UserService
}

type UserController interface {
	GetUsers(c *gin.Context)
	AddUser(c *gin.Context)
}

//NewUserController: constructor, dependency injection from user service 
func NewUserController(s domain.UserService) UserController {
	return &userController{
		userService: s,
	}
}

func (u *userController) GetUsers(c *gin.Context) {
	users, err := u.userService.FindAll()
	if err != nil {
		sentry.CaptureException(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while getting users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (u *userController) AddUser(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		sentry.CaptureException(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.ID = rand.Int63()

	if err1 := (u.userService.Validate(&user)); err1 != nil {
		sentry.CaptureException(err1)
		c.JSON(http.StatusBadRequest, gin.H{"error": err1.Error()})
		return 
	}

	if ageValidation := (u.userService.ValidateAge(&user)); ageValidation != true {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid DOB"})
		return
	}
	
	u.userService.Create(&user)
	c.JSON(http.StatusOK, user)
}