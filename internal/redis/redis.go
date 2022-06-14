package redis

import (
	"DockerPostgreExample/internal/dto"
	"DockerPostgreExample/internal/logger"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

type RDB struct {
	client *redis.Client
}

func NewRDB(rdb *redis.Client) *RDB {
	return &RDB{
		client: rdb,
	}
}
	
func (rdb *RDB) AddDataObj(ctx context.Context, id int, data1 string, data2 string, created_at time.Time) {
	res := rdb.client.HSet(ctx, redisID(id), "data1", data1, "data2", data2, "created_at", created_at)
	logger.Log.Info().Msgf("redis add data> %v", res)

	if id%2 == 0 { // pubsub example
		msg := fmt.Sprintf("congratulations to obj №%d", id)
		rdb.client.Publish(ctx, "congrats", msg)
	}
}

func (rdb *RDB) RemoveDataObj(ctx context.Context, id int) {
	rdb.client.HDel(ctx, redisID(id))
}

func (rdb *RDB) UpdateDataObj(ctx context.Context, id int, data1 string, data2 string) {
	rdb.client.HMGet(ctx, redisID(id), "data1", data1, "data2", data2)
}

func (rdb *RDB) GetObjById(ctx context.Context, id int) dto.Obj {
	var obj dto.Obj

	res, err := rdb.client.HGetAll(ctx, redisID(id)).Result()
	if err != nil {
		logger.Log.Error().Err(err).Msg("")
	} else {
		obj = dto.Obj{
			ID:    id,
			Data1: res["data1"],
			Data2: res["data2"],
		}
	}
	return obj
}

func (rdb *RDB) Sub(ctx context.Context, chName string) <-chan *redis.Message {
	sub := rdb.client.Subscribe(ctx, chName)
	return sub.Channel()
}

func redisID(id int) string {
	return "obj:" + strconv.Itoa(id)
}
