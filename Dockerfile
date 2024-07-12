FROM golang:1.22 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . ./

WORKDIR /app/cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o /api

# Deploy the application binary into a lean image
FROM alpine:3.19 AS build-release-stage

WORKDIR /

COPY --from=build-stage /api /api

ENTRYPOINT ["/api"]
