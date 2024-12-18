# Generates apo/movie.proto
gen_movie:
	@protoc --go_out=.                                    \
		--go_opt=paths=source_relative                      \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    api/movie.proto
