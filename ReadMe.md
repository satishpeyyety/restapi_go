# Developing PostGre-based REST API
## Install required libraries
```
go get -u github.com/gin-gonic/gin        # For HTTP server
go get -u github.com/jackc/pgx/v4         # PostgreSQL driver
go get -u github.com/jackc/pgx/v4/pgxpool  # Connection pool
go get -u github.com/joho/godotenv        # For environment variables
```
go.mod file will updated with all the dependencies automatically.

## Setup table in Yugabyte
```
CREATE TABLE employees (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
name TEXT NOT NULL,
position TEXT NOT NULL,
salary FLOAT NOT NULL
);
```
