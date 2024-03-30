package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
)

func UploadView(c *gin.Context) {
	data, err := os.ReadFile("./static/view/upload.html")
	if err != nil {
		c.JSON(http.StatusInternalServerError, "文件读取错误")
		return
	}
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, string(data))
}
func UploadFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusInternalServerError, "获取文件失败")
		return
	}
	newFile, err := os.Create("./tmp/" + header.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("创建文件失败, 详情:%v", err))
		return
	}
	defer newFile.Close()
	_, err = io.Copy(newFile, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "文件写入失败")
		return
	}
	//c.JSON(http.StatusOK, "上传成功")
	c.Redirect(http.StatusMovedPermanently, "/api/storage/v1/upload/successful")
}
func UploadSuccessful(c *gin.Context) {
	c.JSON(http.StatusOK, "上传成功")
}
