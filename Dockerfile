FROM golang:1.17 as builder


WORKDIR /app

COPY /main ./main
COPY /services ./services
COPY /poller ./poller
COPY go.mod .
COPY go.sum .
RUN go mod tidy
RUN ls -l /app
RUN go build -o myapp ./main/cmd.go
RUN ls -l /app

FROM alpine:3.14
# This line fixed alpine problem of finding folders 
RUN apk add libc6-compat
WORKDIR /app
COPY --from=builder /app/myapp .
RUN chmod +x myapp
RUN ls -l /app
EXPOSE 8080

# Command to run the executable
CMD ["./myapp"]
# RUN ./myapp