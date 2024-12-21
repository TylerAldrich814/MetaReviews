## Docker/consul.Dockerfile

FROM mysql:8.0

ENV MYSQL_ROOT_PASSWORD=root_password
ENV MYSQL_DATABASE=metadb
ENV MYSQL_USER=metauser
ENV MYSQL_PASSWORD=metapass

EXPOSE 3306
