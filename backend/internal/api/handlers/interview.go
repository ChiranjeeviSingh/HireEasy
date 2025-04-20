package handlers

import (
	"backend/internal/models"
	"backend/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateInterviewH handles the creation of a new interview
func CreateInterviewH(ctx *gin.Context) {
	
	var req models.CreateInterviewRequest
	
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	interview, err := services.CreateInterview(ctx, &req); 

	if err != nil {
		if err == services.ErrOverlappingSlot {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create interview", "msg": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, interview)
}

// DeleteInterviewH handles the deletion of an interview
func DeleteInterviewH(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := services.DeleteInterview(ctx, id); err != nil {
		if err == services.ErrInterviewNotFound {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete interview", "msg": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Interview deleted successfully"})
}

// ListInterviewsH retrieves all interviews for the authenticated user
func ListAllInterviewsH(ctx *gin.Context) {
	interviews, err := services.ListAllInterviews(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve interviews", "msg": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, interviews)
}

//SubmitFeedbackH handles the submission of interview feedback by 
func SubmitFeedbackH(ctx *gin.Context) {
	
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var req models.FeedbackRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.SubmitFeedback(ctx, &req, id); err != nil {
		if err == services.ErrInterviewNotFound || err == services.ErrUnauthorizedFeedback {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit feedback", "msg": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Feedback submitted successfully"})
}