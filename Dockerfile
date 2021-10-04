FROM postgis/postgis:13-3.1-alpine
COPY var/data /var/lib/postgresql/data
