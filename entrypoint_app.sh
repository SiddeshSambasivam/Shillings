DB_HOST="db"
DB_PORT="3306"

wait-for "${DB_HOST}:${DB_PORT}" -- "$@"

CompileDaemon --build="go build -buildvcs=false -o bin/app ./cmd/ApplicationServer" --command=./bin/app
