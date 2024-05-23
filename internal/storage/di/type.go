package di

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const LocalStorage = "local_storage"

// CustomStorage 存储
type CustomStorage interface {
	// MakeBucket 创建存储桶
	MakeBucket(bucketName string) error

	// GetObject 获取存储对象
	GetObject(bucketName, objectName string, offset, length int64) ([]byte, error)

	// PutObject 上传存储对象
	PutObject(bucketName, objectName, filePath, contentType string) error

	// DeleteObject 删除存储对象
	DeleteObject(bucketName, objectName string) error

	// Type 类别
	Type() string
}

func GetExtension(filename string) string {
	ext := filepath.Ext(filename)
	if ext == "" {
		return ""
	}
	return strings.ToLower(ext[1:])
}

// selectBucketBySuffix .
func selectBucketBySuffix(filename string) string {
	suffix := GetExtension(filename)
	if suffix == "" {
		return ""
	}
	switch suffix {
	case "jpg", "jpeg", "png", "gif", "bmp":
		return "image"
	case "mp4", "avi", "wmv", "mpeg":
		return "video"
	case "mp3", "wav", "flac":
		return "audio"
	case "pdf", "doc", "docx", "ppt", "pptx", "xls", "xlsx":
		return "doc"
	case "zip", "rar", "tar", "gz", "7z":
		return "archive"
	default:
		return "unknown"
	}
}

type BucketInfo struct {
	Bucket      string
	StorageName string
}

// GetMetaInfo 根据文件名 获取文件相关元数据信息
func GetMetaInfo(filename, uidStr string) (*BucketInfo, error) {
	bucket := selectBucketBySuffix(filename)
	name := filepath.Base(filename)
	name = url.PathEscape(name)
	storageName := fmt.Sprintf("%s.%s", uidStr, GetExtension(filename))
	// 在本地创建uid的目录
	if err := os.MkdirAll(path.Join(LocalStorage, uidStr), 0755); err != nil {
		return nil, err
	}
	return &BucketInfo{
		Bucket:      bucket,
		StorageName: storageName,
	}, nil
}
