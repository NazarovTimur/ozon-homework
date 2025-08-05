package config

import (
	"os"
)

const (
	EnvPort            = "PORT"
	GrpcPort           = "GRPC_PORT"
	DatabaseMasterDSN  = "MASTER_DB_DSN"
	DatabaseReplicaDSN = "REPLICA_DB_DSN"
)

type Config struct {
	Server   ServerConfig
	Grpc     GrpcConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Port string
}

type GrpcConfig struct {
	Port string
}

type DatabaseConfig struct {
	MasterDSN  string
	ReplicaDSN string
}

func NewConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv(EnvPort, "PORT"),
		},
		Grpc: GrpcConfig{
			Port: getEnv(GrpcPort, "GRPC_PORT"),
		},
		Database: DatabaseConfig{
			MasterDSN:  getEnv(DatabaseMasterDSN, "postgres://postgres:masterpass@localhost:5432/appdb?sslmode=disable"),
			ReplicaDSN: getEnv(DatabaseReplicaDSN, "postgres://postgres:replicapass@localhost:5432/appdb?sslmode=disable"),
		},
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
