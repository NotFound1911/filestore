package domain

import "time"

// FileMeta 文件元数据
type FileMeta struct {
	Sha1     string     `json:"sha1"`
	Name     string     `json:"name"`
	Size     int64      `json:"size"`
	Address  string     `json:"address"`
	Status   string     `json:"status"`
	CreateAt *time.Time `json:"create_at"`
	UpdateAt *time.Time `json:"update_at"`
}
