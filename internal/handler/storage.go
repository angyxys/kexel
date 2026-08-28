package handler

import (
	"net/http"
	"strconv"

	"github.com/angyxys/kexel/internal/service"
	"github.com/gin-gonic/gin"
)

type StorageHandler struct {
	storageServ *service.StorageService
}

func NewStorageHandler(storageServ *service.StorageService) *StorageHandler {
	return &StorageHandler{
		storageServ: storageServ,
	}
}

type UpdateBannerRequest struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	DisplayOrder int    `json:"display_order"`
	IsActive     bool   `json:"is_active"`
}

// UploadBanner uploads a new banner image
func (h *StorageHandler) UploadBanner(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
			"status":  http.StatusUnauthorized,
		})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "error reading form",
			"status":  http.StatusBadRequest,
		})
		return
	}

	files := form.File["image"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "no image file provided",
			"status":  http.StatusBadRequest,
		})
		return
	}

	file, err := files[0].Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "error opening file",
			"status":  http.StatusBadRequest,
		})
		return
	}
	defer file.Close()

	bannerType := c.PostForm("type")
	title := c.PostForm("title")
	description := c.PostForm("description")

	if bannerType == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "banner type is required",
			"status":  http.StatusBadRequest,
		})
		return
	}

	banner, err := h.storageServ.UploadBanner(c, userID.(uint), file, files[0].Filename, bannerType, title, description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"status":  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusCreated, banner)
}

// GetUserBanners returns all banners for the authenticated user
func (h *StorageHandler) GetUserBanners(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
			"status":  http.StatusUnauthorized,
		})
		return
	}

	banners, err := h.storageServ.GetUserBanners(c, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to fetch banners",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   banners,
		"status": http.StatusOK,
	})
}

// GetBannersByType returns active banners of a specific type
func (h *StorageHandler) GetBannersByType(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
			"status":  http.StatusUnauthorized,
		})
		return
	}

	bannerType := c.Param("type")
	if bannerType == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "banner type is required",
			"status":  http.StatusBadRequest,
		})
		return
	}

	banners, err := h.storageServ.GetActiveBannersByType(c, userID.(uint), bannerType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to fetch banners",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   banners,
		"status": http.StatusOK,
	})
}

// UpdateBanner updates banner metadata
func (h *StorageHandler) UpdateBanner(c *gin.Context) {
	bannerIDStr := c.Param("id")
	bannerID, err := strconv.ParseUint(bannerIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid banner id",
			"status":  http.StatusBadRequest,
		})
		return
	}

	var req UpdateBannerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
		return
	}

	if err := h.storageServ.UpdateBannerInfo(c, uint(bannerID), req.Title, req.Description, req.DisplayOrder, req.IsActive); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to update banner",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "banner updated",
		"status":  http.StatusOK,
	})
}

// DeleteBanner deletes a banner
func (h *StorageHandler) DeleteBanner(c *gin.Context) {
	bannerIDStr := c.Param("id")
	bannerID, err := strconv.ParseUint(bannerIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid banner id",
			"status":  http.StatusBadRequest,
		})
		return
	}

	if err := h.storageServ.DeleteBanner(c, uint(bannerID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to delete banner",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "banner deleted",
		"status":  http.StatusOK,
	})
}
