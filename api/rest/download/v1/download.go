package v1

import (
	"fmt"
	file_managerv1 "github.com/NotFound1911/filestore/api/proto/gen/file_manager/v1"
	ldi "github.com/NotFound1911/filestore/internal/logger/di"
	"github.com/NotFound1911/filestore/internal/web/jwt"
	serv "github.com/NotFound1911/filestore/pkg/server"
	"github.com/NotFound1911/filestore/service/download/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

type Handler struct {
	jwt.Handler
	logger   ldi.Logger
	fsClient file_managerv1.FileManagerServiceClient
	service  service.DownloadService
}

func NewHandler(service service.DownloadService, hdl jwt.Handler, fsClient file_managerv1.FileManagerServiceClient, logger ldi.Logger) *Handler {
	return &Handler{
		Handler:  hdl,
		logger:   logger,
		fsClient: fsClient,
		service:  service,
	}
}

// DownloadURLHandler 生成下载链接
func (h *Handler) DownloadURLHandler(ctx *gin.Context, req DownloadURLHandlerReq, uc jwt.UserClaims) (serv.Result, error) {
	fileHash := req.FileHash
	uId := uc.UId
	// 查询文件元数据
	res, err := h.fsClient.GetFileMeta(ctx, &file_managerv1.GetFileMetaReq{FileSha1: fileHash})
	if err != nil {
		h.logger.Error(fmt.Sprintf("%v 获取元数据:%s失败", uId, fileHash))
		return serv.Result{
			Code: -1,
			Msg:  fmt.Sprintf("元数据获取失败"),
		}, err
	}
	fmt.Printf("res:%v\n", res)
	// todo
	uri := fmt.Sprintf("bucket=%s&sha1=%s&name=%s&address=%s",
		res.GetFileMeta().Bucket,
		res.GetFileMeta().Sha1,
		res.GetFileMeta().StorageName,
		res.GetFileMeta().Address,
	)
	return serv.Result{
		Code: 2000,
		Msg:  "OK",
		Data: uri,
	}, nil
}

// Download 文件下载
func (h *Handler) Download(ctx *gin.Context) {
	// todo
	// 获取文件名
	fileName := ctx.Query("filename")
	if fileName == "" {
		ctx.String(http.StatusBadRequest, "请提供文件名")
		return
	}

	// 检查文件是否存在
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		ctx.String(http.StatusNotFound, "文件不存在")
		return
	}

	// 设置响应头，告诉浏览器这是一个要下载的文件
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.File(fileName)
}

func (h *Handler) RegisterDownloadRoutes(core *gin.Engine) {
	dl := core.Group("/api/storage/v1/download")
	dl.POST("/dDownload-url", serv.WrapBodyAndClaims(h.DownloadURLHandler))
	dl.POST("/download", h.Download)
}
