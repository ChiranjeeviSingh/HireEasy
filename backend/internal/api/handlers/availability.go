package handlers

import (
	"backend/internal/models"
	"backend/internal/services"
	"net/http"
	"github.com/gin-gonic/gin"
	"strconv"
)

// CreateAvailabilityH handles the creation of a new availability slot
func CreateAvailabilityH(ctx *gin.Context) {
	var availability models.Availability
	if err := ctx.ShouldBindJSON(&availability); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.CreateAvailability(ctx, &availability); err != nil {

		if err == services.ErrMaxDurationExceeded || err == services.ErrInvalidTime || err == services.ErrOverlappingSlot || err == services.ErrSlotScheduled {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create availability slot", "msg": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, availability)
}

// DeleteAvailabilityH deletes a specific availability slot
func DeleteAvailabilityH(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	
	if err := services.DeleteAvailability(ctx, id); err != nil {

		if err == services.ErrNoAvailability{
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}


		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete availability slot", "msg": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Availability slot deleted successfully"})
}



// GetMyAvailabilityH retrieves the availability slots for the authenticated user
func GetMyAvailabilityH(ctx *gin.Context) {
	slots, err := services.GetMyAvailability(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve availability slots", "msg": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, slots)
}

// GetUserAvailabilityH retrieves the availability slots for a specific user
func GetUserAvailabilityH(ctx *gin.Context) {

	userName := ctx.Param("user_name")

	slots, err := services.GetUserAvailability(ctx, userName)
	if err != nil {

		if err==services.ErrUsernameNotFound{
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve availability slots", "msg": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, slots)
}

// GetAvailablePeopleH retrieves all available people
func GetAllAvailabilityH(ctx *gin.Context) {
	
	allSlots, err := services.GetAllAvailability(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve available people", "msg": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, allSlots)
}