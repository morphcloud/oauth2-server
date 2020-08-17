FROM golang:alpine as builder
ENV GO111MODULE=on
WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./main.go

FROM scratch
COPY --from=builder /build/app /app
COPY .env .
EXPOSE 9096
ENTRYPOINT ["/app"]
