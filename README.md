
## CandleStream Service


#### Running the Application in Development Environment

### Setup

1. **Clone the Repository**

   ```bash
   git clone https://your-repository-url.git
   cd your-repository-directory
   
2. Start docker compose

   ```bash
    docker-compose up
    ```
   
3. Database Reset - (Database Initialization)
   ```bash
   make reset
   ```
3. Run application

   ```bash
   make start
   ```

#### Running the Application in Containerized Environment


### Setup

1. **Clone the Repository**

   ```bash
   git clone https://your-repository-url.git
   cd your-repository-directory

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



