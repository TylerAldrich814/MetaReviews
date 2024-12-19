sql_dev:
	@docker exec -i MetaMovieContainer mysql \
		-h localhost -P 3606                   \
		-protocol=tcp                          \
		-uroot                                 \
		-proot_password                        \
		MetaMovieDB < $(FILE)

sql_schema:
	@docker exec -i MetaMovieContainer mysql \
		-uroot                                 \
		-p$(PASSW)                             \
		MetaMovieDB < $(FILE)

sql_show_tables:
	@docker exec -i MetaMovieContainer mysql \
		MetaMovieDB                            \
		-h localhost -P 3606                   \
		-protocol=tcp                          \
		-uroot                                 \
		-proot_password                        \
		-e "SHOW tables"

help:
	@echo "Usage:"
	@echo "  ## sql_dev: For running MySQL files on the development DB."
	@echo " make sql_dev FILE=path/to/sql/schema.sql"
	@echo " make sql_schema PASSW=<password> FILE=path/to/sql/schema.sql"
