package minio

import "github.com/NotFound1911/filestore/config"

type Service struct {
	Endpoint        string
	AccessKeyId     string
	SecretAccessKey string
	UseSSL          bool
}

func NewService(conf *config.Configuration) *Service {
	return &Service{
		Endpoint:        conf.Storage.Minio.Endpoint,
		AccessKeyId:     conf.Storage.Minio.AccessKeyId,
		SecretAccessKey: conf.Storage.Minio.SecretAccessKey,
		UseSSL:          conf.Storage.Minio.UseSSL,
	}
}
