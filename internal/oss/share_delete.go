package oss

import (
	"context"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func DeleteShare(ctx context.Context, name string) error {
	client := Client()

	continuationToken := ""
	for {
		ossListOptions := []oss.Option{
			oss.WithContext(ctx),
			oss.MaxKeys(1000),
			oss.Prefix("shares/" + name),
		}
		if continuationToken != "" {
			ossListOptions = append(ossListOptions, oss.ContinuationToken(continuationToken))
		}
		res, err := client.ListObjectsV2(ossListOptions...)
		if err != nil {
			return err
		}

		var keys []string
		for _, obj := range res.Objects {
			keys = append(keys, obj.Key)
		}
		if len(keys) == 0 {
			break
		}
		ossDeleteOptions := []oss.Option{
			oss.WithContext(ctx),
			oss.DeleteObjectsQuiet(true),
		}
		_, err = client.DeleteObjects(keys, ossDeleteOptions...)

		if !res.IsTruncated {
			break
		}
		continuationToken = res.NextContinuationToken
	}

	shareCache.Delete(name)
	return nil
}
