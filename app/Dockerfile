FROM golang:1.21.5

WORKDIR /webapp
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/main.go
CMD ["ls"]
CMD ["/main.exe"]