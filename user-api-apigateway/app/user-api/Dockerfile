##
## Build
##
FROM golang:1.18.1-alpine3.15  AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

##
## Deploy
##
FROM alpine:3.15.4
WORKDIR /root/
COPY --from=build /app/main .
COPY --from=build /app/.env .
EXPOSE 8000
ENTRYPOINT ["./main"]