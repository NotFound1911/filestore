package grpc

import (
	"context"
	file_managerv1 "github.com/NotFound1911/filestore/api/proto/gen/file_manager/v1"
	"github.com/NotFound1911/filestore/service/file_manager/domain"
	"github.com/NotFound1911/filestore/service/file_manager/service"
)

type FileManagerServiceServer struct {
	file_managerv1.UnimplementedFileManagerServiceServer
	svc service.FileManagerService
}

func (f *FileManagerServiceServer) InsertIfNotExistFileMeta(ctx context.Context, in *file_managerv1.InsertIfNotExistFileMetaReq) (*file_managerv1.InsertIfNotExistFileMetaResp, error) {
	id, err := f.svc.InsertFileMeta(ctx, f.toFileMetaDomain(in))
	return &file_managerv1.InsertIfNotExistFileMetaResp{Id: id}, err
}

func (f *FileManagerServiceServer) InsertUserFile(ctx context.Context, in *file_managerv1.InsertUserFileReq) (*file_managerv1.InsertUserFileResp, error) {
	id, err := f.svc.InsertUserFile(ctx, f.toUserFileDomain(in))
	return &file_managerv1.InsertUserFileResp{Id: id}, err
}

func (f *FileManagerServiceServer) GetFileMetaByUserId(ctx context.Context, in *file_managerv1.GetFileMetaByUserIdReq) (*file_managerv1.GetFileMetaByUserIdResp, error) {
	tmp, err := f.svc.GetFileMetaByUserId(ctx, in.GetUid())
	res := make([]*file_managerv1.FileMeta, 0, len(tmp))
	for _, v := range tmp {
		meta := f.toFileMetaProto(v)
		res = append(res, &meta)
	}
	return &file_managerv1.GetFileMetaByUserIdResp{FileMeta: res}, err
}
func (f *FileManagerServiceServer) GetFileMeta(ctx context.Context, in *file_managerv1.GetFileMetaReq) (*file_managerv1.GetFileMetaResp, error) {
	res, err := f.svc.GetFileMeta(ctx, in.GetFileSha1())
	meta := f.toFileMetaProto(res)
	return &file_managerv1.GetFileMetaResp{
		FileMeta: &meta,
	}, err
}
func NewFileManagerServiceServer(svc service.FileManagerService) *FileManagerServiceServer {
	return &FileManagerServiceServer{svc: svc}
}

func (f *FileManagerServiceServer) toFileMetaDomain(in *file_managerv1.InsertIfNotExistFileMetaReq) domain.FileMeta {
	return domain.FileMeta{
		Sha1:    in.FileMeta.Sha1,
		Size:    in.FileMeta.Size,
		Address: in.FileMeta.Address,
		Type:    in.FileMeta.Type,
	}
}
func (f *FileManagerServiceServer) toUserFileDomain(in *file_managerv1.InsertUserFileReq) domain.UserFile {
	updateTimestamp := in.UserFile.UpdateAt.AsTime()
	return domain.UserFile{
		UserId:   in.UserFile.UserId,
		FileSize: in.UserFile.FileSize,
		FileName: in.UserFile.FileName,
		FileSha1: in.UserFile.FileSha1,
		UpdateAt: &updateTimestamp,
	}
}
func (f *FileManagerServiceServer) toFileMetaProto(meta domain.FileMeta) file_managerv1.FileMeta {
	return file_managerv1.FileMeta{
		Size:    meta.Size,
		Sha1:    meta.Sha1,
		Address: meta.Address,
		Type:    meta.Type,
	}
}
