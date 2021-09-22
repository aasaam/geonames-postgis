FROM postgis/postgis:13-3.1
COPY temp/postgresql /var/lib/postgresql/data
