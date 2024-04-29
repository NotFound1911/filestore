package cache

import (
	"context"
	"github.com/NotFound1911/filestore/service/upload/domain"
	cachemocks "github.com/NotFound1911/filestore/service/upload/repository/cache/mocks"
	"github.com/golang/mock/gomock"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRedisChunkCache_GetChunks(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(controller *gomock.Controller) redis.Cmdable
		ctx      context.Context
		uploadId int64
		chunks   []domain.Chunk

		wantErr error
	}{
		{
			name:     "正常获取",
			uploadId: 3,
			ctx:      context.Background(),
			mock: func(controller *gomock.Controller) redis.Cmdable {
				cmd := cachemocks.NewMockCmdable(controller)
				res := redis.NewCmd(context.Background())
				res.SetVal(`[{
		"id": 1,
		"u_id": 2,
		"upload_id": 3,
		"file_name": "text.txt",
		"sha1": "6bee56615ac5fe5b5c37d7289ab2ffaa13442fb1",
		"size": 1024,
		"create_at": "2024-04-26T15:04:05Z",
		"update_at": "2024-04-26T15:04:05Z",
		"status": "start",
		"count": 2
	},
	{
		"id": 2,
		"u_id": 2,
		"upload_id": 3,
		"file_name": "text.txt",
		"sha1": "6bee56615ac5fe5b5c37d7289ab2ffaa13442fb2",
		"size": 1024,
		"create_at": "2024-04-26T15:04:05Z",
		"update_at": "2024-04-26T15:04:05Z",
		"status": "start",
		"count": 2
	}
]`)
				cmd.EXPECT().Eval(gomock.Any(), luaGetChunks,
					[]string{(&RedisChunkCache{}).pattern(3)}).Return(res)
				return cmd
			},
			chunks: func() []domain.Chunk {
				timeStr := "2024-04-26 15:04:05"
				// 定义时间格式
				layout := "2006-01-02 15:04:05"
				// 使用 time.Parse 解析时间字符串
				t, err := time.Parse(layout, timeStr)
				if err != nil {
					panic(err)
				}
				res := []domain.Chunk{
					{
						Id:       1,
						UId:      2,
						UploadId: 3,
						FileName: "text.txt",
						Sha1:     "6bee56615ac5fe5b5c37d7289ab2ffaa13442fb1",
						Size:     1024,
						CreateAt: &t,
						UpdateAt: &t,
						Status:   "start",
						Count:    2,
					},
					{
						Id:       2,
						UId:      2,
						UploadId: 3,
						FileName: "text.txt",
						Sha1:     "6bee56615ac5fe5b5c37d7289ab2ffaa13442fb2",
						Size:     1024,
						CreateAt: &t,
						UpdateAt: &t,
						Status:   "start",
						Count:    2,
					},
				}
				return res
			}(),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			cache := NewChunkCache(tc.mock(ctrl))
			cs, err := cache.GetChunks(tc.ctx, tc.uploadId)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.chunks, cs)
		})
	}
}
