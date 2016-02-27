FROM alpine:3.3
MAINTAINER Nick Van Wiggeren nick@facepwn.com

RUN apk --update add ca-certificates curl
COPY bogon /bogon

EXPOSE 9001
CMD []
ENTRYPOINT ["/bogon"]