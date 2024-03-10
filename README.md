# Yet another notes app
___
This is learning project. So take everything with grain of salt and don't emulate things without a thought.

This is [go](https://go.dev/) + [gin](https://gin-gonic.com/) implementation of backend
For more details about project read frontend [README.md](https://github.com/KonradOliwer/yana-fe-react/)


### TODO for production ready
In case of productions configuration (like dp user and password) should at least land into environment variables. Alternatively there can be a configuration file.
If we want higher level of security we might want to look for more sophisticated tools - like using AWS secret (if our infrastructure is on AWS).

## Using app
### Requirements
- [go](https://go.dev/)
- [docker](https://www.docker.com/)
### Running
Start db (this will drop DB on finishing process)
```bash
docker run --name postgres -e POSTGRES_PASSWORD=password -e POSTGRES_USER=user -e POSTGRES_DB=yana -p 5432:5432 --rm postgres
```

Run with debug mode
```bash
go run main.go
```

### Testing
```bash
bash run_test.sh
```