package local

import (
	"fmt"
	"github.com/NotFound1911/filestore/config"
	ldi "github.com/NotFound1911/filestore/internal/logger/di"
	sdi "github.com/NotFound1911/filestore/internal/storage/di"
	"io"
	"os"
	"path"
	"path/filepath"
)

const (
	localStorage = sdi.LocalStorage
	Name         = "LOCAL"
)

// Storage 本地存储
type Storage struct {
	RootPath string
	logger   ldi.Logger
}

func (s *Storage) MakeBucket(bucketName string) error {
	dirName := path.Join(s.RootPath, bucketName)
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		if err := os.MkdirAll(dirName, 0755); err != nil {
			return err
		}
	}
	return nil
}

func (s *Storage) GetObject(bucketName, objectName string, offset, length int64) ([]byte, error) {
	objectPath := path.Join(s.RootPath, bucketName, objectName)
	file, err := os.Open(objectPath)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed to open file:%v", err))
		return nil, err
	}
	defer file.Close()
	_, err = file.Seek(offset, io.SeekStart)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Error:%v", err))
		return nil, err
	}
	buffer := make([]byte, length)
	_, err = file.Read(buffer)
	if err != nil && err != io.EOF {
		s.logger.Error(fmt.Sprintf("Read Error:%v", err))
		return nil, err
	}
	return buffer, nil
}

func (s *Storage) PutObject(bucketName, objectName, filePath, contentType string) error {
	// copy 数据到 具体的目录
	// 打开源文件
	sourceFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	objectPath := path.Join(s.RootPath, bucketName, objectName)
	// 新建文件夹
	// 定义新文件夹的路径
	folderPath := filepath.Dir(objectPath)
	// 检查文件夹是否已经存在
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		// 文件夹不存在，创建新文件夹
		err := os.Mkdir(folderPath, 0755) // 0755 权限表示用户具有读、写、执行权限，组和其他用户具有读、执行权限
		if err != nil {
			s.logger.Error(fmt.Sprintf("Failed to create folder:%v", err))
			return err
		}
		s.logger.Info(fmt.Sprintf("Folder:%s created successfully.", folderPath))
	}
	file, err := os.Create(objectPath)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed to create file:%v", err))
		return err
	}
	defer file.Close()

	// 复制文件内容
	_, err = io.Copy(file, sourceFile)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed to copy file:%v", err))
		return err
	}
	return nil
}

func (s *Storage) DeleteObject(bucketName, objectName string) error {
	objectPath := path.Join(s.RootPath, bucketName, objectName)
	err := os.RemoveAll(objectPath)
	return err
}

func (s *Storage) Type() string {
	return Name
}

type StorageOption func(storage *Storage)

func NewStorage(conf *config.Configuration, opts ...StorageOption) sdi.CustomStorage {
	s := &Storage{
		RootPath: localStorage,
	}
	if conf.Storage.Local.Dir != "" {
		s.RootPath = conf.Storage.Local.Dir
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}
func WithLogger(logger ldi.Logger) StorageOption {
	return func(storage *Storage) {
		storage.logger = logger
	}
}
