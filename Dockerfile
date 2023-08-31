# FROM golang:latest
 
# WORKDIR /app
 
# # Effectively tracks changes within your go.mod file
# COPY go.mod go.sum .
 
# RUN go mod download
 
# # Copies your source code into the app directory
# COPY . .
 
# RUN go build -o bin/app -v ./cmd/segments-users-service
 
# EXPOSE 8080
 
# CMD [ “/bin/app ]

FROM golang:alpine as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

FROM golang:alpine as builder
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -tags migrate -o /bin/app ./cmd/segments-users-service
    # go build -o bin/app ./cmd/segments-users-service

# RUN chmod +x wait-for.sh

# RUN go mod tidy
FROM scratch
COPY --from=builder /app/config /config
COPY --from=builder /bin/app /app
EXPOSE 8080
CMD ["/app"]