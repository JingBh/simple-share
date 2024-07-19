package oss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/spf13/viper"
	"strings"
	"sync"
)

func newClient(internal bool) (*oss.Bucket, error) {
	endpoint := viper.GetString("oss.endpoint")
	if !internal {
		endpointPublic := viper.GetString("oss.endpoint_public")
		if endpointPublic != "" {
			endpoint = endpointPublic
		}
	}
	akId := viper.GetString("oss.access_key_id")
	akSecret := viper.GetString("oss.access_key_secret")
	cname := !strings.Contains(endpoint, "aliyuncs")

	client, err := oss.New(endpoint, akId, akSecret, oss.UseCname(cname))
	if err != nil {
		return nil, err
	}

	bucket := viper.GetString("oss.bucket")
	return client.Bucket(bucket)
}

var Client = sync.OnceValue(func() *oss.Bucket {
	bucket, err := newClient(true)
	if err != nil {
		panic(err)
	}
	return bucket
})

var PublicClient = sync.OnceValue(func() *oss.Bucket {
	bucket, err := newClient(false)
	if err != nil {
		panic(err)
	}
	return bucket
})
