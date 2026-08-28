package handler

import (
	"net/http"

	"github.com/angyxys/kexel/internal/service"
	"github.com/gin-gonic/gin"
)

type TOTPHandler struct {
	totpServ *service.TOTPService
}

func NewTOTPHandler(totpServ *service.TOTPService) *TOTPHandler {
	return &TOTPHandler{
		totpServ: totpServ,
	}
}

type SetupTOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type VerifyTOTPRequest struct {
	Code string `json:"code" binding:"required,len=6"`
}

type VerifyBackupCodeRequest struct {
	Code string `json:"code" binding:"required"`
}

// GetTOTPSetup generates a new TOTP setup for the user
func (h *TOTPHandler) GetTOTPSetup(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
			"status":  http.StatusUnauthorized,
		})
		return
	}

	var req SetupTOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
		return
	}

	setup, err := h.totpServ.GenerateTOTPSecret(c, userID.(uint), req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"status":  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, setup)
}

// VerifyTOTP verifies and enables TOTP for the user
func (h *TOTPHandler) VerifyTOTP(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
			"status":  http.StatusUnauthorized,
		})
		return
	}

	var req VerifyTOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
		return
	}

	if err := h.totpServ.EnableTOTP(c, userID.(uint), req.Code); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "2FA enabled successfully",
		"status":  http.StatusOK,
	})
}

// GetTOTPStatus returns TOTP status for the authenticated user
func (h *TOTPHandler) GetTOTPStatus(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
			"status":  http.StatusUnauthorized,
		})
		return
	}

	status, err := h.totpServ.GetTOTPStatus(c, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to get TOTP status",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, status)
}

// DisableTOTP disables 2FA for the user
func (h *TOTPHandler) DisableTOTP(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
			"status":  http.StatusUnauthorized,
		})
		return
	}

	if err := h.totpServ.DisableTOTP(c, userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to disable 2FA",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "2FA disabled",
		"status":  http.StatusOK,
	})
}

// VerifyTOTPCode verifies a TOTP code (used during login)
func (h *TOTPHandler) VerifyTOTPCode(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
			"status":  http.StatusUnauthorized,
		})
		return
	}

	var req VerifyTOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
		return
	}

	valid, err := h.totpServ.VerifyTOTPCode(c, userID.(uint), req.Code)
	if err != nil || !valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "invalid code",
			"status":  http.StatusUnauthorized,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "code verified",
		"status":  http.StatusOK,
	})
}

// VerifyBackupCode verifies a backup code
func (h *TOTPHandler) VerifyBackupCode(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
			"status":  http.StatusUnauthorized,
		})
		return
	}

	var req VerifyBackupCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
		return
	}

	valid, err := h.totpServ.VerifyBackupCode(c, userID.(uint), req.Code)
	if err != nil || !valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "invalid backup code",
			"status":  http.StatusUnauthorized,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "backup code verified",
		"status":  http.StatusOK,
	})
}
