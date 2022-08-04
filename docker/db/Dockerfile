FROM mysql:8.0

ENV MYSQL_ROOT_PASSWORD=root
ENV MYSQL_DATABASE=shillings
ENV MYSQL_USER=admin
ENV MYSQL_PASSWORD=admin

COPY ./docker/db/init.sql /docker-entrypoint-initdb.d/init.sql

