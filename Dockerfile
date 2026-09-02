FROM alpine:3.24.1

COPY sm /usr/local/bin/sm
RUN chmod +x /usr/local/bin/sm

RUN mkdir /workdir
WORKDIR /workdir

ENTRYPOINT [ "/usr/local/bin/sm" ]