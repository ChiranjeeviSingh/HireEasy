package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "backend/internal/services"
)


func LoginH(ctx *gin.Context) {
    var loginReq services.LoginRequest

    if err := ctx.ShouldBindJSON(&loginReq); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid input", "error": err.Error()})
        return
    }

    token, err := services.Login(ctx, &loginReq)
    if err != nil {
        ctx.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized", "error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"token": token})
}


func RegisterH(ctx *gin.Context) {
    var registerReq services.RegisterRequest
    if err := ctx.ShouldBindJSON(&registerReq); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid input", "error": err.Error()})
        return
    }

    authResponse, err := services.Register(ctx, &registerReq)
    if err != nil {
        if err == services.ErrEmailExists {
            ctx.JSON(http.StatusBadRequest, gin.H{"msg": "Bad request", "error": err.Error()})
            return
        }
        ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "Could not create user", "error": err.Error()})
        return
    }

    ctx.JSON(http.StatusCreated, authResponse)
}