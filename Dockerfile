FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY bigip_exporter .
EXPOSE 9142
ENTRYPOINT ["/bigip_exporter"]
