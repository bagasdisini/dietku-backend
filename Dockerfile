# Step #1
FROM golang:1.22 as firststage
LABEL description="Dietku Backend API"
LABEL maintainer="Bagas <mbagas221@gmail.com>"
WORKDIR /build/
COPY . /build
ENV CGO_ENABLED=0
RUN go get
RUN go build -o dietku-backend

# Step #2
FROM alpine:latest
WORKDIR /app/
RUN apk update
RUN apk upgrade
RUN apk add ca-certificates && update-ca-certificates
RUN apk add --no-cache tzdata gcompat
ENV TZ=Asia/Jakarta
COPY --from=firststage /build/dietku-backend .
CMD ["./dietku-backend"]

