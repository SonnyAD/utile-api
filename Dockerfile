FROM golang:latest AS build
WORKDIR /go/src

COPY go.mod .
COPY *.go ./
RUN go get
RUN CGO_ENABLED=0 go build -o server .

FROM scratch AS runtime
COPY --from=build /go/src/server /
EXPOSE 3000/tcp
CMD ["/server"]
