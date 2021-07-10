FROM golang:latest AS build
WORKDIR /go/src

COPY go.mod .
COPY server.go .
RUN go get
RUN CGO_ENABLED=0 go build server.go

FROM scratch AS runtime
COPY --from=build /go/src/server /
EXPOSE 3000/tcp
CMD ["/server"]
