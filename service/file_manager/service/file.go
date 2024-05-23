package service

import (
	"context"
	"github.com/NotFound1911/filestore/service/file_manager/domain"
	"github.com/NotFound1911/filestore/service/file_manager/repository"
)

type FileManagerService interface {
	InsertFileMeta(ctx context.Context, meta domain.FileMeta) (int64, error)
	InsertUserFile(ctx context.Context, file domain.UserFile) (int64, error)
	GetFileMetaByUserId(ctx context.Context, uid int64) ([]domain.FileMeta, error)
	GetFileMeta(ctx context.Context, sha1 string) (domain.FileMeta, error)
	GetUserIdsByFileSha1(ctx context.Context, sha1 string) ([]int64, error)
}
type fileManagerService struct {
	repo repository.FileManagerRepository
}

func (f *fileManagerService) InsertFileMeta(ctx context.Context, meta domain.FileMeta) (int64, error) {
	return f.repo.CreateFileMeta(ctx, meta)
}

func (f *fileManagerService) InsertUserFile(ctx context.Context, file domain.UserFile) (int64, error) {
	return f.repo.CreateUserFile(ctx, file)
}

func (f *fileManagerService) GetFileMetaByUserId(ctx context.Context, uid int64) ([]domain.FileMeta, error) {
	return f.repo.GetFileMetaByUserId(ctx, uid)
}

func (f *fileManagerService) GetUserIdsByFileSha1(ctx context.Context, sha1 string) ([]int64, error) {
	return f.repo.GetUserIdsByFileSha1(ctx, sha1)
}
func (f *fileManagerService) GetFileMeta(ctx context.Context, sha1 string) (domain.FileMeta, error) {
	return f.repo.GetFileMeta(ctx, sha1)
}
func NewFileManagerService(repo repository.FileManagerRepository) FileManagerService {
	return &fileManagerService{
		repo: repo,
	}
}
