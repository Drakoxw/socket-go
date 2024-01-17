# GO SOCKET IO

## RUN
```bash
go run main.go
```

## CONEXION: 
> ws://localhost:8080/ws?id=666&room=strategy-rom

## DOCKER
```bash
docker build -t socket:lastest .
```

```bash
docker run -p 8080:8080 socket:lastest
```

```bash
docker compose up --build -d
```