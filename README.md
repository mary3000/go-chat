# go-chat

Simple golang implementation of chat-server using sqlite and gorm

## Usage

### Server

Server lives on `9000` port.

#### With docker

1. `cd chat/cmd/server`
2. `docker-compose up`

#### Without docker

1. `cd chat/cmd/server`
2. `go run main.go` 

### Client

- Use curl-like requests, see [Avito assignment](https://github.com/avito-tech/backend-trainee-assignment#%D0%BE%D1%81%D0%BD%D0%BE%D0%B2%D0%BD%D1%8B%D0%B5-api-%D0%BC%D0%B5%D1%82%D0%BE%D0%B4%D1%8B)
- Use cli

#### Cli

#### With docker

1. `cd chat/cmd/client`
2. `bash client.sh`. To force image rebuilding, run `bash client.sh -f`

#### Without docker

1. `cd chat/cmd/client`
2. `go run main.go` 

Under client, you can do commands:
- `/r name` - register user with `name`
- `/l name` - login under `name` user
- `/cr name` - register chat named `name`
- `/cl name` - login into chat named `name`

After you both logged in user and chat, you can start typing messages!

Conversation example:  
<img src="https://i.ibb.co/HF17sJp/Screenshot-2020-08-18-at-16-35-24.png" alt="alice" width="400px"> 
<img src="https://i.ibb.co/gPV8CXH/Screenshot-2020-08-18-at-16-36-01.png" alt="bob" width="400px"> 