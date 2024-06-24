package v1

import (
	"bytes"
	"fmt"
	file_managerv1 "github.com/NotFound1911/filestore/api/proto/gen/file_manager/v1"
	ldi "github.com/NotFound1911/filestore/internal/logger/di"
	mdi "github.com/NotFound1911/filestore/internal/mq/di"
	sdi "github.com/NotFound1911/filestore/internal/storage/di"
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
	"path"
	"strconv"
	"time"
)

type Handler struct {
	jwt.Handler
	service      service.UploadService
	fsClient     file_managerv1.FileManagerServiceClient
	storage      sdi.CustomStorage
	messageQueue mdi.MessageQueue
	logger       ldi.Logger
}

type HandlerOption func(handler *Handler)

type DiHandler struct {
	Storage      sdi.CustomStorage
	MessageQueue mdi.MessageQueue
	Logger       ldi.Logger
}

func NewHandler(service service.UploadService, hdl jwt.Handler,
	fsClient file_managerv1.FileManagerServiceClient,
	diHandler DiHandler, opts ...HandlerOption) *Handler {
	handler := &Handler{
		service:      service,
		Handler:      hdl,
		fsClient:     fsClient,
		storage:      diHandler.Storage,
		messageQueue: diHandler.MessageQueue,
		logger:       diHandler.Logger,
	}
	for _, opt := range opts {
		opt(handler)
	}
	return handler
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
		h.logger.Error(fmt.Sprintf("Failed to get file data, err:%s\n", err.Error()))
		return serv.Result{
			Code: -1,
			Msg:  fmt.Sprintf("获取文件失败:%v", err),
		}, err
	}
	// 3. 构建文件元信息
	now := time.Now()
	upInfo := domain.Upload{
		UId:      uc.UId,
		FileName: head.Filename,
		FileSha1: util.Sha1(buf.Bytes()),
		FileSize: int64(len(buf.Bytes())),
		CreateAt: &now,
		Status:   domain.UploadStatusStart,
		Type:     domain.UploadTypeSingle,
	}
	id, err := h.service.Upload(ctx, upInfo)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Failed to upload, err:%s\n", err.Error()))
		return serv.Result{
			Code: -1,
			Msg:  fmt.Sprintf("创建文件失败:%v", err),
		}, err
	}
	// 4. 将文件写入临时存储位置
	// todo 配置文件
	location := "./tmp/" + upInfo.FileSha1 + "." + head.Filename
	newFile, err := os.Create(location)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Failed to create file, err:%s\n", err.Error()))
		return serv.Result{
			Code: -1,
			Msg:  fmt.Sprintf("创建文件失败:%v", err),
		}, err
	}
	defer newFile.Close()
	nByte, err := newFile.Write(buf.Bytes())
	if int64(nByte) != upInfo.FileSize || err != nil {
		h.logger.Error(fmt.Sprintf("Failed to save data into file, writtenSize:%d, err:%s\n", nByte, err.Error()))
		return serv.Result{
			Code: -1,
			Msg:  fmt.Sprintf("文件写入失败,需要字节数:%d,err:%s", nByte, err.Error()),
		}, err
	}
	// 5. 同步或异步将文件转移到自定义对象存储
	storageMetaInfo, err := sdi.GetMetaInfo(location, upInfo.FileSha1)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Failed to get storage meta info,err:%s\n", err.Error()))
		return serv.Result{
			Code: -1,
			Msg:  fmt.Sprintf("获取文件存储元数据失败, err:%s", err.Error()),
		}, err
	}
	if !h.messageQueue.Enable() { // 同步写入
		if err = h.storage.PutObject(storageMetaInfo.Bucket, storageMetaInfo.StorageName, location, ""); err != nil {
			return serv.Result{
				Code: -1,
				Msg:  fmt.Sprintf("文件存储失败, err:%s", err.Error()),
			}, err
		}
	} else {
		// 消息队列 异步写入
		msg := mdi.Message{
			Topic: mdi.TopicName,
			Value: nil,
			Headers: []mdi.Header{
				{
					Key:   mdi.HeaderBucket,
					Value: storageMetaInfo.Bucket,
				},
				{
					Key:   mdi.HeaderStorageName,
					Value: storageMetaInfo.StorageName,
				},
				{
					Key:   mdi.HeaderLocation,
					Value: location,
				},
			},
		}
		err := h.messageQueue.SendMessage(&msg)
		if err != nil {
			return serv.Result{
				Code: -1,
				Msg:  fmt.Sprintf("消息队列发送失败, err:%s", err.Error()),
			}, err
		}
	}
	// 6.  更新文件表记录
	_, err = h.fsClient.InsertIfNotExistFileMeta(ctx.Request.Context(),
		&file_managerv1.InsertIfNotExistFileMetaReq{
			FileMeta: &file_managerv1.FileMeta{
				Sha1:        upInfo.FileSha1,
				Size:        upInfo.FileSize,
				Address:     location,
				Type:        domain.UploadTypeSingle,
				Bucket:      storageMetaInfo.Bucket,
				StorageName: storageMetaInfo.StorageName,
			},
		})
	if err != nil {
		h.logger.Error(fmt.Sprintf("Failed to insert file upInfo, err:%s\n", err.Error()))
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
				UserId:   uc.UId,
				FileName: head.Filename,
				FileSha1: upInfo.FileSha1,
				FileSize: upInfo.FileSize,
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
	err = h.service.UpdateStatus(ctx, id, domain.UploadStatusFinished)
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
				UserId:   uc.UId,
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

// InitMultiUploadFile 初始化分块上传
func (h *Handler) InitMultiUploadFile(ctx *gin.Context, req InitMultiUploadFileReq, uc jwt.UserClaims) (serv.Result, error) {
	// 初始化分块上传的初始化信息
	now := time.Now()
	uploadInfo := domain.Upload{
		UId:      uc.UId,
		FileName: req.FileName,
		FileSha1: req.FileSha1,
		FileSize: req.FileSize,
		CreateAt: &now,
		Status:   domain.UploadStatusStart,
		Type:     domain.UploadTypeMulti,
	}
	id, err := h.service.Upload(ctx, uploadInfo)
	if err != nil {
		log.Printf("Failed to upload, err:%s\n", err.Error())
		return serv.Result{
			Code: -1,
			Msg:  fmt.Sprintf("创建文件失败:%v", err),
		}, err
	}
	uploadInfo.Id = id
	return serv.Result{
		Code: 2000,
		Msg:  "OK",
		Data: uploadInfo,
	}, nil
}

// MultiUploadFilePart 分块传输
func (h *Handler) MultiUploadFilePart(ctx *gin.Context, uc jwt.UserClaims) (serv.Result, error) {
	chunkIdStr := ctx.Query("chunk_id")
	chunkId, err := strconv.ParseInt(chunkIdStr, 10, 64)
	if err != nil {
		fmt.Println("Error converting string to int:", err)
		return serv.Result{
			Code: -1,
			Msg:  "chunk id 错误",
		}, err
	}
	uploadIdStr := ctx.Query("upload_id")
	uploadId, err := strconv.ParseInt(uploadIdStr, 10, 64)
	if err != nil {
		fmt.Println("Error converting string to int64:", err)
		return serv.Result{
			Code: -1,
			Msg:  "upload id 错误",
		}, err
	}
	fileName := ctx.Query("file_name")
	chunkSha1 := ctx.Query("chunk_sha1")
	chunkSizeStr := ctx.Query("chunk_size")
	chunkSize, err := strconv.ParseInt(chunkSizeStr, 10, 64)
	if err != nil {
		fmt.Println("Error converting string to int64:", err)
		return serv.Result{
			Code: -1,
			Msg:  "upload id 错误",
		}, err
	}
	countStr := ctx.Query("count")
	count, err := strconv.ParseInt(countStr, 10, 64)
	if err != nil {
		fmt.Println("Error converting string to int:", err)
		return serv.Result{
			Code: -1,
			Msg:  "count 错误",
		}, err
	}

	// todo 配置
	fPath := fmt.Sprintf("./tmp/%d/%d", uploadId, chunkId)
	os.MkdirAll(path.Dir(fPath), 0744)
	fd, err := os.Create(fPath)
	if err != nil {
		return serv.Result{
			Code: -1,
			Msg:  fmt.Sprintf("文件创建失败:%v", err),
		}, err
	}
	defer fd.Close()
	now := time.Now()
	chunk := domain.Chunk{
		Id:       chunkId,
		UploadId: uploadId,
		UId:      uc.UId,
		FileName: fileName,
		Sha1:     chunkSha1,
		Size:     chunkSize,
		CreateAt: &now,
		UpdateAt: &now,
		Status:   domain.UploadStatusStart,
		Count:    count,
	}
	if err := h.service.SetFileChunk(ctx, chunk); err != nil {
		log.Printf("set file chunk failed,err:%v", err)
	}
	file, err := ctx.FormFile("file")
	if err != nil {
		return serv.Result{
			Code: -1,
			Msg:  fmt.Sprintf("获取文件失败:%v", err),
		}, err
	}
	src, err := file.Open()
	if err != nil {
		return serv.Result{
			Code: -1,
			Msg:  fmt.Sprintf("文件打开失败:%v", err),
		}, err
	}
	if _, err := io.Copy(fd, src); err != nil {
		return serv.Result{
			Code: -1,
			Msg:  fmt.Sprintf("文件存储失败:%v", err),
		}, err
	}
	// check sha1
	sha1, err := util.GetFileSha1(fPath)
	if err != nil {
		return serv.Result{
			Code: -1,
			Msg:  fmt.Sprintf("%s生成md5 失败,err:%v", fPath, err),
		}, err
	}
	if sha1 != chunkSha1 {
		return serv.Result{
			Code: -1,
			Msg:  fmt.Sprintf("sha1 错误,期望：%s,实际:%s", chunkSha1, sha1),
		}, nil
	}
	// 更新redis缓存状态
	now = time.Now()
	chunk.UpdateAt = &now
	chunk.Status = domain.UploadStatusFinished
	if err := h.service.SetFileChunk(ctx, chunk); err != nil {
		log.Printf("set file chunk failed,err:%v", err)
	}
	return serv.Result{
		Code: 2000,
		Msg:  "OK",
	}, nil
}

// MultiUploadFileMerge 通知分块上传完成，合并
func (h *Handler) MultiUploadFileMerge(ctx *gin.Context, req MultiUploadFileMergeReq, uc jwt.UserClaims) (serv.Result, error) {
	cs, err := h.service.GetChunks(ctx, req.UploadId)
	if err != nil {
		return serv.Result{
			Code: -1,
			Msg:  fmt.Sprintf("获取文件分片信息失败:%v", err),
		}, err
	}
	// 校验数量
	cnt := cs[0].Count
	if cnt != int64(len(cs)) {
		return serv.Result{
			Code: -1,
			Msg:  fmt.Sprintf("分片信息不完整，期望:%d,实际:%d", cnt, len(cs)),
		}, nil
	}
	// 校验状态
	for i := 0; int64(i) < cnt; i++ {
		if cs[i].Status != domain.UploadStatusFinished {
			return serv.Result{
				Code: -1,
				Msg:  "切片状态未完成",
			}, nil
		}
	}
	// TODO：合并分块, 可以将ceph当临时存储，合并时将文件写入ceph;
	// 也可以不用在本地进行合并，转移的时候将分块append到ceph/oss即可
	srcPath := fmt.Sprintf("./tmp/%d/", req.UploadId)
	filePaths, err := util.GetAllFilesInDirectory(srcPath)
	if err != nil {
		return serv.Result{
			Code: -1,
			Msg:  fmt.Sprintf("获取文件列表失败，err:%v", err),
		}, err
	}
	location := fmt.Sprintf("./tmp/%s.%s", req.FileSha1, req.FileName) // todo 文件bucket
	if err = util.Merge(filePaths, location); err != nil {
		return serv.Result{
			Code: -1,
			Msg:  fmt.Sprintf("文件合并失败:%v", err),
		}, err
	}
	// sha1 校验
	sha1, err := util.GetFileSha1(location)
	if err != nil {
		return serv.Result{
			Code: -1,
			Msg:  fmt.Sprintf("文件sha1获取失败:%v", err),
		}, err
	}
	if sha1 != req.FileSha1 {
		return serv.Result{
			Code: -1,
			Msg:  fmt.Sprintf("文件sha1校验失败，期望:%s，实际:%s", req.FileSha1, sha1),
		}, nil
	}
	// 5. 同步或异步将文件转移到自定义对象存储

	storageMetaInfo, err := sdi.GetMetaInfo(location, req.FileSha1)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Failed to get storage meta info,err:%s\n", err.Error()))
		return serv.Result{
			Code: -1,
			Msg:  fmt.Sprintf("获取文件存储元数据失败, err:%s", err.Error()),
		}, err
	}
	if !h.messageQueue.Enable() { // 同步写入
		h.storage.PutObject(storageMetaInfo.Bucket, storageMetaInfo.StorageName, location, "")
	} else {
		// 消息队列 异步写入
		msg := mdi.Message{
			Topic: mdi.TopicName,
			Value: nil,
			Headers: []mdi.Header{
				{
					Key:   mdi.HeaderBucket,
					Value: storageMetaInfo.Bucket,
				},
				{
					Key:   mdi.HeaderStorageName,
					Value: storageMetaInfo.StorageName,
				},
				{
					Key:   mdi.HeaderLocation,
					Value: location,
				},
			},
		}
		err := h.messageQueue.SendMessage(&msg)
		if err != nil {
			return serv.Result{
				Code: -1,
				Msg:  fmt.Sprintf("消息队列发送失败, err:%s", err.Error()),
			}, err
		}
	}
	// 更新唯一文件表及用户文件表
	err = h.service.UpdateStatus(ctx, req.UploadId, domain.UploadStatusFinished)
	if err != nil {
		h.logger.Warn(fmt.Sprintf("文件任务列表更新失败:%v\n", err))
	}
	_, err = h.fsClient.InsertIfNotExistFileMeta(ctx.Request.Context(),
		&file_managerv1.InsertIfNotExistFileMetaReq{
			FileMeta: &file_managerv1.FileMeta{
				Sha1:    req.FileSha1,
				Size:    req.FileSize,
				Address: location, // todo 维护
				Type:    domain.UploadTypeSingle,
			},
		})
	if err != nil {
		h.logger.Error(fmt.Sprintf("Failed to insert file meta, err:%s\n", err.Error()))
		return serv.Result{
			Code: -1,
			Msg:  fmt.Sprintf("创建文件元数据失败:%v", err),
		}, err
	}
	upTime := time.Now()
	t := timestamppb.New(upTime)
	_, err = h.fsClient.InsertUserFile(ctx.Request.Context(),
		&file_managerv1.InsertUserFileReq{
			UserFile: &file_managerv1.UserFile{
				UserId:   uc.UId,
				FileName: req.FileName,
				FileSha1: req.FileSha1,
				FileSize: req.FileSize,
				UpdateAt: t,
			},
		},
	)
	if err != nil {
		return serv.Result{
			Code: -1,
			Msg:  fmt.Sprintf("用户文件列表更新失败:%v", err),
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
	ug.POST("/init-multi-upload-file", serv.WrapBodyAndClaims(h.InitMultiUploadFile))
	ug.POST("/multi-upload-file-part", serv.WrapClaims(h.MultiUploadFilePart))
	ug.POST("/multi-upload-file-merge", serv.WrapBodyAndClaims(h.MultiUploadFileMerge))
}
