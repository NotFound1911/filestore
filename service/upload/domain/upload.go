package domain

import "time"

const (
	UploadTypeSingle = "Single"
	UploadTypeMulti  = "Multi"

	UploadStatusStart    = "Start"
	UploadStatusFinished = "Finished"
)

// Upload 文件上传信息
type Upload struct {
	Id       int64      `json:"id"`        // 上传id
	UId      int64      `json:"u_id"`      // 用户id
	FileName string     `json:"file_name"` // 文件名称
	FileSha1 string     `json:"file_sha1"` // 文件sha1
	FileSize int64      `json:"file_size"` // 文件size
	CreateAt *time.Time `json:"create_at"` // 上传开始时间
	UpdateAt *time.Time `json:"update_at"` // 上传更新时间
	Status   string     `json:"status"`    // 上传状态
	Type     string     `json:"type"`      // 任务类别：single 不分片， multi 分片
}

// Chunk 分块对象
type Chunk struct {
	Id       int64      `json:"id"`        // 分块序号
	UId      int64      `json:"u_id"`      // 用户id
	UploadId int64      `json:"upload_id"` // 上传任务关联id
	FileName string     `json:"file_name"` // 文件名
	Sha1     string     `json:"sha1"`      // 文件sha1
	Size     int64      `json:"size"`      // 文件size
	CreateAt *time.Time `json:"create_at"` // 上传开始时间
	UpdateAt *time.Time `json:"update_at"` // 上传更新时间
	Status   string     `json:"status"`    // 上传状态
	Count    int64      `json:"count"`     // 分块总数
}
