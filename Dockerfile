FROM golang:1.23-alpine as build

RUN addgroup -S app && adduser -S app -G app

WORKDIR /go/src/app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -trimpath -ldflags '-extldflags "-static"' -tags timetzdata -o app ./cmd/web

FROM scratch

COPY --from=build /go/src/app/app /bin/app

USER app

ENTRYPOINT ["/bin/app"]