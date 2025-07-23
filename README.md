# Financial Data Backend

A backend service for stock time-series data from the Alpha Vantage API. Designed using Clean Architecture principles with support for both SQL and NoSQL storage.
### API Endpoints
| Method | Endpoint        | Description                         |
| ------ | --------------- | ----------------------------------- |
| GET    | `/symbols`      | Get selection of symbols, given they match keyed url query argument "keywords"     |
| POST   | `/data/:symbol` | Fetch and store new stock data from up to last 2-3 weeks     |
| DELETE | `/data/:symbol` | Delete a symbol and its stored data |
| GET    | `/data`         | Retrieve all stored stock data      |
### Tech Stack
* Language: Go (Gin, testing and mocking packages)
* Storage Options: MongoDB Atlas (NoSQL), PostgreSQL
* Other Tools: GitHub, Postman
### Key Features
* REST API for fetching and managing stock data (daily interval, via Alpha Vantage)
* Clean Architecture: separated handler, usecase, repository layers
* Timeout middleware (for MongoDB Atlas cloud latency)
* Centralised error-handling middleware (all branches)
* Unit tests with mocks for core logic (ongoing expansion planned)
### Branch Overview
* main: Mirrors the mongodb branch and includes working REST endpoints with MongoDB Atlas
* mongodb: Uses MongoDB Atlas for storage; includes timeout middleware
* postgresql: Uses PostgreSQL; excludes timeout middleware
* Each branch contains 4 functional endpoints and is independently usable

