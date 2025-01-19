FROM golang:1.23.3

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /shirinec ./src/cmd/app/main.go

CMD ["/shirinec"]
