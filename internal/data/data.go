package data

import (
	"DockerPostgreExample/internal/dto"
	"DockerPostgreExample/internal/logger"
	"DockerPostgreExample/internal/redis"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Manager struct {
	pool    *pgxpool.Pool
	redisDB *redis.RDB
}

func NewManager(pool *pgxpool.Pool, rdb *redis.RDB) *Manager {
	return &Manager{
		pool:    pool,
		redisDB: rdb,
	}
}

func (m *Manager) GetAllData(ctx context.Context) ([]dto.Obj, error) {
	var data []dto.Obj
	conn, err := m.pool.Acquire(ctx) // get connection from pgx pool
	if err != nil {
		logger.Log.Fatal().Stack().Err(err).Msg("Unable to acquire a database connection")
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, "SELECT * FROM important_data")
	if err != nil {
		logger.Log.Fatal().Stack().Err(err).Msg("")

		return data, err
	}
	defer rows.Close()

	for rows.Next() {
		var obj dto.Obj
		err = rows.Scan(&obj.ID, &obj.Data1, &obj.Data2, &obj.CreatedAt)
		if err != nil {
			logger.Log.Fatal().Stack().Err(err).Msg("")

			return data, err
		}
		data = append(data, obj)
	}
	return data, nil
}

func (m *Manager) AddDataObj(ctx context.Context, obj dto.Obj) error {
	//postgres
	conn, err := m.pool.Acquire(ctx)
	if err != nil {
		logger.Log.Fatal().Stack().Err(err).Msg("Unable to acquire a database connection")
	}
	defer conn.Release()

	err = conn.QueryRow(ctx,
		"INSERT INTO important_data (data1, data2) VALUES ($1, $2) RETURNING id, created_at",
		obj.Data1, obj.Data2).Scan(&obj.ID, &obj.CreatedAt)
	if err != nil {
		return err
	}
	//logger.Log.Info().Msgf("postgres add data => %v", obj)

	//redis
	m.redisDB.AddDataObj(ctx, obj.ID, obj.Data1, obj.Data2, obj.CreatedAt)

	return nil
}

func (m *Manager) RemoveDataObj(ctx context.Context, id int) error {
	conn, err := m.pool.Acquire(ctx)
	if err != nil {
		logger.Log.Fatal().Stack().Err(err).Msg("Unable to acquire a database connection")
	}
	defer conn.Release()
	res, err := conn.Exec(ctx, "DELETE FROM important_data WHERE id=$1", id)
	if err != nil {
		return err
	}
	logger.Log.Info().Msgf("result => %v", res)

	m.redisDB.RemoveDataObj(ctx, id) //redis
	return err
}

func (m *Manager) UpdateDataObj(ctx context.Context, obj dto.Obj) error {
	conn, err := m.pool.Acquire(ctx)
	if err != nil {
		logger.Log.Fatal().Stack().Err(err).Msg("Unable to acquire a database connection")
	}
	defer conn.Release()

	res, err := conn.Exec(ctx, "UPDATE important_data SET data1=$1, data2=$2 WHERE id=$3",
		obj.Data1, obj.Data2, obj.ID)
	if err != nil {
		return err
	}
	logger.Log.Info().Msgf("result => %v", res)

	m.redisDB.UpdateDataObj(ctx, obj.ID, obj.Data1, obj.Data2)
	return nil
}

func (m *Manager) GetObjById(ctx context.Context, id int) (dto.Obj, error) {
	var obj dto.Obj

	// return obj from redis if exists
	redisObj := m.redisDB.GetObjById(ctx, id)
	if redisObj != obj {
		logger.Log.Info().Msgf("return obj id=%d from redis", id)
		return redisObj, nil
	}

	// else return obj from postgres if exists
	conn, err := m.pool.Acquire(ctx)
	if err != nil {
		logger.Log.Fatal().Stack().Err(err).Msg("Unable to acquire a database connection")
		return obj, err
	}
	defer conn.Release()
	row := conn.QueryRow(ctx, "SELECT * FROM important_data WHERE id=$1", id)
	if err != nil {
		return obj, err
	}
	err = row.Scan(&obj.ID, &obj.Data1, &obj.Data2, &obj.CreatedAt)
	if err != nil {
		logger.Log.Error().Stack().Err(err).Msg("")
	}
	logger.Log.Info().Msgf("return obj id=%d from postgres", id)
	return obj, err
}
