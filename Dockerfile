FROM alpine

RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app

COPY ./bin/pow-test /usr/src/app

EXPOSE 8080

CMD [ "/usr/src/app/pow-test" ]
