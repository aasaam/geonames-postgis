FROM postgis/postgis:14-3.1-alpine
COPY var/data /var/lib/postgresql/data
