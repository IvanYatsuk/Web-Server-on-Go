FROM golang:1.17-alpine
WORKDIR /Web-Server-on-Go
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o main .
CMD ["./main"]