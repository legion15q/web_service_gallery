FROM golang:alpine AS build
#/go/src/... потому что там лежат все файлы go
#WORKDIR автоматически создает папку 
WORKDIR /go/src/my_app_in_wm
COPY ./ ./
# Можно сгенерировать go.mod, но он был уже скопирован
RUN go mod init
RUN go mod tidy
RUN go get ./...
# Создаем папку, откуда будут подтягиваться картины
WORKDIR /images/
# Возвращаемся в исходную
WORKDIR /go/src/my_app_in_wm
RUN go build -o main .
#порт контейнера, который слушает main.go. Не зависит от порта в браузере
EXPOSE 8080
# Лучше использовать ENTRYPOINT, вместо CMD
ENTRYPOINT [ "./main" ]
# docker build -t go-app .
# docker run --name=go-web-app -p 8080:80 go-app 
# То есть слушаем порт 80 контейнера, который маппится с портом 8080 из браузера