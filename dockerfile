FROM golang:1.22.3-alpine
WORKDIR /Web-Server-on-Go
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o main .
CMD ["./main"]