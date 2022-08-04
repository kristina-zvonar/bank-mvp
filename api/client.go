package api

import (
	db "bank-mvp/db/sqlc"
	"database/sql"
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

type getClientRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getClient(ctx *gin.Context) {
	var req getClientRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	client, err := server.store.GetClient(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, client)
}

type updateClientRequest struct {
	ID int64 `json:"id" binding:"required,min=1"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	CountryID int64 `json:"country_id"`
	Active bool `json:"active"`
}

func (server *Server) updateClient(ctx *gin.Context) {
	var req updateClientRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateClientParams {
		ID: req.ID,
		FirstName: req.FirstName,
		LastName: req.LastName,
		CountryID: req.CountryID,
		Active: req.Active,
	}
	client, err := server.store.UpdateClient(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, client)
}

type listClientsParams struct {
	Page int32 `form:"page" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=100"`
}

func (server *Server) ListClients(ctx *gin.Context) {
	var req listClientsParams
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListClientsParams {
		Limit: req.PageSize,
		Offset: (req.Page - 1) * req.PageSize,
	}
	clients, err := server.store.ListClients(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, clients)
}