
-- CREATE SCHEMA IF NOT EXISTS base

CREATE TABLE IF NOT EXISTS important_data(
    id SERIAL NOT NULL PRIMARY KEY,
    data1 TEXT,
    data2 VARCHAR(50),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- create indexes example
create index if not exists important_data_data1_index on important_data (data1);
create index if not exists important_data_data2_index on important_data (data2);
create index if not exists important_data_created_at_index on important_data (created_at);