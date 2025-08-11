### Stage 1: Build app binary ###
# base Go binary
FROM golang:1.24 AS builder

# create dir for subsequent commands
WORKDIR /app

### copy project files from local dir into /app inside container ###
COPY go.mod ./
COPY go.sum ./
RUN go mod download
# copy all source code
COPY . .

# compile Go code into binary named 'tictacgo'
RUN go build -o tictacgo

# listens on port 8080
EXPOSE 8080

# run
CMD ["./tictacgo"]