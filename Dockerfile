FROM golang:onbuild
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go build -o renderer .
CMD ["/app/renderer"]

FROM alpine:3.7
RUN apk add --no-cache ca-certificates
RUN mkdir /app /volumes_in /volumes_out
COPY --from=0 /app/renderer /app/renderer
WORKDIR /app
ENTRYPOINT ["/app/renderer"]
