package local

import (
	"fmt"
	"github.com/NotFound1911/filestore/config"
	"github.com/NotFound1911/filestore/internal/storage/di"
	"io"
	"os"
	"path"
	"path/filepath"
)

const (
	localStorage = di.LocalStorage
	Name         = "LOCAL"
)

// Storage 本地存储
type Storage struct {
	RootPath string
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
		fmt.Println("Failed to open file:", err)
		return nil, err
	}
	defer file.Close()
	_, err = file.Seek(offset, io.SeekStart)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	buffer := make([]byte, length)
	_, err = file.Read(buffer)
	if err != nil && err != io.EOF {
		fmt.Println("Error:", err)
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
			fmt.Println("Failed to create folder:", err)
			return err
		}
		fmt.Println("Folder created successfully.")
	}
	file, err := os.Create(objectPath)
	if err != nil {
		fmt.Println("Failed to create file:", err)
		return err
	}
	defer file.Close()

	// 复制文件内容
	_, err = io.Copy(file, sourceFile)
	if err != nil {
		fmt.Println("Failed to copy file:", err)
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

func NewStorage(conf *config.Configuration, opts ...StorageOption) di.CustomStorage {
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
