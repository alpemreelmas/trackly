package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type AppConfig struct {
	Port                  string `mapstructure:"port" yaml:"port"`
	CouchbaseUrl          string `mapstructure:"couchbase_url" yaml:"couchbase_url"`
	CouchbaseUsername     string `mapstructure:"couchbase_username" yaml:"couchbase_username"`
	CouchbasePassword     string `mapstructure:"couchbase_password" yaml:"couchbase_password"`
	AzureConnectionString string `mapstructure:"azure_connection_string" yaml:"azure_connection_string"`
	CosmosDBEndpoint      string `mapstructure:"cosmosdb_endpoint" yaml:"cosmosdb_endpoint"`
	CosmosDBKey           string `mapstructure:"cosmosdb_key" yaml:"cosmosdb_key"`
	CosmosDBDatabase      string `mapstructure:"cosmosdb_database" yaml:"cosmosdb_database"`
	CosmosDBContainer     string `mapstructure:"cosmosdb_container" yaml:"cosmosdb_container"`
}

func Read() *AppConfig {
	viper.SetConfigName("config")      // name of config file (without extension)
	viper.SetConfigType("yaml")        // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("$PWD/config") // call multiple times to add many search paths
	viper.AddConfigPath(".")           // optionally look for config in the working directory
	viper.AddConfigPath("/config")     // optionally look for config in the working directory
	viper.AddConfigPath("./config")    // optionally look for config in the working directory
	err := viper.ReadInConfig()        // Find and read the config file
	if err != nil {                    // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	var appConfig AppConfig
	err = viper.Unmarshal(&appConfig)
	if err != nil {
		panic(fmt.Errorf("fatal error unmarshalling config: %w", err))
	}

	return &appConfig
}
