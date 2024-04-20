package repository

import (
	"context"
	"github.com/NotFound1911/filestore/service/file_manager/domain"
	"github.com/NotFound1911/filestore/service/file_manager/repository/dao"
)

type FileManagerRepository interface {
	CreateFileMeta(ctx context.Context, f domain.FileMeta) (int64, error)
	CreateUserFile(ctx context.Context, u domain.UserFile) (int64, error)
	UpdateFileMeta(ctx context.Context, f domain.FileMeta) error
	GetFileMetaByUserId(ctx context.Context, uid int64) ([]domain.FileMeta, error)
	GetUserIdsByFileSha1(ctx context.Context, sha1 string) ([]int64, error)
}
type fileManagerRepo struct {
	dao dao.FileManagerDao
	// cache
}

func (repo *fileManagerRepo) CreateFileMeta(ctx context.Context, f domain.FileMeta) (int64, error) {
	return repo.dao.InsertFileMeta(ctx, repo.toFileMetaEntity(f))
}

func (repo *fileManagerRepo) CreateUserFile(ctx context.Context, u domain.UserFile) (int64, error) {
	return repo.dao.InsertUserFile(ctx, repo.toUserFileEntity(u))
}

func (repo *fileManagerRepo) UpdateFileMeta(ctx context.Context, f domain.FileMeta) error {
	return repo.dao.UpdateFileMetaById(ctx, f.Id, map[string]interface{}{
		"type":    f.Type,
		"address": f.Address,
	})
}

func (repo *fileManagerRepo) GetFileMetaByUserId(ctx context.Context, uid int64) ([]domain.FileMeta, error) {
	tmp, err := repo.dao.GetFileMetasByUserId(ctx, uid)
	if err != nil {
		return []domain.FileMeta{}, err
	}
	res := make([]domain.FileMeta, 0, len(tmp))
	for _, v := range tmp {
		res = append(res, repo.toFileMetaDomain(v))
	}
	return res, nil
}

func (repo *fileManagerRepo) GetUserIdsByFileSha1(ctx context.Context, sha1 string) ([]int64, error) {
	return repo.dao.GetUserIdsByFileSha1(ctx, sha1)
}

func NewFileManagerRepository(dao dao.FileManagerDao) FileManagerRepository {
	return &fileManagerRepo{
		dao: dao,
	}
}
func (repo *fileManagerRepo) toFileMetaEntity(u domain.FileMeta) dao.FileMetaInfo {
	return dao.FileMetaInfo{
		Id:      u.Id,
		Sha1:    u.Sha1,
		Size:    u.Size,
		Address: u.Address,
		Type:    u.Type,
	}
}

func (repo *fileManagerRepo) toUserFileEntity(u domain.UserFile) dao.UserFileInfo {
	return dao.UserFileInfo{
		Id:       u.Id,
		UserId:   u.UserId,
		FileSize: u.FileSize,
		FileName: u.FileName,
		FileSha1: u.FileSha1,
		UpdateAt: u.UpdateAt,
	}
}
func (repo *fileManagerRepo) toFileMetaDomain(u dao.FileMetaInfo) domain.FileMeta {
	return domain.FileMeta{
		Id:      u.Id,
		Sha1:    u.Sha1,
		Size:    u.Size,
		Address: u.Address,
		Type:    u.Type,
	}
}
func (repo *fileManagerRepo) toUserFileDomain(u dao.UserFileInfo) domain.UserFile {
	return domain.UserFile{
		Id:       u.Id,
		UserId:   u.UserId,
		FileSize: u.FileSize,
		FileName: u.FileName,
		FileSha1: u.FileSha1,
		UpdateAt: u.UpdateAt,
	}
}
