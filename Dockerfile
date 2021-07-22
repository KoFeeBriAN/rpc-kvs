#build stage
FROM golang:alpine AS builder
RUN apk add --no-cache git
WORKDIR /go/src/app
COPY . .
RUN go get -d -v .
RUN go build -o /go/bin/frontend -v ./cmd/frontend/main.go
RUN go build -o /go/bin/storage -v ./cmd/storage/main.go

#final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN apk add --no-cache bash
COPY --from=builder /go/bin/frontend /frontend
COPY --from=builder /go/bin/storage /storage
COPY --from=builder /go/src/app/wrapper_script.sh /wrapper_script.sh

COPY --from=builder /go/src/app/config /config

RUN ["chmod", "+x", "/frontend"]
RUN ["chmod", "+x", "/storage"]
RUN ["chmod", "+x", "/wrapper_script.sh"]

CMD [ "./wrapper_script.sh" ]

LABEL Name=monolith Version=0.0.1
EXPOSE 8000
EXPOSE 50051
