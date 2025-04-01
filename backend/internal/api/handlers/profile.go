package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "backend/internal/services"
	"backend/internal/models"
)

func CreateProfileH(ctx *gin.Context) {
	var profileReq models.Profile
	if err := ctx.ShouldBindJSON(&profileReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid input", "error": err.Error()})
		return
	}

	err := services.CreateProfile(ctx, &profileReq)
	if err != nil {

		if err == services.ErrProfileExists{
			ctx.JSON(http.StatusBadRequest, gin.H{"msg": "Bad request", "error": err.Error()})
			return} 
		
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "Could not create profile", "error": err.Error()})
		return
	}


	ctx.JSON(http.StatusCreated, profileReq)
}

func GetMyProfileH(ctx *gin.Context) {
	profile, err := services.GetMyProfile(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "Could not get profile", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, profile)
}

func GetUserProfileH(ctx *gin.Context) {
		
	userName := ctx.Param("user_name")

	profile, err := services.GetUserProfile(ctx, userName)
	if err != nil {

		if err==services.ErrUsernameNotFound{
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "Could not get profile", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, profile)
}

func UpdateMyProfileH(ctx *gin.Context) {
	var profileReq models.Profile
	if err := ctx.ShouldBindJSON(&profileReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid input", "error": err.Error()})
		return
	}

	err := services.UpdateMyProfile(ctx, &profileReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "Could not update profile", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, profileReq)
}