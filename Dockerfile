# syntax=docker/dockerfile:1

FROM golang:1.22.5 AS build-stage
  WORKDIR /app

  COPY go.mod go.sum ./
  RUN go mod download

  COPY . .
  RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd/main.go

# Deploy the application binary into a lean image
FROM scratch AS build-release-stage
  WORKDIR /

  COPY --from=build-stage /api /api
  COPY development.env .

  EXPOSE 8080

  ENTRYPOINT ["/api"]