package Handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"main/Authentication"
	"main/Database"
	"net/http"
)

type signupRequestBody struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name" binding:"required"`
}

type loginRequestBody struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (s *Server) HandleSignup(c *gin.Context) {
	var reqData signupRequestBody

	// unmarshal json
	err := c.BindJSON(&reqData)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "can not unmarshal json"})
		return
	}

	//Validate user data
	if mes, err := s.Validation.ValidateData(reqData.Email, reqData.Name, reqData.Password); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Validation error", "data": mes})
		return
	}

	// check user duplicate
	result, err := s.Database.CheckUserDuplicateByEmail(reqData.Email)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "something went wrong"})
		return
	}
	if result {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "user already exist"})
		return
	}

	// Create new user in database
	user := &Database.User{
		Name:     reqData.Name,
		Email:    reqData.Email,
		Password: reqData.Password,
	}
	if err := s.Database.CreateNewUser(user); err != nil {
		s.Logger.Errorf("can not create new user error:%v", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "something went wrong "})
		return
	}
	c.IndentedJSON(http.StatusCreated, gin.H{"message": "User successfully created"})
	return
}
func (s *Server) HandleLogin(c *gin.Context) {
	var reqData loginRequestBody

	// unmarshal json
	err := c.BindJSON(&reqData)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "can not unmarshal json"})
		return
	}

	// Authenticate of user
	err = s.Authentication.AuthenticateUserWithCredentials(Authentication.Credentials{
		Email:    reqData.Email,
		Password: reqData.Password,
	})
	if err != nil {
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": err.Error()})
		return
	}

	// Create the JWT token
	token, err := s.Authentication.GenerateJwtToken(reqData.Email)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "can not Create token"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"access_token": *token})
	return

}
func (s *Server) HandleCheckLogin(c *gin.Context) {
	token := c.Request.Header.Get("access_token")
	//cookie, err := c.Request.Cookie("access_token")
	//if err != nil {
	//	return
	//}
	//fmt.Println(cookie.Value)
	if token == "" {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	user, err := s.Authentication.CheckToken(token)
	if err != nil {
		if err.Error() == "something went wrong" {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err})
			fmt.Println("sg")
			return
		}
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": err})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": fmt.Sprintf("wellcome %v", user.Email)})

}
