package oss

import (
	"context"
	"github.com/jingbh/simple-share/internal/models"
	"sync"
	"time"
)

type shareCacheRecord struct {
	Share *models.Share
	Time  time.Time
}

var shareCache sync.Map

func GetShareCached(ctx context.Context, name string) (*models.Share, error) {
	if v, ok := shareCache.Load(name); ok {
		record := v.(*shareCacheRecord)
		if time.Since(record.Time) < time.Hour {
			return record.Share, nil
		}
	}

	res, err := GetShare(ctx, name)
	if err != nil {
		return nil, err
	}
	shareCache.Store(name, &shareCacheRecord{
		Share: res,
		Time:  time.Now(),
	})

	return res, nil
}
