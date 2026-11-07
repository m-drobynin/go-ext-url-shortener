# migrate cmd example

migrate create -ext sql -dir ./migrations -seq create_url_table   
migrate -database "postgres://mdrobynin:@localhost:5432/mdrobynin?sslmode=disable" -path ./migrations up


# run with db example

go run go run cmd/shortener/main.go -d "host=localhost user=mdrobynin password= dbname=mdrobynin sslmode=disable"
go run go run cmd/shortener/main.go -f "./file_storage.json"
