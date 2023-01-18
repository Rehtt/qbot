FROM golang:1.19-alpine AS builder

RUN go env -w GO111MODULE=auto \
  && go env -w CGO_ENABLED=0 \
  && go env -w GOPROXY=https://goproxy.cn,direct

WORKDIR /build

RUN set -ex  \
    && sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories \
    && apk add --no-cache git jq curl \
    && latest=$(curl https://api.github.com/repos/Mrs4s/go-cqhttp/releases/latest | jq -r '.tag_name') \
    && git clone -b ${latest} --depth=1 https://github.com/Mrs4s/go-cqhttp.git \
    && cd go-cqhttp \
    && go build -ldflags "-s -w -extldflags '-static'" -o cqhttp

FROM alpine:latest

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories \
    && apk add --no-cache ffmpeg tzdata

COPY --from=builder /build/go-cqhttp/cqhttp /usr/bin/cqhttp
RUN chmod +x /usr/bin/cqhttp

WORKDIR /bot

ENTRYPOINT [ "/usr/bin/cqhttp" ]