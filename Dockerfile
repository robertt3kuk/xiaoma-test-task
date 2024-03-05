FROM golang:latest AS builder
WORKDIR /server/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app cmd/app/main.go 

FROM gcr.io/distroless/static
WORKDIR /bin/
COPY --from=builder /server/app .
COPY --from=builder /server/config/ ./config/
CMD [ "./app" ]
EXPOSE 8000
