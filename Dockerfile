FROM golang:alpine as builder
RUN apk --no-cache add build-base tzdata ca-certificates
WORKDIR /app
COPY go.mod go.sum /app/
RUN go mod download
COPY . /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags=go_json -o dist/app .
RUN mkdir -p /app/storage

FROM gcr.io/distroless/base-debian12 as serve
WORKDIR /app

COPY --from=builder /app/dist/app /app/dist/app
COPY --from=builder /app/storage /app/storage

ENV GIN_MODE=release
ENV TZ=Asia/Bangkok
EXPOSE 8080
USER nonroot:nonroot

ENTRYPOINT ["/app/dist/app"]
CMD ["http"]