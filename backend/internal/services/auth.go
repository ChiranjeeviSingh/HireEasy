package services

import (
    "context"
    "errors"
    "time"
    "github.com/golang-jwt/jwt/v4"
    "golang.org/x/crypto/bcrypt"
    "backend/internal/database"
    "backend/internal/config"
    "backend/internal/models"
)

var (
    ErrInvalidCredentials = errors.New("invalid credentials")
    ErrEmailExists       = errors.New("email already exists")
)

type RegisterRequest struct {
    Username string `json:"username" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
    Role string     `json:"role" binding:"required"`
    CompanyName string `json:"company_name" binding:"required"` 
}

type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
    Token string `json:"token"`
    User  models.User `json:"user"`
}

func Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error) {
    db := database.GetDB()

    // Check if email exists
    var exists bool
    err := db.GetContext(ctx, &exists, 
        "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", req.Email)
    if err != nil {
        return nil, err
    }
    if exists {
        return nil, ErrEmailExists
    }

    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }

    // Create user
    var userId int
    err = db.GetContext(ctx, &userId,
        `INSERT INTO users (email, password_hash, username, role, company_name) 
         VALUES ($1, $2, $3, $4, $5) 
         RETURNING id`, req.Email, string(hashedPassword), req.Username, req.Role, req.CompanyName)
    if err != nil {
        return nil, err
    }

    // Generate token
    token, err := generateToken(userId)
    if err != nil {
        return nil, err
    }

    return &AuthResponse{
        Token: token,
        User: models.User {
            Username: req.Username,
            Email: req.Email,
        },
    }, nil
}

func Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error) {

    var user models.User

    db := database.GetDB()

    err := db.GetContext(ctx, &user,
        `SELECT id, email, password_hash, username 
         FROM users 
         WHERE email = $1`, req.Email)
    if err != nil {
        return nil, ErrInvalidCredentials
    }

    // Verify password
    err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
    if err != nil {
        return nil, ErrInvalidCredentials
    }

    // Generate token
    token, err := generateToken(user.ID)
    if err != nil {
        return nil, err
    }

    return &AuthResponse{
        Token: token,
        User: models.User {
            Username: user.Username,
            Email: req.Email,
        },
    }, nil
}

func generateToken(userID int) (string, error) {
    
    claims := jwt.MapClaims{
        "user_id": userID,
        "exp":     time.Now().Add(24 * time.Hour).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(config.GetConfig().JWTSecret))
}