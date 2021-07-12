module github.com/eventflowdb/eventflowdb

go 1.16

require (
	github.com/dgraph-io/badger/v3 v3.2103.0
	github.com/gofiber/adaptor/v2 v2.1.7
	github.com/gofiber/fiber/v2 v2.13.0
	github.com/google/uuid v1.2.0
	github.com/jackc/pgx/v4 v4.12.0 // indirect
	github.com/joho/godotenv v1.3.0
	github.com/kr/pretty v0.2.0 // indirect
	github.com/prometheus/client_golang v1.11.0
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	go.etcd.io/bbolt v1.3.6
	google.golang.org/genproto v0.0.0-20210624195500-8bfb893ecb84 // indirect
	google.golang.org/grpc v1.38.0
	google.golang.org/protobuf v1.27.0
	gorm.io/driver/postgres v1.1.0
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.21.11
)
