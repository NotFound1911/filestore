package repository

import (
	"context"
	"fmt"
	"github.com/NotFound1911/filestore/domain"
	"github.com/NotFound1911/filestore/repository/dao"
)

type FileRepository interface {
	CreateFileMeta(ctx context.Context, f domain.FileMeta) error
	UpdateFileMeta(ctx context.Context, f domain.FileMeta) error
	GetFileMetaBySha1(ctx context.Context, sha1 string) (domain.FileMeta, error)
	DeleteFileMetaBySha1(ctx context.Context, sha1 string) error
}
type CachedFileRepository struct {
	dao dao.FileDAO
	//cache   cache.FileCache
	logFunc func(format string, a ...any)
}

func (repo *CachedFileRepository) CreateFileMeta(ctx context.Context, f domain.FileMeta) error {
	return repo.dao.Insert(ctx, repo.toEntity(f))
}

func (repo *CachedFileRepository) UpdateFileMeta(ctx context.Context, f domain.FileMeta) error {
	return repo.dao.UpdateBySha1(ctx, repo.toEntity(f))
}

func (repo *CachedFileRepository) GetFileMetaBySha1(ctx context.Context, sha1 string) (domain.FileMeta, error) {
	return repo.GetFileMetaBySha1(ctx, sha1)
}

func (repo *CachedFileRepository) DeleteFileMetaBySha1(ctx context.Context, sha1 string) error {
	return repo.DeleteFileMetaBySha1(ctx, sha1)
}
func (repo *CachedFileRepository) toDomain(f dao.FileMeta) domain.FileMeta {
	return domain.FileMeta{
		Sha1:     f.Sha1,
		Name:     f.Name,
		Size:     f.Size,
		Address:  f.Address,
		Status:   f.Status,
		CreateAt: f.CTime,
		UpdateAt: f.UTime,
	}
}

func (repo *CachedFileRepository) toEntity(f domain.FileMeta) dao.FileMeta {
	return dao.FileMeta{
		Sha1:    f.Sha1,
		Name:    f.Name,
		Size:    f.Size,
		Address: f.Address,
		Status:  f.Status,
		CTime:   f.CreateAt,
		UTime:   f.UpdateAt,
	}
}

// NewCachedFileRepository todo 缓存模式
func NewCachedFileRepository(
	dao dao.FileDAO,
	// c cache.UserCache
) FileRepository {
	return &CachedFileRepository{
		dao: dao,
		//cache: c,
		logFunc: func(format string, a ...any) {
			fmt.Printf(format, a)
		},
	}
}
