
$echo lsof -i:3000 |awk 'NR==2' | awk '{print $2}'
lsof -i:3000 |awk 'NR==2' | awk '{print $2}' |xargs kill -INT
#PREST_DEBUG=true PREST_AUTH_ENABLED=false PREST_PG_HOST=localhost PREST_PG_USER=postgres PREST_PG_PASS=123456 PREST_PG_DATABASE=prest PREST_PG_PORT=5432 PORT=6000 go run main.go

PREST_CONF=/Users/andy/GoLang/src/prest/adapters/prest.toml  go run main.go

