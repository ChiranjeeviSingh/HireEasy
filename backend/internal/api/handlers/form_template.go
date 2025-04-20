package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "backend/internal/models"
    "backend/internal/services"
	"database/sql"
)

func CreateFormTemplateH(ctx *gin.Context) {
	var template models.FormTemplate
	if err := ctx.ShouldBindJSON(&template); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}


	if err := services.CreateFormTemplate(ctx, &template); err != nil {

        if err == services.ErrFormTemplateIdExists {
            ctx.JSON(http.StatusBadRequest, gin.H{"msg": "Bad request", "error": err.Error()})
            return
        }
        ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to create form template", "error": err.Error()})
        return

    }

	ctx.JSON(http.StatusCreated, template)
}


func GetFormTemplateH(ctx *gin.Context) {
	templateID := ctx.Param("form_template_id")
	template, err := services.GetFormTemplateById(ctx,templateID)
	if err != nil {

		if err == sql.ErrNoRows {
            ctx.JSON(http.StatusNotFound, gin.H{"error": "Form Template not found"})
            return
        }

		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to retrieve form template", "error": err.Error()})
        return

	}

	ctx.JSON(http.StatusOK, template)
}

func ListFormTemplatesH(ctx *gin.Context) {
	templates, err := services.GetFormTemplatesByUserId(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to retrieve all form templates", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, templates)
}

func DeleteFormTemplateH(ctx *gin.Context) {
	templateID := ctx.Param("form_template_id")
	if err := services.DeleteFormTemplate(ctx, templateID); err != nil {

		if err == services.ErrFormTemplateIdDoesNotExists {
            ctx.JSON(http.StatusBadRequest, gin.H{"msg": "Bad request", "error": err.Error()})
            return
        }

        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete template"})
        return

	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Form template deleted successfully"})
}
