package minio

import (
	"bytes"
	"context"
	"fmt"
	ldi "github.com/NotFound1911/filestore/internal/logger/di"
	sdi "github.com/NotFound1911/filestore/internal/storage/di"
	m "github.com/NotFound1911/filestore/pkg/minio"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
)

type Storage struct {
	client *minio.Client
	logger ldi.Logger
}

func (s *Storage) MakeBucket(bucketName string) error {
	ctx := context.Background()
	isExist, err := s.client.BucketExists(ctx, bucketName)
	if err != nil {
		s.logger.Error(fmt.Sprintf("BucketExists %s failed,err:%v", bucketName, err))
		return err
	}
	if isExist {
		s.logger.Debug(fmt.Sprintf("%v is exist", bucketName))
		return nil
	}
	return s.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: "cn-north-1"})
}

func (s *Storage) GetObject(bucketName, objectName string, offset, length int64) ([]byte, error) {
	ctx := context.Background()
	obj, err := s.client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	// 注意，这里需要关闭流，否则会造成资源占用，minio会hang住
	defer func(obj *minio.Object) {
		err := obj.Close()
		if err != nil {
			s.logger.Warn(fmt.Sprintf("%+v closed failed,err:%v", obj, err))
		}
	}(obj)
	stat, err := obj.Stat()
	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed to stat object %s/%s: %v", bucketName, objectName, err))
		return nil, err
	}
	s.logger.Debug(fmt.Sprintf("Object size: %d", stat.Size))
	_, err = obj.Seek(offset, io.SeekStart)
	if err != nil {
		s.logger.Error(fmt.Sprintf("%s %s seek %d failed,err:%v", bucketName, objectName, offset, err))
		return nil, err
	}
	if length == -1 {
		// 读取所有剩余内容
		buffer := bytes.NewBuffer(nil)
		tempBuffer := make([]byte, 1024)
		for {
			n, err := obj.Read(tempBuffer)
			if err != nil && err != io.EOF {
				s.logger.Error(fmt.Sprintf("Read failed: %v", err))
				return nil, err
			}
			if n == 0 {
				break
			}
			buffer.Write(tempBuffer[:n])
		}
		return buffer.Bytes(), nil
	}
	data := make([]byte, length)
	n, err := obj.Read(data) // 这里有read，就需要close
	return data[:n], err
}

func (s *Storage) PutObject(bucketName, objectName, filePath, contentType string) error {
	ctx := context.Background()
	// 上传对象到存储桶
	_, err := s.client.FPutObject(ctx, bucketName, objectName, filePath,
		minio.PutObjectOptions{ContentType: contentType, NumThreads: 10})
	if err != nil {
		s.logger.Error(fmt.Sprintf("Storage PutObject %s %s %s %s failed,err:%v",
			bucketName, objectName, filePath, contentType, err))
	}
	return err
}

func (s *Storage) DeleteObject(bucketName, objectName string) error {
	ctx := context.Background()
	// 删除存储桶中的对象
	err := s.client.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
	return err
}

func NewStorage(service *m.Service, logger ldi.Logger) sdi.CustomStorage {
	client, err := minio.New(service.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(service.AccessKeyId, service.SecretAccessKey, ""),
		Secure: service.UseSSL,
	})
	if err != nil {
		panic(err)
	}
	return &Storage{
		client: client,
		logger: logger,
	}
}
