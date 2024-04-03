# Create and Run PostgreSQL Container and Initialize Database
```bash
docker run --name postgres-tictactoe \
  -e POSTGRES_PASSWORD=mysecretpassword \
  -e POSTGRES_DB=tictactoe \
  -v ./init.sql \
  -p 5432:5432 \
  -d postgres
```



# Next Steps
1. Get Server Working
2. Test Backend Docker Container (Need to make sure postgres container is up first)
3. Get Frontend Container Working
4. Test Docker Compose
5. Modify helm charts to actually use the containers
6. Create a ci/cd pipeline for testing and deploying the application to a kubernetes cluster on GCP