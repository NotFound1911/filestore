package repository

import (
	"context"
	"github.com/NotFound1911/filestore/service/upload/domain"
	"github.com/NotFound1911/filestore/service/upload/repository/cache"
	"github.com/NotFound1911/filestore/service/upload/repository/dao"
	"time"
)

type UploadRepository interface {
	Create(ctx context.Context, u domain.Upload) (int64, error)
	UpdateStatus(ctx context.Context, u domain.Upload) error
	FindUploadById(ctx context.Context, id int64) (domain.Upload, error)
	GetUploadInfosByUId(ctx context.Context, uid int64) ([]domain.Upload, error)

	SetFileChunk(ctx context.Context, c domain.Chunk) error
	GetFileChunk(ctx context.Context, uploadId int64, id int64) (domain.Chunk, error)
	DelFileChunk(ctx context.Context, uploadId int64, id int64) error
	GetChunks(ctx context.Context, uploadId int64) ([]domain.Chunk, error)
}

type uploadRepo struct {
	dao   dao.UploadDao
	cache cache.ChunkCache
	// logFunc func(format string, a ...any)
}

func (repo *uploadRepo) SetFileChunk(ctx context.Context, c domain.Chunk) error {
	return repo.cache.Set(ctx, c)
}

func (repo *uploadRepo) GetChunks(ctx context.Context, uploadId int64) ([]domain.Chunk, error) {
	return repo.cache.GetChunks(ctx, uploadId)
}

func (repo *uploadRepo) GetFileChunk(ctx context.Context, uploadId int64, id int64) (domain.Chunk, error) {
	return repo.cache.Get(ctx, uploadId, id)
}

func (repo *uploadRepo) DelFileChunk(ctx context.Context, uploadId int64, id int64) error {
	return repo.cache.Del(ctx, uploadId, id)
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

func NewUploadRepository(dao dao.UploadDao, chunkCache cache.ChunkCache) UploadRepository {
	return &uploadRepo{
		dao:   dao,
		cache: chunkCache,
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
		Type:     u.Type,
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
		Type:     u.Type,
	}
}
