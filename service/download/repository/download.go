package repository

import (
	"context"
	"github.com/NotFound1911/filestore/service/download/domain"
	"github.com/NotFound1911/filestore/service/download/repository/dao"
	"time"
)

type DownloadRepository interface {
	Create(ctx context.Context, d domain.Download) (int64, error)
	UpdateStatus(ctx context.Context, u domain.Download) error
	FindUploadById(ctx context.Context, id int64) (domain.Download, error)
	GetDownloadInfosByUId(ctx context.Context, uid int64) ([]domain.Download, error)
}

type downloadRepo struct {
	dao dao.DownloadDao
}

func NewDownloadRepository(dao dao.DownloadDao) DownloadRepository {
	return &downloadRepo{
		dao: dao,
	}
}

func (d *downloadRepo) Create(ctx context.Context, down domain.Download) (int64, error) {
	return d.dao.Insert(ctx, d.toEntity(down))
}

func (d *downloadRepo) UpdateStatus(ctx context.Context, down domain.Download) error {
	return d.dao.UpdateById(ctx, down.Id, map[string]interface{}{
		"status": down.Status,
		"u_time": time.Now(),
	})
}

func (d *downloadRepo) FindUploadById(ctx context.Context, id int64) (domain.Download, error) {
	res, err := d.dao.FindById(ctx, id)
	if err != nil {
		return domain.Download{}, err
	}
	return d.toDomain(res), err
}

func (d *downloadRepo) GetDownloadInfosByUId(ctx context.Context, uid int64) ([]domain.Download, error) {
	tmp, err := d.dao.GetDownloadInfosByUid(ctx, uid)
	if err != nil {
		return []domain.Download{}, err
	}
	res := make([]domain.Download, 0, len(tmp))
	for _, v := range tmp {
		res = append(res, d.toDomain(v))
	}
	return res, nil
}
func (d *downloadRepo) toEntity(down domain.Download) dao.DownloadInfo {
	return dao.DownloadInfo{
		Id:       down.Id,
		UId:      down.UId,
		FileName: down.FileName,
		FileSize: down.FileSize,
		CTime:    down.CreateAt,
		UTime:    down.UpdateAt,
		Status:   down.Status,
	}
}
func (d *downloadRepo) toDomain(down dao.DownloadInfo) domain.Download {
	return domain.Download{
		Id:       down.Id,
		UId:      down.UId,
		FileName: down.FileName,
		FileSize: down.FileSize,
		CreateAt: down.CTime,
		UpdateAt: down.UTime,
		Status:   down.Status,
	}
}
