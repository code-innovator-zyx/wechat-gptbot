FROM golang:1.20-alpine3.18 AS builder
RUN apk --update --no-cache add musl-dev gcc
RUN go env -w GOPROXY=https://goproxy.cn,direct
WORKDIR /code
COPY . /code
ENV CGO_ENABLED=1
RUN go fmt && go mod tidy

RUN go build .


FROM 1003941268/python3.11-alpine:latest
RUN sed -i "s/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g" /etc/apk/repositories
RUN apk --update --no-cache add tzdata
# 设置时区为上海
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

WORKDIR /code
COPY --from=builder /code/wechat-gptbot /code/wechat-gptbot
COPY --from=builder /code/streamlit_app/ /code/streamlit_app/
COPY --from=builder /code/.streamlit/ /code/.streamlit/
COPY --from=builder /code/config/config.json.template /code/config/config.json
COPY --from=builder /code/config/prompt.conf.template /code/config/prompt.conf

RUN chown -R root:root /code
CMD [ "/code/wechat-gptbot" ]