package data

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
	"time"
)

type Manager struct {
	pool *pgxpool.Pool
	log  zerolog.Logger
}

func NewManager(pool *pgxpool.Pool, log zerolog.Logger) *Manager {
	return &Manager{
		pool: pool,
		log:  log,
	}
}

type Obj struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Data1     string    `json:"data1"`
	Data2     string    `json:"data2"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

func (m *Manager) GetAllData(ctx context.Context) ([]Obj, error) {
	var data []Obj
	conn, err := m.pool.Acquire(ctx) // get connection from pgx pool
	if err != nil {
		m.log.Fatal().Stack().Err(err).Msg("Unable to acquire a database connection")
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, "SELECT * FROM important_data")
	if err != nil {
		m.log.Fatal().Stack().Err(err).Msg("")

		return data, err
	}
	defer rows.Close()

	for rows.Next() {
		var obj Obj
		err = rows.Scan(&obj.ID, &obj.Data1, &obj.Data2, &obj.CreatedAt)
		if err != nil {
			m.log.Fatal().Stack().Err(err).Msg("")

			return data, err
		}
		data = append(data, obj)
	}
	return data, nil
}

func (m *Manager) AddDataObj(ctx context.Context, obj Obj) error {
	obj.ID = uuid.New()
	obj.CreatedAt = time.Now()

	conn, err := m.pool.Acquire(ctx)
	if err != nil {
		m.log.Fatal().Stack().Err(err).Msg("Unable to acquire a database connection")
	}
	defer conn.Release()

	_, err = conn.Exec(ctx,
		"INSERT INTO important_data (id, data1, data2, created_at) VALUES ($1, $2, $3, $4)",
		obj.ID, obj.Data1, obj.Data2, obj.CreatedAt)
	if err != nil {
		return err
	}
	m.log.Info().Msgf("add data => %v", obj)
	return nil
}

func (m *Manager) RemoveDataObj(ctx context.Context, id uuid.UUID) error {
	conn, err := m.pool.Acquire(ctx)
	if err != nil {
		m.log.Fatal().Stack().Err(err).Msg("Unable to acquire a database connection")
	}
	defer conn.Release()
	res, err := conn.Exec(ctx, "DELETE FROM important_data WHERE id=$1", id)
	if err != nil {
		return err
	}
	m.log.Debug().Msgf("result => %v", res)
	return err
}

func (m *Manager) UpdateDataObj(ctx context.Context, obj Obj) error {
	conn, err := m.pool.Acquire(ctx)
	if err != nil {
		m.log.Fatal().Stack().Err(err).Msg("Unable to acquire a database connection")
	}
	defer conn.Release()

	res, err := conn.Exec(ctx, "UPDATE important_data SET data1=$1, data2=$2 WHERE id=$3",
		obj.Data1, obj.Data2, obj.ID)
	if err != nil {
		return err
	}
	m.log.Debug().Msgf("result => %v", res)
	return nil
}
