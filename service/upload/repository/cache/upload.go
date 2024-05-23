package cache

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/NotFound1911/filestore/service/upload/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

//go:embed lua/get_chunks.lua
var luaGetChunks string

var ErrKeyNotExist = redis.Nil

//go:generate mockgen -source=./upload.go -package=cachemocks -destination=./mocks/upload.mock.go ChunkCache
//go:generate $GOPATH/bin/mockgen -destination=mocks/mock_redis_cmdable.gen.go -package=cachemocks github.com/redis/go-redis/v9 Cmdable
type ChunkCache interface {
	Get(ctx context.Context, uploadId int64, id int64) (domain.Chunk, error)
	GetChunks(ctx context.Context, uploadId int64) ([]domain.Chunk, error)
	Set(ctx context.Context, c domain.Chunk) error
	Del(ctx context.Context, uploadId int64, id int64) error
}

type RedisChunkCache struct {
	cmd        redis.Cmdable
	expiration time.Duration
}

func (r *RedisChunkCache) key(uploadId int64, id int64) string {
	return fmt.Sprintf("MP:Chunk:%d:%d", uploadId, id)
}
func (r *RedisChunkCache) pattern(uploadId int64) string {
	return fmt.Sprintf("MP:Chunk:%d:*", uploadId)
}

func (r *RedisChunkCache) GetChunks(ctx context.Context, uploadId int64) ([]domain.Chunk, error) {
	res, err := r.cmd.Eval(ctx, luaGetChunks, []string{r.pattern(uploadId)}).Result()
	if err != nil {
		return nil, err
	}
	resStrs, ok := res.([]interface{})
	if !ok {
		return nil, fmt.Errorf("%v is not type []interface{}", res)
	}
	cs := make([]domain.Chunk, 0, len(resStrs))
	for i := 0; i < len(resStrs); i++ {
		c := domain.Chunk{}
		str, ok := resStrs[i].(string)
		if !ok {
			return nil, fmt.Errorf("%v is not type string", resStrs[i])
		}
		if err := json.Unmarshal([]byte(str), &c); err != nil {
			return nil, err
		}
		cs = append(cs, c)
	}
	return cs, nil
}

func (r *RedisChunkCache) Get(ctx context.Context, uploadId int64, id int64) (domain.Chunk, error) {
	key := r.key(uploadId, id)
	data, err := r.cmd.Get(ctx, key).Result()
	if err != nil {
		return domain.Chunk{}, err
	}
	var c domain.Chunk
	err = json.Unmarshal([]byte(data), &c)
	return c, err
}

func (r *RedisChunkCache) Set(ctx context.Context, c domain.Chunk) error {
	key := r.key(c.UploadId, c.Id)
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return r.cmd.Set(ctx, key, data, r.expiration).Err()
}

func (r *RedisChunkCache) Del(ctx context.Context, uploadId int64, id int64) error {
	return r.cmd.Del(ctx, r.key(uploadId, id)).Err()
}

func NewChunkCache(cmd redis.Cmdable) ChunkCache {
	return &RedisChunkCache{
		cmd:        cmd,
		expiration: time.Minute * 15,
	}
}
