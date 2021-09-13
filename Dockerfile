FROM golang:alpine AS build
WORKDIR /go/src

RUN apk add -U --no-cache ca-certificates

COPY go.mod .
COPY *.go ./
RUN go get
RUN CGO_ENABLED=0 go build -o server .


FROM scratch AS runtime

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/src/server /

EXPOSE 3000/tcp
CMD ["/server"]
