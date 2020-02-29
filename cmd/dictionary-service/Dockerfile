FROM scratch

MAINTAINER Adrian Lapierre <al@alapierre.io>

EXPOSE 9098

ADD dictionary-service *.*.sql /

# COPY resources/passwd /etc/passwd
#RUN addgroup -g 666 -S app && adduser -u 666 -S -G app app
#USER app

CMD ["/dictionary-service"]
