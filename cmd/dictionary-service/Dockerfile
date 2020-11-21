FROM scratch

MAINTAINER Adrian Lapierre <al@alapierre.io>

EXPOSE 9098

ADD dictionary-service *.*.sql /

COPY resources/passwd /etc/passwd
COPY resources/groups /etc/groups

USER app

CMD ["/dictionary-service"]
