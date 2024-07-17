package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	_ "github.com/golang-jwt/jwt/v4"
	"log"
	"nedorezov/pkg/db"
	"nedorezov/pkg/domain/account"
	"net/http"
	"strconv"
	"time"
)

var jwtKey = []byte("secret_key")

func (h *Handler) Registration(c *gin.Context) {
	var acc account.Registration
	if err := c.BindJSON(&acc); err != nil {
		c.JSON(400, gin.H{"error bind Registration account": err.Error()})
		log.Printf("error bind Registration account %v", err.Error())
		return
	}
	userId, _ := c.Cookie("id")
	result, err := h.service.Registration(c.Request.Context(), userId, acc)
	if err != nil {
		switch err.Error() {
		case db.ErrValidate.Error():
			c.JSON(400, gin.H{"error bind Registration account": err.Error()})
			log.Printf("error bind Registration account %v", err.Error())
			break
		case db.ErrAuthorize.Error():
			c.JSON(403, gin.H{"Already Authorized": err.Error()})
			log.Println("Already Auth on Register")
			break
		case db.ErrDuplicate.Error():
			c.JSON(409, "Email already exist")
			log.Printf("Register email failed %v", err.Error())
			break
		default:
			c.JSON(500, gin.H{"error bind Registration account": err.Error()})
			log.Printf("error service Registration account %v", err.Error())
		}
		return
	}

	c.JSON(201, result)
	log.Printf("Success registration: %v", result)
}

func (h *Handler) Login(c *gin.Context) {
	var acc account.Login
	if err := c.BindJSON(&acc); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	log.Printf("Login: %v", acc)

	id, err := h.service.Login(c.Request.Context(), acc)
	if err != nil {
		switch err.Error() {
		case db.ErrNotExist.Error():
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Account does not exist"})
			log.Printf("Account does not exist")
			break
		default:
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Printf("error service Login error %v", err.Error())
		}
		return
	}

	token, err := createToken(strconv.FormatInt(id, 10))
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Error creating token"})
		return
	}
	sameSite := http.SameSiteNoneMode // Важно установить SameSite=None с Secure=true

	//c.SetCookie("token", token, 3600*72, "/", "", true, true)
	//c.SetSameSite(sameSite)
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 72),
		Path:     "/",
		Domain:   "",
		SameSite: sameSite,
		Secure:   true,
		HttpOnly: true,
	})

	c.IndentedJSON(http.StatusOK, gin.H{"id": id})
}

func (h *Handler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("token")

		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		token, err := parseToken(tokenString)
		if err != nil || !token.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("userId", claims["id"])

		c.Next()
	}
}

func createToken(id string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	})

	return token.SignedString(jwtKey)
}

func parseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil

	})
}
