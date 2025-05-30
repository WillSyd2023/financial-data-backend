--CREATE DATABASE stockfeed_db;

--DROP SCHEMA public CASCADE;
--CREATE SCHEMA public;

CREATE TABLE symbols (
    symbol_id BIGSERIAL PRIMARY KEY,
    symbol VARCHAR UNIQUE NOT NULL,
    last_refreshed DATE NOT NULL
);

CREATE TABLE ohlcv_per_day (
    ohlcv_id BIGSERIAL PRIMARY KEY,
    record_day DATE NOT NULL,
    open_price NUMERIC(12, 4) NOT NULL,
    high_price NUMERIC(12, 4) NOT NULL,
    low_price NUMERIC(12, 4) NOT NULL,
    close_price NUMERIC(12, 4) NOT NULL,
    volume INTEGER NOT NULL,
    symbol_id BIGINT NOT NULL REFERENCES symbols(symbol_id)
       ON DELETE CASCADE
);