FROM docker.io/alpine:latest

COPY entrypoint.sh /entrypoint.sh
RUN apk add --no-cache git && \
    wget -O /usr/local/bin/git-chglog \
        https://github.com/git-chglog/git-chglog/releases/download/0.9.1/git-chglog_linux_amd64 && \
    chmod 755 /entrypoint.sh /usr/local/bin/git-chglog

ENTRYPOINT [ "/entrypoint.sh" ]