package service

import (
	"context"
	"github.com/NotFound1911/filestore/service/upload/domain"
	"github.com/NotFound1911/filestore/service/upload/repository"
)

type UploadService interface {
	Upload(ctx context.Context, u domain.Upload) (int64, error)
	UpdateStatus(ctx context.Context, id int64, status string) error
	FindUploadById(ctx context.Context, id int64) (domain.Upload, error)
	GetUploadInfosByUId(ctx context.Context, uid int64) ([]domain.Upload, error)

	SetFileChunk(ctx context.Context, c domain.Chunk) error
	GetFileChunk(ctx context.Context, uploadId int64, id int64) (domain.Chunk, error)
	DelFileChunk(ctx context.Context, uploadId int64, id int64) error
	GetChunks(ctx context.Context, uploadId int64) ([]domain.Chunk, error)
}
type uploadService struct {
	repo repository.UploadRepository
}

func (s *uploadService) GetChunks(ctx context.Context, uploadId int64) ([]domain.Chunk, error) {
	return s.repo.GetChunks(ctx, uploadId)
}

func (s *uploadService) SetFileChunk(ctx context.Context, c domain.Chunk) error {
	return s.repo.SetFileChunk(ctx, c)
}

func (s *uploadService) GetFileChunk(ctx context.Context, uploadId int64, id int64) (domain.Chunk, error) {
	return s.repo.GetFileChunk(ctx, uploadId, id)
}

func (s *uploadService) DelFileChunk(ctx context.Context, uploadId int64, id int64) error {
	return s.repo.DelFileChunk(ctx, uploadId, id)
}

func (s *uploadService) Upload(ctx context.Context, u domain.Upload) (int64, error) {
	return s.repo.Create(ctx, u)
}

func (s *uploadService) UpdateStatus(ctx context.Context, id int64, status string) error {
	u := domain.Upload{
		Id:     id,
		Status: status,
	}
	return s.repo.UpdateStatus(ctx, u)
}

func (s *uploadService) FindUploadById(ctx context.Context, id int64) (domain.Upload, error) {
	return s.repo.FindUploadById(ctx, id)
}

func (s *uploadService) GetUploadInfosByUId(ctx context.Context, uid int64) ([]domain.Upload, error) {
	return s.repo.GetUploadInfosByUId(ctx, uid)
}

func NewUploadService(repo repository.UploadRepository) UploadService {
	return &uploadService{
		repo: repo,
	}
}
