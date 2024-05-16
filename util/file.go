package util

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Merge 合并多个文件
func Merge(filePaths []string, outputPath string) error {
	outputFile, err := os.Create(outputPath)
	if err != nil {
		fmt.Println("create failed: ", err)
		return err
	}
	defer outputFile.Close()

	for i, filePath := range filePaths {
		err = func() error {
			fmt.Printf("第 %d 个临时文件路径为：%s.\n", i+1, filePath)
			inputFile, openErr := os.Open(filePath)
			if openErr != nil {
				fmt.Println("打开失败：", openErr)
				return err
			}
			defer inputFile.Close()

			_, copyErr := io.Copy(outputFile, inputFile)
			if copyErr != nil {
				fmt.Println("copy 失败：", copyErr)
				return err
			}
			return nil
		}()
		if err != nil {
			return err
		}
	}
	return nil
}
func GetAllFilesInDirectory(dirPath string) ([]string, error) {
	var files []string

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 如果当前路径是文件而不是目录，则将其添加到文件列表中
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}
