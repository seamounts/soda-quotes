FROM golang:1.13-alpine as builder
ENV GOPROXY=https://goproxy.io 
RUN apk add gcc g++ make libffi-dev openssl-dev libtool
COPY . $GOPATH/src/soda-quetes
RUN cd $GOPATH/src/soda-quetes && \
    go build -o /app/soda-quetes && \
    cp -r db /app/db && cp -r views /app/views && cp -r statics /app/statics

FROM alpine:3.10.3
WORKDIR /app
COPY --from=builder /app/soda-quetes /app/
COPY --from=builder /app/db /app/db
COPY --from=builder /app/views /app/views
COPY --from=builder /app/statics /app/statics

ENTRYPOINT [ "./soda-quetes" ]



