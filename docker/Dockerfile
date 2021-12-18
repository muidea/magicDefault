FROM alpine:3.13.4

LABEL Author="rangh"
LABEL Application="magicDefault"

RUN mkdir -p /var/app

WORKDIR /var/app

ADD magicDefault /var/app/magicDefault
ADD bootstrap.sh /var/app/bootstrap.sh

RUN chmod +x /var/app/magicDefault
RUN chmod +x /var/app/bootstrap.sh

ENTRYPOINT [ "/var/app/bootstrap.sh" ]

EXPOSE 8880
