# Create and Run PostgreSQL Container and Initialize Database
```bash
docker run --name postgres-tictactoe \
  -e POSTGRES_PASSWORD=mysecretpassword \
  -e POSTGRES_DB=tictactoe \
  -v ./init.sql \
  -p 5432:5432 \
  -d postgres
```