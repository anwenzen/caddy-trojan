FROM golang:1.22-alpine AS builder
COPY ./ /caddy-trojan
RUN apk add git \
    && go install github.com/caddyserver/xcaddy/cmd/xcaddy@latest \
    && xcaddy build --output /caddy \
        --with github.com/imgk/caddy-trojan \
        --replace github.com/imgk/caddy-trojan=/caddy-trojan


FROM alpine:latest AS dist

# See https://caddyserver.com/docs/conventions#file-locations for details
ENV XDG_CONFIG_HOME /config
ENV XDG_DATA_HOME /data
ENV TZ Asia/Shanghai

COPY --from=builder /caddy /usr/bin/caddy
ADD https://raw.githubusercontent.com/caddyserver/dist/master/config/Caddyfile /etc/caddy/Caddyfile
ADD https://raw.githubusercontent.com/caddyserver/dist/master/welcome/index.html /usr/share/caddy/index.html

RUN apk add tzdata \
    && ln -sf /usr/share/zoneinfo/${TZ} /etc/localtime \
    && echo ${TZ} > /etc/timezone \
    && rm -rf /var/cache/apk/*

VOLUME /config
VOLUME /dataclear

EXPOSE 80
EXPOSE 443
EXPOSE 443/udp
EXPOSE 2019

WORKDIR /srv

CMD ["caddy", "run", "--config", "/etc/caddy/Caddyfile", "--adapter", "caddyfile"]