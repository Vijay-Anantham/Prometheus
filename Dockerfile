FROM golang:1.17 as builder


WORKDIR /app

COPY ./main/cmd.go .
COPY go.mod .
COPY go.sum .
RUN go mod tidy
RUN go build -o myapp .
RUN ls -l /app

FROM alpine:3.14
# This line fixed alpine problem of finding folders 
RUN apk add libc6-compat
WORKDIR /app
COPY --from=builder /app/myapp .
RUN chmod +x myapp
RUN ls -l /app
EXPOSE 8090

# Command to run the executable
CMD ["./myapp"]
# RUN ./myapp