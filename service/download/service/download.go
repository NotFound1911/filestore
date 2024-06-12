package service

import (
	"context"
	"github.com/NotFound1911/filestore/service/download/domain"
	"github.com/NotFound1911/filestore/service/download/repository"
)

type DownloadService interface {
	Download(ctx context.Context, d domain.Download) (int64, error)
	UpdateStatus(ctx context.Context, id int64, status domain.Status) error
	FindDownloadById(ctx context.Context, id int64) (domain.Download, error)
	GetDownloadInfosByUId(ctx context.Context, uid int64) ([]domain.Download, error)
}

func NewDownloadService(repo repository.DownloadRepository) DownloadService {
	return &downloadService{
		repo: repo,
	}
}

type downloadService struct {
	repo repository.DownloadRepository
}

func (d *downloadService) Download(ctx context.Context, down domain.Download) (int64, error) {
	return d.repo.Create(ctx, down)
}

func (d *downloadService) UpdateStatus(ctx context.Context, id int64, status domain.Status) error {
	down := domain.Download{
		Id:     id,
		Status: status,
	}
	return d.repo.UpdateStatus(ctx, down)
}

func (d *downloadService) FindDownloadById(ctx context.Context, id int64) (domain.Download, error) {
	return d.repo.FindUploadById(ctx, id)
}

func (d *downloadService) GetDownloadInfosByUId(ctx context.Context, uid int64) ([]domain.Download, error) {
	return d.repo.GetDownloadInfosByUId(ctx, uid)
}
