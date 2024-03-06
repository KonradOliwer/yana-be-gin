## Running
Start db
```bash
docker run --name postgres -e POSTGRES_PASSWORD=password -e POSTGRES_USER=user -e POSTGRES_DB=yana -p 5432:5432 --rm postgres
```

Run GO app
```bash
go run -C src main.go
```

## Testing
```bash
bash run_test.sh
```