package domain

import "time"

// Download 下载信息
type Download struct {
	Id       int64      `json:"id"`        // 下载id
	UId      int64      `json:"u_id"`      // 用户id
	FileName string     `json:"file_name"` // 文件名称
	FileSize int64      `json:"file_size"` // 文件size
	CreateAt *time.Time `json:"create_at"` // 下载开始时间
	UpdateAt *time.Time `json:"update_at"` // 下载更新时间
	Status   Status     `json:"status"`    // 下载状态
}

type Status uint32

const (
	Start Status = iota
	Finished
)
