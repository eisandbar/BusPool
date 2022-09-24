FROM golang:1.19
WORKDIR /client
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN go build -o ./
CMD ["./client"]