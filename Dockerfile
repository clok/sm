FROM alpine:3.14.0

COPY sm /usr/local/bin/sm
RUN chmod +x /usr/local/bin/sm

RUN mkdir /workdir
WORKDIR /workdir

ENTRYPOINT [ "/usr/local/bin/sm" ]