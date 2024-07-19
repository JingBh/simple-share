package oss

import (
	"context"
	"encoding/json"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/jingbh/simple-share/internal/models"
	"github.com/jingbh/simple-share/internal/utils"
	"net/url"
	"strconv"
	"strings"
)

// CreateShareOptions Request to create a single shared file in the OSS store.
// If `Source` is provided, the file is already uploaded to `uploads/` and should be moved to the destination.
type CreateShareOptions struct {
	Type     string // `file`, `directory`, `text`, `url`
	Text     string
	Source   string
	Name     string // name of the file
	Path     string // path to save the file, after `shares/`
	Password string
	Expiry   int
	Creator  *models.ShareCreator
}

func CreateShare(ctx context.Context, options CreateShareOptions) error {
	client := Client()

	ossOptions := []oss.Option{
		oss.WithContext(ctx),
		oss.CacheControl("private, max-age=86400"),
		oss.Meta("Share-Type", options.Type),
		oss.Meta("Share-Expiry", strconv.Itoa(options.Expiry)),
	}
	if options.Name != "" {
		nameEncoded := url.PathEscape(options.Name)
		ossOptions = append(ossOptions, oss.ContentDisposition("attachment; filename=\""+nameEncoded+"\"; filename*=UTF-8''"+nameEncoded))
	}
	if options.Creator != nil {
		creatorJsonBytes, err := json.Marshal(options.Creator)
		if err == nil {
			creatorJson := string(creatorJsonBytes)
			ossOptions = append(ossOptions, oss.Meta("Share-Creator", creatorJson))
		}
	}
	if options.Password != "" {
		passwordHashed, err := utils.HashPassword(options.Password)
		if err != nil {
			return err
		}
		ossOptions = append(ossOptions, oss.Meta("Share-Password", passwordHashed))
	}
	if options.Expiry > 0 {
		ossOptions = append(ossOptions, oss.Meta("Share-Expiry", strconv.Itoa(options.Expiry)))
		ossOptions = append(ossOptions, oss.SetTagging(oss.Tagging{Tags: []oss.Tag{
			{Key: "period", Value: strconv.Itoa(options.Expiry)},
		}}))
	}

	// no need to add retry here, as the source file is not deleted,
	// the client can actively retry
	if options.Source != "" {
		ossOptions = append(ossOptions, oss.TaggingDirective(oss.TaggingReplace))
		_, err := client.CopyObject(
			"uploads/"+options.Source,
			"shares/"+options.Path,
			ossOptions...,
		)
		return err
	} else {
		md5 := utils.MD5HashBase64([]byte(options.Text))
		ossOptions = append(ossOptions, oss.ContentMD5(md5))
		err := client.PutObject(
			"shares/"+options.Path,
			strings.NewReader(options.Text),
			ossOptions...,
		)
		return err
	}
}
