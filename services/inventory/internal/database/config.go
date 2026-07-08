package database

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

var DBConfig = Config{
	Host:     "localhost",
	Port:     "5433",
	User:     "postgres",
	Password: "postgres",
	DBName:   "inventorydb",
	SSLMode:  "disable",
}