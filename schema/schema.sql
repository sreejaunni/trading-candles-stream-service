
DROP TABLE IF EXISTS candlesticks;

CREATE TABLE IF NOT EXISTS candlesticks (
     id SERIAL PRIMARY KEY,
     symbol VARCHAR(10) NOT NULL,
     open DECIMAL(28, 10) NOT NULL,
     high DECIMAL(28, 10) NOT NULL,
     low DECIMAL(28, 10) NOT NULL,
     close DECIMAL(28, 10) NOT NULL,
     volume DECIMAL(28, 10) NOT NULL,
     start_timestamp TIMESTAMP NOT NULL,
     end_timestamp  TIMESTAMP NOT NULL
    );

