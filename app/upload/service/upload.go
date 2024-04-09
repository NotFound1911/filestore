package service

import (
	"context"
	"github.com/NotFound1911/filestore/app/upload/domain"
	"github.com/NotFound1911/filestore/app/upload/repository"
)

type UploadService interface {
	Upload(ctx context.Context, u domain.Upload) (int64, error)
	UpdateStatus(ctx context.Context, id int64, status string) error
	FindUploadById(ctx context.Context, id int64) (domain.Upload, error)
	GetUploadInfosByUId(ctx context.Context, uid int64) ([]domain.Upload, error)
	// todo 妙传 分片
}
type uploadService struct {
	repo repository.UploadRepository
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
