package oss

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/jingbh/simple-share/internal/models"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type GetShareContentOptions struct {
	Name    string
	FileId  string
	Headers http.Header
}

type GetShareContentResult struct {
	Reader  io.ReadCloser
	Headers http.Header
}

func GetShare(ctx context.Context, name string) (*models.Share, error) {
	client := Client()

	key := "shares/" + name
	res, err := client.GetObjectDetailedMeta(key, oss.WithContext(ctx))
	if err != nil {
		var ossErr oss.ServiceError
		if errors.As(err, &ossErr) && ossErr.Code == "NoSuchKey" {
			return nil, nil
		}
		return nil, err
	}

	shareType := res.Get(oss.HTTPHeaderOssMetaPrefix + "Share-Type")
	expiry, _ := strconv.Atoi(res.Get(oss.HTTPHeaderOssMetaPrefix + "Share-Expiry"))
	size, _ := strconv.ParseInt(res.Get(oss.HTTPHeaderContentLength), 10, 64)
	creatorJson := res.Get(oss.HTTPHeaderOssMetaPrefix + "Share-Creator")
	var creator *models.ShareCreator = nil
	if creatorJson != "" {
		creator = new(models.ShareCreator)
		_ = json.Unmarshal([]byte(creatorJson), creator)
	}
	createdAt, _ := http.ParseTime(res.Get(oss.HTTPHeaderLastModified))

	var files models.ShareFiles = nil
	if shareType == "directory" {
		// get file tree and calculate total size
		filesJsonReader, err := client.GetObject(key, oss.WithContext(ctx))
		if err == nil {
			err = json.NewDecoder(filesJsonReader).Decode(&files)
		}

		var dirSize int64 = 0
		continuationToken := ""
		for {
			ossListOptions := []oss.Option{
				oss.WithContext(ctx),
				oss.MaxKeys(1000),
				oss.Prefix(key),
			}
			if continuationToken != "" {
				ossListOptions = append(ossListOptions, oss.ContinuationToken(continuationToken))
			}

			dirRes, err := client.ListObjectsV2(ossListOptions...)
			if err != nil {
				break
			}
			for _, object := range dirRes.Objects {
				if files != nil {
					for fileKey, file := range files {
						if object.Key == key+".d/"+file.Id+".bin" {
							files[fileKey].Size = object.Size
							break
						}
					}
				}
				dirSize += object.Size
			}
			if !dirRes.IsTruncated {
				break
			}
			continuationToken = dirRes.NextContinuationToken
		}
		if dirSize > 0 {
			size = dirSize
		}
	} else if shareType == "file" {
		filename := res.Get(oss.HTTPHeaderOssMetaPrefix + "Share-Filename")
		files = models.ShareFiles{{
			Path: filename,
			Size: size,
		}}
	}

	return &models.Share{
		Type:      shareType,
		Name:      name,
		Password:  res.Get(oss.HTTPHeaderOssMetaPrefix + "Share-Password"),
		Expiry:    expiry,
		Size:      size,
		CreatedAt: createdAt,
		Files:     files,
		Creator:   creator,
	}, nil
}

func ListShares(ctx context.Context, continuationToken string) ([]*models.Share, string, error) {
	client := Client()

	var result []*models.Share
	ossOptions := []oss.Option{
		oss.WithContext(ctx),
		oss.MaxKeys(10),
		oss.Prefix("shares/"),
		oss.Delimiter("/"),
	}
	if continuationToken != "" {
		ossOptions = append(ossOptions, oss.ContinuationToken(continuationToken))
	}
	res, err := client.ListObjectsV2(ossOptions...)
	if err != nil {
		return nil, "", err
	}

	for _, object := range res.Objects {
		name := strings.TrimPrefix(object.Key, res.Prefix)
		share, _ := GetShareCached(ctx, name)
		if share != nil {
			result = append(result, share)
		}
	}
	if res.IsTruncated {
		return result, res.NextContinuationToken, nil
	} else {
		return result, "", nil
	}
}

func GetShareContent(ctx context.Context, options GetShareContentOptions) (*GetShareContentResult, error) {
	client := Client()

	key := "shares/" + options.Name
	if options.FileId != "" {
		key += ".d/" + options.FileId + ".bin"
	}

	ossOptions := []oss.Option{
		oss.WithContext(ctx),
	}
	if options.Headers != nil {
		for k, vs := range options.Headers {
			for _, v := range vs {
				ossOptions = append(ossOptions, oss.SetHeader(k, v))
			}
		}
	}

	res, err := client.DoGetObject(&oss.GetObjectRequest{
		ObjectKey: key,
	}, ossOptions)
	if err != nil {
		var ossErr oss.ServiceError
		if errors.As(err, &ossErr) && ossErr.Code == "NoSuchKey" {
			return nil, nil
		}
		return nil, err
	}
	resp := res.Response

	return &GetShareContentResult{
		Reader:  resp.Body,
		Headers: resp.Headers,
	}, nil
}
