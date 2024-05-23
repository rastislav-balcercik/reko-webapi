package reconvalescence

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rastislav-balcercik/reko-webapi/internal/db_service"
)

type ticketUpdater = func(
	ctx *gin.Context,
	ticket *ReconvalescenceTicket,
) (updatedTicket *ReconvalescenceTicket, responseContent interface{}, status int)

func updateTicketFunc(ctx *gin.Context, updater ticketUpdater) {
	value, exists := ctx.Get("db_service")
	if !exists {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db_service not found",
				"error":   "db_service not found",
			})
		return
	}

	db, ok := value.(db_service.DbService[ReconvalescenceTicket])
	if !ok {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db_service context is not of type db_service.DbService",
				"error":   "cannot cast db_service context to db_service.DbService",
			})
		return
	}

	ticketId := ctx.Param("entryId")

	ticket, err := db.FindDocument(ctx, ticketId)

	switch err {
	case nil:
		// continue
	case db_service.ErrNotFound:
		ctx.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  "Not Found",
				"message": "Ticket not found",
				"error":   err.Error(),
			},
		)
		return
	default:
		ctx.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to load ticket from database",
				"error":   err.Error(),
			})
		return
	}

	if !ok {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "Failed to cast ticket from database",
				"error":   "Failed to cast ticket from database",
			})
		return
	}

	updatedTicket, responseObject, status := updater(ctx, ticket)

	if updatedTicket != nil {
		err = db.UpdateDocument(ctx, ticketId, updatedTicket)
	} else {
		err = nil // redundant but for clarity
	}

	switch err {
	case nil:
		if responseObject != nil {
			ctx.JSON(status, responseObject)
		} else {
			ctx.AbortWithStatus(status)
		}
	case db_service.ErrNotFound:
		ctx.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  "Not Found",
				"message": "Ticket was deleted while processing the request",
				"error":   err.Error(),
			},
		)
	default:
		ctx.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to update ticket in database",
				"error":   err.Error(),
			})
	}

}
