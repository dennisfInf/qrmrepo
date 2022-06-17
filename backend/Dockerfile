FROM golang:1.18 as build-env

ENV CGO_ENABLED=0

COPY . /app/
WORKDIR /app/

RUN go build -o backend -ldflags="-extldflags=-static" .

FROM gcr.io/distroless/static

COPY --from=build-env --chown=nonroot:nonroot \
    /app/backend \
    /app/

EXPOSE 80 443

USER nonroot:nonroot

ENTRYPOINT [ "/app/backend" ]
CMD [ "serve" ]