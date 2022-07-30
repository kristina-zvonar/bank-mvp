package api

import (
	db "bank-mvp/db/sqlc"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createClientRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	CountryID int64  `json:"country_id" binding:"required,min=1,max=228"`
	UserID    int64  `json:"user_id" binding:"required,min=1"`
}

func (server *Server) createClient(ctx *gin.Context) {
	var req createClientRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateClientParams {
		FirstName: req.FirstName,
		LastName: req.LastName,
		CountryID: req.CountryID,
		UserID: req.UserID,
	}

	client, err := server.store.CreateClient(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			log.Println(pqErr)
			switch pqErr.Code.Name() {
			case "foreign_key_violation":
				customErr := fmt.Errorf("country or user id does not exist")
				ctx.JSON(http.StatusForbidden, errorResponse(customErr))
				return
			}

		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, client)
}