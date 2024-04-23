package v1

import (
	"bytes"
	"fmt"
	file_managerv1 "github.com/NotFound1911/filestore/api/proto/gen/file_manager/v1"
	"github.com/NotFound1911/filestore/internal/web/jwt"
	serv "github.com/NotFound1911/filestore/pkg/server"
	"github.com/NotFound1911/filestore/service/upload/domain"
	"github.com/NotFound1911/filestore/service/upload/service"
	"github.com/NotFound1911/filestore/util"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"log"
	"os"
	"time"
)

type Handler struct {
	jwt.Handler
	service  service.UploadService
	fsClient file_managerv1.FileManagerServiceClient
}

func NewHandler(service service.UploadService, hdl jwt.Handler, fsClient file_managerv1.FileManagerServiceClient) *Handler {
	return &Handler{
		service:  service,
		Handler:  hdl,
		fsClient: fsClient,
	}
}

// UploadFile 文件上传
func (h *Handler) UploadFile(ctx *gin.Context, uc jwt.UserClaims) (serv.Result, error) {
	// 1. 从form表单中获得文件内容句柄
	file, head, err := ctx.Request.FormFile("file")
	if err != nil {
		log.Printf("Failed to get form data, err:%s\n", err.Error())
		return serv.Result{
			Code: -1,
			Msg:  fmt.Sprintf("获取form data失败:%v", err),
		}, err
	}
	defer file.Close()
	// 2. 把文件内容转为[]byte
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		log.Printf("Failed to get file data, err:%s\n", err.Error())
		return serv.Result{
			Code: -1,
			Msg:  fmt.Sprintf("获取文件失败:%v", err),
		}, err
	}
	// 3. 构建文件元信息
	now := time.Now()
	// todo 构造元数据
	meta := domain.Upload{
		FileName: head.Filename,
		FileSha1: util.Sha1(buf.Bytes()),
		FileSize: int64(len(buf.Bytes())),
		CreateAt: &now,
		Status:   "开始上传",
	}
	id, err := h.service.Upload(ctx, meta)
	if err != nil {
		log.Printf("Failed to upload, err:%s\n", err.Error())
		return serv.Result{
			Code: -1,
			Msg:  fmt.Sprintf("创建文件失败:%v", err),
		}, err
	}
	// 4. 将文件写入临时存储位置
	// toto 配置文件
	location := "./tmp/" + meta.FileSha1
	newFile, err := os.Create(location)
	if err != nil {
		log.Printf("Failed to create file, err:%s\n", err.Error())
		return serv.Result{
			Code: -1,
			Msg:  fmt.Sprintf("创建文件失败:%v", err),
		}, err
	}
	defer newFile.Close()
	nByte, err := newFile.Write(buf.Bytes())
	if int64(nByte) != meta.FileSize || err != nil {
		log.Printf("Failed to save data into file, writtenSize:%d, err:%s\n", nByte, err.Error())
		return serv.Result{
			Code: -1,
			Msg:  fmt.Sprintf("文件写入失败,需要字节数:%d,err:%s", nByte, err.Error()),
		}, err
	}
	// 5. 同步或异步将文件转移到Ceph/OSS
	// todo
	// 6.  更新文件表记录
	// todo 元数据更新
	// 元数据更新  妙传判断
	_, err = h.fsClient.InsertIfNotExistFileMeta(ctx.Request.Context(),
		&file_managerv1.InsertIfNotExistFileMetaReq{
			FileMeta: &file_managerv1.FileMeta{
				Sha1:    meta.FileSha1,
				Size:    meta.FileSize,
				Address: location,
				Type:    "upload-test",
			},
		})
	if err != nil {
		log.Printf("Failed to insert file meta, err:%s\n", err.Error())
		return serv.Result{
			Code: -1,
			Msg:  fmt.Sprintf("创建文件元数据失败:%v", err),
		}, err
	}
	// 更新用户文件列表
	upTime := time.Now()
	t := timestamppb.New(upTime)
	_, err = h.fsClient.InsertUserFile(ctx.Request.Context(),
		&file_managerv1.InsertUserFileReq{
			UserFile: &file_managerv1.UserFile{
				UserId:   uc.Uid,
				FileName: head.Filename,
				FileSha1: meta.FileSha1,
				FileSize: meta.FileSize,
				UpdateAt: t,
			},
		},
	)
	if err != nil {
		log.Printf("Failed to insert user file, err:%s\n", err.Error())
		return serv.Result{
			Code: -1,
			Msg:  fmt.Sprintf("创建用户文件失败:%v", err),
		}, err
	}
	// 用户文件更新
	err = h.service.UpdateStatus(ctx, id, "上传完毕")
	if err != nil {
		return serv.Result{
			Code: -1,
			Msg:  fmt.Sprintf("文件状态更新失败,err:%s", err.Error()),
		}, err
	}
	return serv.Result{
		Code: 2000,
		Msg:  "OK",
	}, nil
}

// Resume 秒传
func (h *Handler) Resume(ctx *gin.Context, req ResumeReq, uc jwt.UserClaims) (serv.Result, error) {
	// 1. 解析请求参数
	fileSha1 := req.FileSha1
	fileName := req.FileName

	// 2. 从文件表中查询相同hash的文件记录
	fileMeaResp, err := h.fsClient.GetFileMeta(ctx, &file_managerv1.GetFileMetaReq{
		FileSha1: fileSha1,
	})
	if err != nil {
		log.Printf("get file meta by sha1 failed,err:%v\n", err)
		return serv.Result{
			Code: -1,
			Msg:  fmt.Sprintf("文件元数据查询失败,err:%s", err.Error()),
		}, err
	}
	// 3. 查不到记录则返回秒传失败
	if fileMeaResp.GetFileMeta().Size == 0 {
		return serv.Result{
			Code: -1,
			Msg:  "无上传记录,秒传失败",
		}, err
	}
	// 4. 上传过则将文件信息写入用户文件表， 返回成功
	upTime := time.Now()
	t := timestamppb.New(upTime)
	_, err = h.fsClient.InsertUserFile(ctx.Request.Context(),
		&file_managerv1.InsertUserFileReq{
			UserFile: &file_managerv1.UserFile{
				UserId:   uc.Uid,
				FileName: fileName,
				FileSha1: fileSha1,
				FileSize: fileMeaResp.GetFileMeta().GetSize(),
				UpdateAt: t,
			},
		},
	)
	if err != nil {
		log.Printf("insert user file failed,err:%v\n", err)
		return serv.Result{
			Code: -1,
			Msg:  fmt.Sprintf("用户文件上传失败,err:%s", err.Error()),
		}, err
	}
	return serv.Result{
		Code: 2000,
		Msg:  "OK",
	}, nil
}
func (h *Handler) RegisterUploadRoutes(core *gin.Engine) {
	ug := core.Group("/api/storage/v1/upload")
	ug.POST("/upload-file", serv.WrapClaims(h.UploadFile))
	ug.POST("/resume", serv.WrapBodyAndClaims(h.Resume))
}
