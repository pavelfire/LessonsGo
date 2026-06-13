FROM scratch

COPY bin/portloader /portloader
COPY ports.json /data/ports.json

ENTRYPOINT ["/portloader"]
CMD ["-file", "/data/ports.json"]
