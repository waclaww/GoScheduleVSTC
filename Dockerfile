FROM golang:1.26-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /out/api ./cmd/api \
 && CGO_ENABLED=0 go build -o /out/migrate ./cmd/migrate

FROM alpine:3.22

RUN apk add --no-cache curl ca-certificates tzdata

WORKDIR /app

COPY --from=build /out/ ./

EXPOSE 8000

CMD ["sh", "-c", "./migrate && ./api"]
