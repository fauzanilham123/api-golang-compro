package controllers

import (
	"api-golang-compro/models"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

type LoginInput struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

type RegisterInput struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
    Email    string `json:"email" binding:"required"`
}

func ParseTokenExp(tokenString string) (int64, error) {
    token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
    if err != nil {
        return 0, err
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return 0, fmt.Errorf("failed to parse token claims")
    }

    exp, ok := claims["exp"].(float64)
    if !ok {
        return 0, fmt.Errorf("failed to parse token expiration")
    }

    return int64(exp), nil
}


func Login(c *gin.Context) {
    db := c.MustGet("db").(*gorm.DB)
    var input LoginInput

    if err := c.ShouldBind(&input); err != nil {
        SendError(c, "error", err.Error())
        return
    }

    u := models.User{}
    u.Username = input.Username
    u.Password = input.Password


    token, err := models.LoginCheck(u.Username, u.Password, db)
    if err != nil {
        fmt.Println(err)
        SendError(c, "username or password is incorrect.", err.Error())
        return
    }

    // Fetch the user from the database to get the email
    if err := db.Where("username = ?", u.Username).First(&u).Error; err != nil {
        SendError(c, "failed to fetch user details.", err.Error())
        return
    }

    user := map[string]string{
        "username": u.Username,
        "email":    u.Email,
    }

    tokenExp, err := ParseTokenExp(token)
    if err != nil {
        SendError(c, "failed to parse token expiration.", err.Error())
        return
    }

    // Konversi waktu kedaluwarsa (exp) ke format yang lebih mudah dibaca
    tokenExpReadable := time.Unix(tokenExp, 0).Format("2006-01-02 15:04:05 MST")

    // Send the user details along with the token in the response
    response := gin.H{
        "user":  user,
        "token": token,
        "token_exp": tokenExpReadable,
    }

    // Simpan id_user dalam variabel
    idUser := u.ID
    
	SendResponse(c, response, "login success")
    activityMessage := "Login" 
    activitylogin(c,activityMessage,idUser)
}


func Register(c *gin.Context) {
    db := c.MustGet("db").(*gorm.DB)
    var input RegisterInput

    if err := c.ShouldBind(&input); err != nil {
        SendError(c, "error", err.Error())
        return
    }

    // Check if email already exists in the database
    var existingUser models.User
    if err := db.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
        SendError(c, "Email already registered", "Email already registered")
        return
    }

    u := models.User{}
    u.Username = input.Username
    u.Email = input.Email
    u.Password = input.Password

    savedUser, err := u.SaveUser(db)

    if err != nil {
        SendError(c, "error", err.Error())
        return
    }

    user := map[string]string{
        "username": input.Username,
        "email":    input.Email,
    }
    SendResponse(c, user, "registration success")
    activityMessage := "Register" 
    activitylogin(c,activityMessage,savedUser.ID)
}

