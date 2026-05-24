FROM alpine:latest

WORKDIR /app

COPY main .
COPY frontend ./frontend

EXPOSE 8080

CMD ["./main"]