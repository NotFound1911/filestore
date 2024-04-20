package repository

import (
	"context"
	"github.com/NotFound1911/filestore/service/upload/domain"
	"github.com/NotFound1911/filestore/service/upload/repository/dao"
	"time"
)

type UploadRepository interface {
	Create(ctx context.Context, u domain.Upload) (int64, error)
	UpdateStatus(ctx context.Context, u domain.Upload) error
	FindUploadById(ctx context.Context, id int64) (domain.Upload, error)
	GetUploadInfosByUId(ctx context.Context, uid int64) ([]domain.Upload, error)
}

type uploadRepo struct {
	dao dao.UploadDao
	// cache
}

func (repo *uploadRepo) Create(ctx context.Context, u domain.Upload) (int64, error) {
	return repo.dao.Insert(ctx, repo.toEntity(u))
}

func (repo *uploadRepo) UpdateStatus(ctx context.Context, u domain.Upload) error {
	return repo.dao.UpdateById(ctx, u.Id, map[string]interface{}{
		"status": u.Status,
		"u_time": time.Now(),
	})
}

func (repo *uploadRepo) FindUploadById(ctx context.Context, id int64) (domain.Upload, error) {
	res, err := repo.dao.FindById(ctx, id)
	if err != nil {
		return domain.Upload{}, err
	}
	return repo.toDomain(res), err
}

func (repo *uploadRepo) GetUploadInfosByUId(ctx context.Context, uid int64) ([]domain.Upload, error) {
	tmp, err := repo.dao.GetUploadInfosByUid(ctx, uid)
	if err != nil {
		return []domain.Upload{}, err
	}
	res := make([]domain.Upload, 0, len(tmp))
	for _, v := range tmp {
		res = append(res, repo.toDomain(v))
	}
	return res, nil
}

func NewUploadRepository(dao dao.UploadDao) UploadRepository {
	return &uploadRepo{
		dao: dao,
	}
}
func (repo *uploadRepo) toDomain(u dao.UploadInfo) domain.Upload {
	return domain.Upload{
		Id:       u.Id,
		UId:      u.UId,
		FileName: u.FileName,
		FileSha1: u.FileSha1,
		FileSize: u.FileSize,
		CreateAt: u.CTime,
		UpdateAt: u.UTime,
		Status:   u.Status,
	}
}

func (repo *uploadRepo) toEntity(u domain.Upload) dao.UploadInfo {
	return dao.UploadInfo{
		Id:       u.Id,
		UId:      u.UId,
		FileName: u.FileName,
		FileSha1: u.FileSha1,
		FileSize: u.FileSize,
		CTime:    u.CreateAt,
		UTime:    u.UpdateAt,
		Status:   u.Status,
	}
}
