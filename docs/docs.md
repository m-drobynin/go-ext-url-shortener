# migrate cmd example

migrate create -ext sql -dir ./migrations -seq create_url_table   
migrate -database "postgres://mdrobynin:@localhost:5432/mdrobynin?sslmode=disable" -path ./migrations up

# run with db example

go run cmd/shortener/main.go -d "host=localhost user=mdrobynin password= dbname=mdrobynin sslmode=disable"
go run cmd/shortener/main.go -f "./file_storage.json"

# handlers curls

curl -i -X POST 'http://localhost:8080/' -H 'Content-Type: text/plain' --data-raw 'https://example.com/article'

curl -i -X POST 'http://localhost:8080/api/shorten' -H 'Content-Type: application/json' --data-raw '{"url":"https://example.com/news"}'

curl -i -X POST 'http://localhost:8080/api/shorten/batch' -H 'Content-Type: application/json' \
  --data-raw '[
    {
      "correlation_id": "request-1",
      "original_url": "https://example.com/catalog"
    },
    {
      "correlation_id": "request-2",
      "original_url": "https://example.com/about"
    }
  ]'