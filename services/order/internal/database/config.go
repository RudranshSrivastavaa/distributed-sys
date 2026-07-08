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
	Port:     "5432",
	User:     "postgres",
	Password: "postgres",
	DBName:   "orderdb",
	SSLMode:  "disable",
}