# GO SOCKET IO

Este proyecto es un servidor de chat implementado en Go utilizando el framework Echo para la gestión de rutas HTTP y el paquete Gorilla WebSocket para la funcionalidad de WebSocket. Proporciona una plataforma de comunicación en tiempo real que permite a los clientes conectarse a salas de chat, enviar mensajes y recibir actualizaciones en tiempo real.

Para los servicios que no pueden conectarse con el Websocket, se proporciona un servicio de conexion HTTP para el envio de mensajes.

### Características Principales
- **WebSocket**: Comunicación bidireccional entre el servidor y los clientes utilizando WebSocket para una experiencia de chat en tiempo real.
- **Salas de Chat**: Los usuarios pueden unirse a salas específicas y comunicarse dentro de ellas.
- **Peticiones REST**: Funcionalidad adicional mediante un endpoint REST para enviar mensajes desde sistemas que no admiten WebSocket directamente.

## Instalacion
```bash
go get -u ./...
```

```bash
go mod tidy
```

## RUN
```bash
go run main.go
```

## CONEXION WS: 
>ws://localhost:8888/ws?id=1&room=rom-1
```
{
	"sender": 1,
	"content": "Hola desde socket"
}
```

## CONEXION HTTP: 
>http://localhost:8888/rest-sw?room=rom-1
```
{
	"sender": 2,
	"content": "Hola desde HTTP"
}
```

## DOCKER
```bash
docker build -t socket:lastest .
```

```bash
docker run -d -p 8888:8888  --name socket-app socket:lastest
```

```bash
docker compose up --build -d
```