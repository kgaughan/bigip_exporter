FROM gcr.io/distroless/static:latest
COPY bigip_exporter .
EXPOSE 9142
ENTRYPOINT ["/bigip_exporter"]
