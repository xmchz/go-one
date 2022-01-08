package storage

type Config interface {
	DbDriver() string
	DbName() string
	DbUrl() string //dsn := "user:password@tcp(127.0.0.1:3306)/sql_test?charset=utf8mb4&parseTime=True"
	DbInfo() string
	MaxIdleConn() int
	MaxOpenConn() int
	MigrationUrl() string
	MigrationVersion() uint
}
