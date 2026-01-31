# Nipple exposer (SDR)

## Building

`go build -o nipple`

## Running

`./nipple`

For local deployment either change config.yaml

```yml
rcon:
  host: 127.0.0.1:27015
  password:  rcon_password

http:
  port: port

server:
  password: server_password
```

Or pass through env variables

```env
RCON_HOST=ip:port
RCON_PASSWORD=password
HTTP_SERVER_PORT=port
SERVER_PASSWORD=password
```
