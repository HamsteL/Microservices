##
## Build
##
FROM golang:1.18.1-alpine3.15  AS build
WORKDIR /app
COPY . .
RUN go mod download
COPY *.go ./
RUN go build -o /k8s_first_try

##
## Deploy
##
FROM alpine:3.15.4
WORKDIR /
COPY --from=build /k8s_first_try /k8s_first_try
EXPOSE 8000
ENTRYPOINT ["./k8s_first_try"]