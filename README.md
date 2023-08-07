# golang-client-server-api

## Pré-Requisitos

- GOLANG 1.19;
- [Composer](https://getcomposer.org);
- [Docker](https://www.docker.com);



## Criar network

* Executar o seguinte comando `docker network create client-server-api-network`

## Subir o container

* Para subir os containers, entrar nas pastas `client` e `server`, executar o seguinte comando `docker compose up -d`

## Atualizar os pacotes

* Executar o seguinte comando (no container) `docker exec -it server-app go mod tidy`

## Para subir o server

* Executar o seguinte comando (no container) `docker exec -it server-app go run main.go`

## Para subir o client

* Executar o seguinte comando (no container) `docker exec -it client-app go run main.go`