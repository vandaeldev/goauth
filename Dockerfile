FROM golang:1.24-alpine

RUN apk --no-cache update

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -v -o /usr/local/bin/app ./

EXPOSE 8888

CMD ["app"]
