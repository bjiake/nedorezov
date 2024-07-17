package handler

import (
	"nedorezov/pkg/db"
	"nedorezov/pkg/domain/account"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (h *Handler) GetBalance(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(401, gin.H{"error": "Authorization required"})
		log.Errorf("Auth error: Authorization required")
		return
	}
	id, ok := userId.(string)
	if !ok {
		c.JSON(401, gin.H{"error": "Invalid user ID"})
		log.Errorf("Invalid user ID")
		return
	}

	ch := make(chan struct{})

	go func() {
		result, err := h.service.GetBalance(c.Request.Context(), id)
		if err != nil {
			switch err.Error() {
			case db.ErrNotExist.Error():
				c.JSON(404, gin.H{"error": err.Error()})
				log.Errorf("ID:%v\terror:%v", id, err.Error())
				break
			default:
				c.JSON(500, gin.H{"error": err.Error()})
				log.Error(err.Error())
			}
		} else {
			c.JSON(200, gin.H{"data": result})
			log.Printf("Get Balance ID:%v %v", id, *result)
		}
		ch <- struct{}{}
	}()

	<-ch
	return
}

func (h *Handler) WithDraw(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(401, gin.H{"error": "Authorization required"})
		log.Errorf("Auth error: Authorization required")
		return
	}
	id, ok := userId.(string)
	if !ok {
		c.JSON(401, gin.H{"error": "Invalid user ID"})
		log.Errorf("Invalid user ID")
		return
	}

	var response account.ResponseBalance
	if err := c.BindJSON(&response); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		log.Errorf("ID:%v\tFail bind balance:%v", id, err.Error())
		return
	}

	ch := make(chan struct{})

	go func() {
		result, err := h.service.WithDraw(c.Request.Context(), id, response.Balance)
		if err != nil {
			switch err.Error() {
			case db.ErrNotValidAmount.Error():
				c.JSON(400, gin.H{"error": err.Error()})
				log.Errorf("ID:%v\tWithdraw:%v", id, err.Error())
				break
			case db.ErrUpdateFailed.Error():
				c.JSON(404, gin.H{"error": err.Error()})
				log.Errorf("ID:%v\tWithdraw:%v", id, err)
				break
			default:
				c.JSON(500, gin.H{"error": err.Error()})
				log.Error(err.Error())
			}
		} else {
			c.JSON(200, gin.H{"data": result})
			log.Printf("Withdraw ID:%v %v", id, *result)
		}
		ch <- struct{}{}
	}()

	<-ch
	return

}

func (h *Handler) Deposit(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(401, gin.H{"error": "Authorization required"})
		log.Errorf("Auth error: Authorization required")
		return
	}
	id, ok := userId.(string)
	if !ok {
		c.JSON(401, gin.H{"error": "Invalid user ID"})
		log.Errorf("Invalid user ID")
		return
	}

	var response account.ResponseBalance
	if err := c.BindJSON(&response); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		log.Errorf("ID:%v\tFail bind balance:%v", id, err.Error())
		return
	}

	ch := make(chan struct{})

	go func() {
		result, err := h.service.Deposit(c.Request.Context(), id, response.Balance)
		if err != nil {
			switch err.Error() {
			case db.ErrNotValidAmount.Error():
				c.JSON(400, gin.H{"error": err.Error()})
				log.Errorf("ID:%v\tDeposit:%v", id, err.Error())
				break
			case db.ErrUpdateFailed.Error():
				c.JSON(404, gin.H{"error": err.Error()})
				log.Errorf("ID:%v\tDeposit:%v", id, err)
				break
			default:
				c.JSON(500, gin.H{"error": err.Error()})
				log.Error(err.Error())
			}
		} else {
			c.JSON(200, gin.H{"data": result})
			log.Printf("Deposit ID:%v %v", id, *result)
		}
		ch <- struct{}{}
	}()

	<-ch
	return
}

func (h *Handler) PutAccount(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(401, gin.H{"error": "Authorization required"})
		log.Errorf("Auth error: Authorization required")
		return
	}
	id, ok := userId.(string)
	if !ok {
		c.JSON(401, gin.H{"error": "Invalid user ID"})
		log.Errorf("Invalid user ID")
		return
	}

	var updateAccount account.Registration
	if err := c.BindJSON(&updateAccount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		log.Errorf("ID:%v\tFail bind balance:%v", id, err.Error())
		return
	}

	result, err := h.service.PutAccount(c.Request.Context(), id, updateAccount)
	if err != nil {
		switch err.Error() {
		case db.ErrUpdateFailed.Error():
			c.JSON(404, gin.H{"error": "Failed to update account"})
			log.Errorf("ID:%v\tPut:%v", id, err.Error())
			break
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Error(err.Error())
		}
		return
	}

	c.JSON(200, gin.H{"data": result})
	log.Printf("Put Account ID:%v %v", id, *result)
	return
}
