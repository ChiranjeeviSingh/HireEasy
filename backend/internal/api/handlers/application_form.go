package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "backend/internal/models"
    "backend/internal/services"
	"github.com/google/uuid"
)


func LinkJobToFormTemplateH(ctx *gin.Context) {
    // Extract job_id from URL
    jobID := ctx.Param("job_id")
   
    // Parse JSON request body
    var request models.LinkJobToFormRequest
    if err := ctx.ShouldBindJSON(&request); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Call the service to link job to form template
    form, err := services.LinkJobToFormTemplate(ctx, jobID, request.FormTemplateID)
    if err != nil {
        if err == services.ErrJobNotFound || err == services.ErrFormTemplateNotFound {
            ctx.JSON(http.StatusBadRequest, gin.H{"msg": "Bad request", "error": err.Error()})
            return
        }
        ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to link job to form template", "error": err.Error()})
        return
    }

    // Return the created form link
    ctx.JSON(http.StatusCreated, gin.H{"form_uuid": form.FormUUID})
}


func UpdateFormStatusH(ctx *gin.Context) {
    // Extract form_uuid from URL
    formUUID := ctx.Param("form_uuid")

	// Validate form_uuid format before proceeding
    if _, err := uuid.Parse(formUUID); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid form UUID format", "error": err.Error()})
        return
    }

    // Parse JSON request body
    var request models.UpdateFormStatusRequest
    if err := ctx.ShouldBindJSON(&request); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Call the service to update form status
    form, err := services.UpdateFormStatus(ctx, formUUID, request.Status)
    if err != nil {
        if err == services.ErrFormNotFound {
            ctx.JSON(http.StatusNotFound, gin.H{"msg": "Form not found", "error": err.Error()})
            return
        }
        ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to update form status", "error": err.Error()})
        return
    }

    // Return the updated form status
    ctx.JSON(http.StatusOK, gin.H{
        "form_uuid": form.FormUUID,
        "status":    form.Status,
    })
}

// GetFormDetailsH retrieves job and form template details based on form_uuid.
func GetFormDetailsH(ctx *gin.Context) {
    formUUID := ctx.Param("form_uuid")

	// Validate form_uuid format before proceeding
    if _, err := uuid.Parse(formUUID); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid form UUID format", "error": err.Error()})
        return
    }

    formResponse, err := services.GetFormDetails(ctx, formUUID)
    if err != nil {
        if err == services.ErrFormNotFound {
            ctx.JSON(http.StatusNotFound, gin.H{"msg": "Form not found", "error": err.Error()})
            return
        }
        ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to fetch form details", "error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, formResponse)
}

// DeleteFormH handles form deletion by form_uuid.
func DeleteFormH(ctx *gin.Context) {
    formUUID := ctx.Param("form_uuid")

	// Validate form_uuid format before proceeding
    if _, err := uuid.Parse(formUUID); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid form UUID format", "error": err.Error()})
        return
    }

    err := services.DeleteForm(ctx, formUUID)
    if err != nil {
        if err == services.ErrFormNotFound {
            ctx.JSON(http.StatusNotFound, gin.H{"msg": "Form not found", "error": err.Error()})
            return
        }
        ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to delete form", "error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "msg":        "Form successfully deleted",
        "form_uuid":  formUUID,
    })
}