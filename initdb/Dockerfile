FROM postgres:13.1
ENV POSTGRES_USER docker
ENV POSTGRES_DB docker
ENV POSTGRES_PASSWORD passw0rd
VOLUME /var/lib/postgresql/data
COPY init.sql /docker-entrypoint-initdb.d/
