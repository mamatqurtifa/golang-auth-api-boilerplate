package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mamatqurtifa/golang-auth-api-boilerplate/database"
	"github.com/mamatqurtifa/golang-auth-api-boilerplate/models"
	"github.com/mamatqurtifa/golang-auth-api-boilerplate/services"
	"github.com/mamatqurtifa/golang-auth-api-boilerplate/utils"
)

type AuthController struct {
	emailService *services.EmailService
}

// NewAuthController creates a new auth controller
func NewAuthController() *AuthController {
	return &AuthController{
		emailService: services.NewEmailService(),
	}
}

// RegisterRequest represents registration request body
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
}

// LoginRequest represents login request body
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// ForgotPasswordRequest represents forgot password request body
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest represents reset password request body
type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// ChangePasswordRequest represents change password request body
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// Register handles user registration
func (ctrl *AuthController) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Check if user already exists
	var existingUser models.User
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		utils.ErrorResponse(c, http.StatusConflict, "Email already registered", "email_exists")
		return
	}

	// Generate verification token
	verificationToken, err := utils.GenerateRandomToken(32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate verification token", err.Error())
		return
	}

	// Create new user
	user := models.User{
		Email:             req.Email,
		Password:          req.Password,
		Name:              req.Name,
		VerificationToken: verificationToken,
		IsEmailVerified:   false,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create user", err.Error())
		return
	}

	// Send verification email
	go func() {
		if err := ctrl.emailService.SendVerificationEmail(user.Email, verificationToken); err != nil {
			// Log error but don't fail the request
			println("Failed to send verification email:", err.Error())
		}
	}()

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate token", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "User registered successfully", gin.H{
		"user":  user,
		"token": token,
	})
}

// Login handles user login
func (ctrl *AuthController) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Find user by email
	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid email or password", "invalid_credentials")
		return
	}

	// Check password
	if !user.CheckPassword(req.Password) {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid email or password", "invalid_credentials")
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate token", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Login successful", gin.H{
		"user":  user,
		"token": token,
	})
}

// GetProfile returns the authenticated user's profile
func (ctrl *AuthController) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", "not_authenticated")
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Profile retrieved successfully", user)
}

// UpdateProfile updates the authenticated user's profile
func (ctrl *AuthController) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", "not_authenticated")
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found", err.Error())
		return
	}

	var req struct {
		Name string `json:"name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Update user
	if req.Name != "" {
		user.Name = req.Name
	}

	if err := database.DB.Save(&user).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update profile", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Profile updated successfully", user)
}

// ForgotPassword handles forgot password request
func (ctrl *AuthController) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Find user by email
	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		// Don't reveal if user exists or not
		utils.SuccessResponse(c, http.StatusOK, "If the email exists, a reset link has been sent", nil)
		return
	}

	// Generate reset token
	resetToken, err := utils.GenerateRandomToken(32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate reset token", err.Error())
		return
	}

	// Set token and expiry (1 hour)
	expiryTime := time.Now().Add(1 * time.Hour)
	user.ResetToken = resetToken
	user.ResetTokenExpiry = &expiryTime

	if err := database.DB.Save(&user).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to save reset token", err.Error())
		return
	}

	// Send reset email
	go func() {
		if err := ctrl.emailService.SendPasswordResetEmail(user.Email, resetToken); err != nil {
			println("Failed to send reset email:", err.Error())
		}
	}()

	utils.SuccessResponse(c, http.StatusOK, "If the email exists, a reset link has been sent", nil)
}

// ResetPassword handles password reset
func (ctrl *AuthController) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Find user by reset token
	var user models.User
	if err := database.DB.Where("reset_token = ?", req.Token).First(&user).Error; err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid or expired reset token", "invalid_token")
		return
	}

	// Check if token is expired
	if user.ResetTokenExpiry == nil || time.Now().After(*user.ResetTokenExpiry) {
		utils.ErrorResponse(c, http.StatusBadRequest, "Reset token has expired", "token_expired")
		return
	}

	// Hash new password
	hashedPassword, err := models.HashPassword(req.NewPassword)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to hash password", err.Error())
		return
	}

	// Update password and clear reset token
	user.Password = hashedPassword
	user.ResetToken = ""
	user.ResetTokenExpiry = nil

	if err := database.DB.Save(&user).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to reset password", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Password reset successfully", nil)
}

// ChangePassword handles password change for authenticated users
func (ctrl *AuthController) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", "not_authenticated")
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "User not found", err.Error())
		return
	}

	// Check old password
	if !user.CheckPassword(req.OldPassword) {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Current password is incorrect", "invalid_password")
		return
	}

	// Hash new password
	hashedPassword, err := models.HashPassword(req.NewPassword)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to hash password", err.Error())
		return
	}

	// Update password
	user.Password = hashedPassword
	if err := database.DB.Save(&user).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to change password", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Password changed successfully", nil)
}

// VerifyEmail handles email verification
func (ctrl *AuthController) VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Verification token is required", "missing_token")
		return
	}

	// Find user by verification token
	var user models.User
	if err := database.DB.Where("verification_token = ?", token).First(&user).Error; err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid verification token", "invalid_token")
		return
	}

	// Update user as verified
	user.IsEmailVerified = true
	user.VerificationToken = ""

	if err := database.DB.Save(&user).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to verify email", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Email verified successfully", nil)
}

// Logout handles user logout (client-side token removal)
func (ctrl *AuthController) Logout(c *gin.Context) {
	// In a JWT-based system, logout is typically handled client-side
	// by removing the token from storage.
	// For additional security, you could implement token blacklisting here.
	utils.SuccessResponse(c, http.StatusOK, "Logged out successfully", nil)
}
