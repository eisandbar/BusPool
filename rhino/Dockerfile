FROM golang:1.19
WORKDIR /rhino
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN go mod download && go build -o ./
CMD ["./rhino"]