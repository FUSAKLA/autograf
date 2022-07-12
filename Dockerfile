FROM alpine

COPY autograf /usr/bin/autograf
COPY Dockerfile /

ENTRYPOINT ["autograf"]
CMD ["--help"]
