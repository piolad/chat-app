package config

import "os"

type Config struct {
	GRPCPort              string
	MongoConnectionString string
	MongoDBName           string
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func Load() Config {
	return Config{
		GRPCPort:              getenv("GRPC_PORT", ":50051"),
		MongoConnectionString: getenv("MONGO_CONNECTION_STRING", "mongodb://adminUser:adminPassword@message-data-centre-db:27017"),
		MongoDBName:           getenv("MONGO_DB_NAME", "message-db"),
	}
}
