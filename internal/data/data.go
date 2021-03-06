package data

import (
	"DockerPostgreExample/internal/logger"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Manager struct {
	pool *pgxpool.Pool
}

func NewManager(pool *pgxpool.Pool) *Manager {
	return &Manager{
		pool: pool,
	}
}

func (m *Manager) GetAllData(ctx context.Context) ([]Obj, error) {
	var data []Obj
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
		var obj Obj
		err = rows.Scan(&obj.ID, &obj.Data1, &obj.Data2, &obj.CreatedAt)
		if err != nil {
			logger.Log.Fatal().Stack().Err(err).Msg("")

			return data, err
		}
		data = append(data, obj)
	}
	return data, nil
}

func (m *Manager) AddDataObj(ctx context.Context, obj Obj) error {
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
	logger.Log.Info().Msgf("add data => %v", obj)
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
	return err
}

func (m *Manager) UpdateDataObj(ctx context.Context, obj Obj) error {
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
	return nil
}
