FROM golang:alpine as builder
RUN mkdir /build 
ADD . /build/
WORKDIR /build 
ENV GOPROXY https://goproxy.io
RUN apk add --no-cache git
RUN go mod download && go mod verify
RUN go build -o main .

FROM alpine
RUN apk add --no-cache ca-certificates
COPY --from=builder /build/main /app/
EXPOSE 1323
ENTRYPOINT ["/app/main"]
