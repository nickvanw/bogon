FROM alpine:3.3
MAINTAINER Nick Van Wiggeren nick@facepwn.com

RUN apk --update add ca-certificates
COPY bogon /bogon

CMD []
ENTRYPOINT ["/bogon"]