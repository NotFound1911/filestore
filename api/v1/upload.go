package v1

import (
	"fmt"
	"github.com/NotFound1911/filestore/meta"
	"github.com/NotFound1911/filestore/util"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"time"
)

// UploadView 上传界面
func UploadView(c *gin.Context) {
	data, err := os.ReadFile("./static/view/upload.html")
	if err != nil {
		c.JSON(http.StatusInternalServerError, "文件读取错误")
		return
	}
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, string(data))
}

// UploadFile 上传文件
func UploadFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusInternalServerError, "获取文件失败")
		return
	}
	defer file.Close()

	fileMeta := meta.File{
		Name:     header.Filename,
		Location: "./tmp/" + header.Filename,
		UploadAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	newFile, err := os.Create(fileMeta.Location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("创建文件失败, 详情:%v", err))
		return
	}
	defer newFile.Close()
	fileMeta.Size, err = io.Copy(newFile, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "文件写入失败")
		return
	}
	newFile.Seek(0, 0)
	fileMeta.Sha1 = util.FileSha1(newFile)
	fmt.Println("fileMeta:", fileMeta)
	meta.UpdateFileMeta(fileMeta)
	c.Redirect(http.StatusMovedPermanently, "/api/storage/v1/file/upload/successful")
}
func UploadSuccessful(c *gin.Context) {
	c.JSON(http.StatusOK, "上传成功")
}

func GetFileMeta(c *gin.Context) {
	filehash := c.PostForm("filehash")
	// todo 存在性校验
	fMeta := meta.GetFileMeta(filehash)
	c.JSON(http.StatusOK, fMeta)
}

func DownLoad(c *gin.Context) {
	filehash := c.PostForm("filehash")
	fMeta := meta.GetFileMeta(filehash)
	//c.Writer.Header().Set("Content-Type", "application/octect-stream")
	fmt.Println("name:", fMeta.Name)
	c.Writer.Header().Set("Content-Disposition", "attachment;filename=\""+fMeta.Name+"\"")
	c.File(fMeta.Location)
}

func UpdateFile(c *gin.Context) {
	// todo 参数校验
	filehash := c.Query("filehash")
	opType := c.Query("op")
	filename := c.Query("filename")
	if opType != "0" {
		c.JSON(http.StatusForbidden, "操作类型错误")
		return
	}
	fMeta := meta.GetFileMeta(filehash)
	fMeta.Name = filename
	meta.UpdateFileMeta(fMeta)
	c.JSON(http.StatusOK, fMeta)
}

func DeleteFile(c *gin.Context) {
	filehash := c.PostForm("filehash")
	fMeta := meta.GetFileMeta(filehash)
	os.Remove(fMeta.Location)
	meta.DeleteFileMeta(filehash)
	c.JSON(http.StatusOK, "删除成功")
}
