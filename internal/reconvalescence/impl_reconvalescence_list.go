package reconvalescence

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rastislav-balcercik/reko-webapi/internal/db_service"
	"go.mongodb.org/mongo-driver/bson"
)

// CreateReconvalescenceTicket - Saves new entry into Reconvalescence list
func (this *implReconvalescenceTicketListAPI) CreateReconvalescenceTicket(ctx *gin.Context) {
	value, exists := ctx.Get("db_service")
	if !exists {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db not found",
				"error":   "db not found",
			})
		return
	}

	db, ok := value.(db_service.DbService[ReconvalescenceTicket])
	if !ok {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db context is not of required type",
				"error":   "cannot cast db context to db_service.DbService",
			})
		return
	}

	ticket := ReconvalescenceTicket{}
	err := ctx.BindJSON(&ticket)
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"status":  "Bad Request",
				"message": "Invalid request body",
				"error":   err.Error(),
			})
		return
	}

	if ticket.Id == "" || ticket.Id == "@new" {
		ticket.Id = uuid.New().String()
	}

	err = db.CreateDocument(ctx, ticket.Id, &ticket)

	switch err {
	case nil:
		ctx.JSON(
			http.StatusCreated,
			ticket,
		)
	case db_service.ErrConflict:
		ctx.JSON(
			http.StatusConflict,
			gin.H{
				"status":  "Conflict",
				"message": "Ticket already exists",
				"error":   err.Error(),
			},
		)
	default:
		ctx.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to create ticket in database",
				"error":   err.Error(),
			},
		)
	}
}

// DeleteReconvalescenceTicket - Deletes specific reconvalescence ticket
func (this *implReconvalescenceTicketListAPI) DeleteReconvalescenceTicket(ctx *gin.Context) {
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
	err := db.DeleteDocument(ctx, ticketId)

	switch err {
	case nil:
		ctx.AbortWithStatus(http.StatusNoContent)
	case db_service.ErrNotFound:
		ctx.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  "Not Found",
				"message": "Ticket not found",
				"error":   err.Error(),
			},
		)
	default:
		ctx.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to delete ticket from database",
				"error":   err.Error(),
			})
	}
}

// GetReconvalescenceList - Provides the reconvalence list
func (this *implReconvalescenceTicketListAPI) GetReconvalescenceList(ctx *gin.Context) {
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

	tickets := []ReconvalescenceTicket{}

	role := ctx.Query("role")
	if role == "doctor" {
		doctorId := ctx.Param("userId")
		t, err := db.FindDocuments(ctx, bson.D{{Key: "doctorid", Value: doctorId}})

		if err != nil {
			ctx.JSON(
				http.StatusBadGateway,
				gin.H{
					"status":  "Bad Gateway",
					"message": "Failed to fetch tickets from database",
					"error":   err.Error(),
				})
			return
		}

		for _, ticket := range t {
			tickets = append(tickets, *ticket)
		}

		ctx.JSON(
			http.StatusOK,
			tickets,
		)
		return
	}
	if role == "patient" {
		patientId := ctx.Param("userId")
		t, err := db.FindDocuments(ctx, bson.D{{Key: "patientid", Value: patientId}})

		if err != nil {
			ctx.JSON(
				http.StatusBadGateway,
				gin.H{
					"status":  "Bad Gateway",
					"message": "Failed to fetch tickets from database",
					"error":   err.Error(),
				})
			return
		}

		for _, ticket := range t {
			tickets = append(tickets, *ticket)
		}

		ctx.JSON(
			http.StatusOK,
			tickets,
		)
	}
}

// GetReconvalescenceTicket - Provides details about reconvalsescence ticket
func (this *implReconvalescenceTicketListAPI) GetReconvalescenceTicket(ctx *gin.Context) {
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
		ctx.JSON(
			http.StatusOK,
			ticket,
		)
	case db_service.ErrNotFound:
		ctx.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  "Not Found",
				"message": "Ticket not found",
				"error":   err.Error(),
			},
		)
	default:
		ctx.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to feth ticket from database",
				"error":   err.Error(),
			})
	}
}

// UpdateReconvalescenceTicket - Updates specific entry
func (this *implReconvalescenceTicketListAPI) UpdateReconvalescenceTicket(ctx *gin.Context) {
	value, exists := ctx.Get("db_service")
	if !exists {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db not found",
				"error":   "db not found",
			})
		return
	}

	db, ok := value.(db_service.DbService[ReconvalescenceTicket])
	if !ok {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db context is not of required type",
				"error":   "cannot cast db context to db_service.DbService",
			})
		return
	}

	ticket := ReconvalescenceTicket{}
	err := ctx.BindJSON(&ticket)
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"status":  "Bad Request",
				"message": "Invalid request body",
				"error":   err.Error(),
			})
		return
	}

	ticketId := ctx.Param("entryId")
	ticket.Id = ticketId
	err = db.UpdateDocument(ctx, ticketId, &ticket)

	switch err {
	case nil:
		ctx.JSON(
			http.StatusOK,
			ticket,
		)
	case db_service.ErrNotFound:
		ctx.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  "Not Found",
				"message": "Ticket not found",
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
			},
		)
	}
}
