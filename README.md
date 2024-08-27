
# Trading Candles Stream Service

## Overview
The Candles Stream Service is a gRPC-based microservice that connects to the Binance API to fetch real-time trade data for cryptocurrency pairs such as BTCUSDT, ETHUSDT, and PEPEUSDT. The service aggregates the trade data into OHLC (Open, High, Low, Close) candlesticks with a 1-minute timeframe and stores them in a PostgreSQL database. It also provides a gRPC API for broadcasting these candlestick bars.

## Features
- **Real-time Data Streaming:** Connects to Binance API to stream real-time trade data.
- **Candlestick Aggregation:** Aggregates trades into 1-minute OHLC candlesticks.
- **gRPC API:** Provides a gRPC service for broadcasting candlestick data.
- **PostgreSQL Storage:** Stores candlestick data in a PostgreSQL database.
- **Dockerized Service:** Runs as a Docker container for easy deployment.

## Installation
To get started with the Candles Stream Service, follow these steps:


## 🚀 **Running the Application in Development Environment**

### Setup

1. **Clone the Repository**

   ```bash
   git clone https://github.com/sreejaunni/trading-candles-stream-service.git
   cd candles-stream-service
   
2. Start docker compose

   ```bash
    docker-compose up
    ```
   
3. Database Reset - (Database Initialization)
   ```bash
   make reset-db
   ```
3. Run application

   ```bash
   make start
   ```

## 🚀 **Running the Application in Containerized Environment**


### Setup

1. **Clone the Repository**

   ```bash
   git clone https://github.com/sreejaunni/trading-candles-stream-service.git
   cd candles-stream-service

2. make build  (build image)
   ```bash
   make build
   ```
3. Start docker compose

   ```bash
    docker-compose up
    ```

3. Database Reset -(Database Initialization)

   ```bash
   docker run --name marketplace-migration-container -e DB_HOST=host.docker.internal marketplace-app ./app reset-database
   ```
4. Run application

   ```bash
   docker run --name marketplace-application-container -e DB_HOST=host.docker.internal marketplace-app ./app server
   ```



