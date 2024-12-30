FROM library/golang:1.22.5-alpine as builder

WORKDIR /app

COPY go.mod go.mod

RUN go mod download

COPY . .

RUN go build -o bin/server ./cmd/server

FROM alpine

COPY --from=builder /app/bin/server /server

CMD [ "/server" ]
