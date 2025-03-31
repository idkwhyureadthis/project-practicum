FROM gomicro/goose

COPY *.sql /migrations/
COPY entrypoint.sh /migrations/

RUN ls -la /migrations

RUN chmod +x /migrations/entrypoint.sh

ENTRYPOINT ["sh", "/migrations/entrypoint.sh"]