package cloudinary

import (
	"context"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/google/uuid"
)

type CloudinaryService interface {
	UploadImage(fileHeader *multipart.FileHeader) (string, error)
}

type cloudinaryService struct {
	cloudName string
	apiKey    string
	apiSecret string
}

func NewCloudinaryService(cloudName, apiKey, apiSecret string) CloudinaryService {
	return &cloudinaryService{
		cloudName: cloudName,
		apiKey:    apiKey,
		apiSecret: apiSecret,
	}
}

func (cs *cloudinaryService) UploadImage(fileHeader *multipart.FileHeader) (string, error) {
	cld, err := cloudinary.NewFromParams(cs.cloudName, cs.apiKey, cs.apiSecret)
	if err != nil {
		return "", err
	}

	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	publicID := uuid.New().String()

	ctx := context.Background()
	uploadResult, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder:   "go_final/reports",
		PublicID: publicID,
	})
	if err != nil {
		return "", err
	}

	return uploadResult.SecureURL, nil
}
