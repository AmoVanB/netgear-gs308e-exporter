ARG ARCH=
FROM ${ARCH}golang:1.15
WORKDIR /opt/netgear-gs308e-exporter
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -mod=vendor -o /go/bin/netgear-gs308e-exporter

FROM ${ARCH}alpine:3.8
WORKDIR /usr/local/bin
COPY --from=0 /go/bin/netgear-gs308e-exporter .
COPY ./config/config.yaml /etc/netgear-gs308e-exporter/config.yaml

ENTRYPOINT ["/usr/local/bin/netgear-gs308e-exporter"]
CMD ["serve"]
