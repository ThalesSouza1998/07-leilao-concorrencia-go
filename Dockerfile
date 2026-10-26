FROM golang:1.22-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /app/auction ./cmd/auction

FROM alpine:3.20

WORKDIR /app

COPY --from=build /app/auction .

EXPOSE 8080

ENTRYPOINT ["./auction"]
