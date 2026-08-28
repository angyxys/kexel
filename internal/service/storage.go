package service

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"time"

	"github.com/angyxys/kexel/internal/database/models"
	"github.com/angyxys/kexel/internal/repository"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type StorageService struct {
	bannerRepo  *repository.BannerRepository
	minioClient *minio.Client
	bucketName  string
}

// NewMinIOClient initializes MinIO client from environment
func NewMinIOClient() (*minio.Client, error) {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:9000"
	}

	useSSL := os.Getenv("MINIO_USE_SSL") == "true"
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	return client, nil
}

func NewStorageService(bannerRepo *repository.BannerRepository, minioClient *minio.Client, bucketName string) *StorageService {
	return &StorageService{
		bannerRepo:  bannerRepo,
		minioClient: minioClient,
		bucketName:  bucketName,
	}
}

// UploadBanner uploads an image to MinIO and creates a banner record
func (s *StorageService) UploadBanner(ctx context.Context, userID uint, file io.Reader, filename string, bannerType string, title string, description string) (*BannerInfo, error) {
	// Read file into memory to get size and check image
	fileData, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	// Validate image
	config, _, err := image.DecodeConfig(bytes.NewReader(fileData))
	if err != nil {
		return nil, fmt.Errorf("invalid image file: %w", err)
	}

	// Generate S3 key
	s3Key := fmt.Sprintf("banners/%d/%d-%s", userID, time.Now().Unix(), filename)

	// Upload to MinIO
	uploadInfo, err := s.minioClient.PutObject(ctx, s.bucketName, s3Key, bytes.NewReader(fileData), int64(len(fileData)), minio.PutObjectOptions{
		ContentType: "image/jpeg", // Should be detected from file
	})
	if err != nil {
		return nil, fmt.Errorf("error uploading to MinIO: %w", err)
	}

	// Create banner record
	banner := &models.Banner{
		UserID:      userID,
		Name:        filename,
		Type:        bannerType,
		Title:       title,
		Description: description,
		S3Key:       s3Key,
		ImageURL:    fmt.Sprintf("https://minio.example.com/%s/%s", s.bucketName, s3Key),
		Width:       config.Width,
		Height:      config.Height,
		FileSize:    uploadInfo.Size,
		MimeType:    "image/jpeg",
		IsActive:    true,
	}

	if err := s.bannerRepo.Create(ctx, banner); err != nil {
		// Clean up uploaded file
		_ = s.minioClient.RemoveObject(ctx, s.bucketName, s3Key, minio.RemoveObjectOptions{})
		return nil, fmt.Errorf("error creating banner record: %w", err)
	}

	return s.getBannerInfo(banner), nil
}

type BannerInfo struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	Type         string    `json:"type"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	ImageURL     string    `json:"image_url"`
	Width        int       `json:"width"`
	Height       int       `json:"height"`
	FileSize     int64     `json:"file_size"`
	IsActive     bool      `json:"is_active"`
	DisplayOrder int       `json:"display_order"`
	CreatedAt    time.Time `json:"created_at"`
}

func (s *StorageService) getBannerInfo(banner *models.Banner) *BannerInfo {
	return &BannerInfo{
		ID:           banner.ID,
		Name:         banner.Name,
		Type:         banner.Type,
		Title:        banner.Title,
		Description:  banner.Description,
		ImageURL:     banner.ImageURL,
		Width:        banner.Width,
		Height:       banner.Height,
		FileSize:     banner.FileSize,
		IsActive:     banner.IsActive,
		DisplayOrder: banner.DisplayOrder,
		CreatedAt:    banner.CreatedAt,
	}
}

// GetUserBanners returns all banners for a user
func (s *StorageService) GetUserBanners(ctx context.Context, userID uint) ([]BannerInfo, error) {
	banners, err := s.bannerRepo.ListUserBanners(ctx, userID)
	if err != nil {
		return nil, err
	}

	infos := make([]BannerInfo, len(banners))
	for i, banner := range banners {
		infos[i] = *s.getBannerInfo(&banner)
	}
	return infos, nil
}

// GetActiveBannersByType returns active banners of a specific type
func (s *StorageService) GetActiveBannersByType(ctx context.Context, userID uint, bannerType string) ([]BannerInfo, error) {
	banners, err := s.bannerRepo.ListActiveBannersByType(ctx, userID, bannerType)
	if err != nil {
		return nil, err
	}

	infos := make([]BannerInfo, len(banners))
	for i, banner := range banners {
		infos[i] = *s.getBannerInfo(&banner)
	}
	return infos, nil
}

// DeleteBanner deletes a banner and removes the file from MinIO
func (s *StorageService) DeleteBanner(ctx context.Context, bannerID uint) error {
	banner, err := s.bannerRepo.GetByID(ctx, bannerID)
	if err != nil {
		return err
	}

	// Delete from MinIO
	if err := s.minioClient.RemoveObject(ctx, s.bucketName, banner.S3Key, minio.RemoveObjectOptions{}); err != nil {
		return fmt.Errorf("error deleting from MinIO: %w", err)
	}

	// Delete from database
	return s.bannerRepo.Delete(ctx, bannerID)
}

// UpdateBannerInfo updates banner metadata
func (s *StorageService) UpdateBannerInfo(ctx context.Context, bannerID uint, title string, description string, displayOrder int, isActive bool) error {
	banner, err := s.bannerRepo.GetByID(ctx, bannerID)
	if err != nil {
		return err
	}

	banner.Title = title
	banner.Description = description
	banner.DisplayOrder = displayOrder
	banner.IsActive = isActive

	return s.bannerRepo.Update(ctx, banner)
}

// GeneratePresignedURL generates a temporary download URL
func (s *StorageService) GeneratePresignedURL(ctx context.Context, s3Key string, expiration time.Duration) (string, error) {
	url, err := s.minioClient.PresignedGetObject(ctx, s.bucketName, s3Key, expiration, nil)
	if err != nil {
		return "", fmt.Errorf("error generating presigned URL: %w", err)
	}
	return url.String(), nil
}
