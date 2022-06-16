FROM golang:1.18 as build-env

ENV CGO_ENABLED=0

COPY . /app/
WORKDIR /app/

RUN go build -o relay -ldflags="-extldflags=-static" .

FROM gcr.io/distroless/static

COPY --from=build-env --chown=nonroot:nonroot \
    /app/relay \
    /app/

USER nonroot:nonroot

ENTRYPOINT [ "/app/relay" ]
CMD [ "serve" ]