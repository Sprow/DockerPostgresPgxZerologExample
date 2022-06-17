package redis

import (
	"DockerPostgreExample/internal/dto"
	"DockerPostgreExample/internal/logger"
	"context"
	"errors"
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

func (rdb *RDB) AddDataObj(ctx context.Context, obj dto.Obj) {
	res := rdb.client.HSet(
		ctx, redisID(obj.ID), "data1", obj.Data1, "data2", obj.Data2, "created_at", obj.CreatedAt)
	logger.Log.Info().Msgf("redis add data> %v", res)

	if obj.ID%2 == 0 { // pubsub example
		msg := fmt.Sprintf("congratulations to obj №%d", obj.ID)
		rdb.client.Publish(ctx, "congrats", msg)
	}
}

func (rdb *RDB) RemoveDataObj(ctx context.Context, id int) {
	rdb.client.HDel(ctx, redisID(id))
}

func (rdb *RDB) UpdateDataObj(ctx context.Context, obj dto.Obj) {
	rdb.client.HMGet(ctx, redisID(obj.ID), "data1", obj.Data1, "data2", obj.Data2)
}

func (rdb *RDB) GetObjById(ctx context.Context, id int) (dto.Obj, error) {
	var obj dto.Obj

	res, err := rdb.client.HGetAll(ctx, redisID(id)).Result()
	if err != nil {
		logger.Log.Error().Err(err).Msg("")
		return obj, err
	}
	if res["data1"] == "" || res["data2"] == "" {
		err = errors.New("no such obj in redis DB")
		return obj, err
	}
	t, err := time.Parse("2006-01-02T15:04:05Z07:00", res["created_at"])
	if err != nil {
		logger.Log.Error().Err(err).Msg("can't parse time")
	}
	obj = dto.Obj{
		ID:        id,
		Data1:     res["data1"],
		Data2:     res["data2"],
		CreatedAt: t,
	}

	return obj, nil
}

func (rdb *RDB) Sub(ctx context.Context, chName string) <-chan *redis.Message {
	sub := rdb.client.Subscribe(ctx, chName)
	return sub.Channel()
}

func redisID(id int) string {
	return "obj:" + strconv.Itoa(id)
}
