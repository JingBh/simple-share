package oss

import (
	"context"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"io"
	"sync"
	"time"
)

// UploadPartSize 20MB
const UploadPartSize = 20 * 1024 * 1024

type uploadInfoRecord struct {
	mu         sync.RWMutex
	FileId     string
	InitResult oss.InitiateMultipartUploadResult
	PartResult []oss.UploadPart
	StartedAt  time.Time
}

var uploadInfoMap sync.Map

// generateFileId generates a UUID as the uploaded file key.
func generateFileId() (string, error) {
	instance, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	return instance.String(), nil
}

func UploadInit(ctx context.Context) (string, error) {
	client := Client()

	fileId, err := generateFileId()
	if err != nil {
		return "", err
	}
	key := "uploads/" + fileId + ".bin"
	ossOptions := []oss.Option{
		oss.WithContext(ctx),
		oss.ContentType("application/octet-stream"),
	}
	res, err := client.InitiateMultipartUpload(key, ossOptions...)
	if err != nil {
		return "", err
	}
	uploadInfoMap.Store(fileId, &uploadInfoRecord{
		FileId:     fileId,
		InitResult: res,
		StartedAt:  time.Now(),
	})

	return fileId, nil
}

func UploadPart(ctx context.Context, fileId string, partNumber int, body io.Reader, size int64) error {
	client := Client()

	uploadInfoVal, ok := uploadInfoMap.Load(fileId)
	if !ok {
		return errors.Errorf("file id not found: %s", fileId)
	}
	uploadInfo := uploadInfoVal.(*uploadInfoRecord)
	uploadInfo.mu.Lock()
	defer uploadInfo.mu.Unlock()

	ossOptions := []oss.Option{
		oss.WithContext(ctx),
	}
	res, err := client.UploadPart(uploadInfo.InitResult, body, size, partNumber, ossOptions...)
	if err != nil {
		return err
	}
	uploadInfo.PartResult = append(uploadInfo.PartResult, res)

	return nil
}

func UploadComplete(ctx context.Context, fileId string) error {
	client := Client()

	uploadInfoVal, ok := uploadInfoMap.Load(fileId)
	if !ok {
		return errors.Errorf("file id not found: %s", fileId)
	}
	uploadInfo := uploadInfoVal.(*uploadInfoRecord)
	uploadInfo.mu.RLock()

	ossOptions := []oss.Option{
		oss.WithContext(ctx),
	}
	_, err := client.CompleteMultipartUpload(uploadInfo.InitResult, uploadInfo.PartResult, ossOptions...)
	if err != nil {
		return err
	}

	uploadInfo.mu.RUnlock()
	uploadInfoMap.Delete(fileId)
	return nil
}

func init() {
	go func() {
		// Set up a coroutine to clean up upload info older than a day
		for {
			uploadInfoMap.Range(func(key, value interface{}) bool {
				uploadInfo := value.(*uploadInfoRecord)
				if time.Since(uploadInfo.StartedAt) > 24*time.Hour {
					uploadInfoMap.Delete(key)
				}
				return true
			})
			time.Sleep(10 * time.Minute)
		}
	}()
}
