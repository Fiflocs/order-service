package config

import (
    "fmt"
    "gopkg.in/yaml.v3"
    "os"
)

type Config struct {
    HTTP struct {
        Address string `yaml:"address"`
    } `yaml:"http"`
    Database struct {
        Host     string `yaml:"host"`
        Port     int    `yaml:"port"`
        User     string `yaml:"user"`
        Password string `yaml:"password"`
        DBName   string `yaml:"dbname"`
        SSLMode  string `yaml:"sslmode"`
    } `yaml:"database"`
    NATS struct {
        URL      string `yaml:"url"`
        ClusterID string `yaml:"cluster_id"`
        ClientID  string `yaml:"client_id"`
        Subject   string `yaml:"subject"`
    } `yaml:"nats"`
}

func Load() *Config {
    var cfg Config
    
    data, err := os.ReadFile("config.yaml")
    if err != nil {
        // Значения по умолчанию
        cfg.HTTP.Address = ":8080"
        cfg.Database.Host = "localhost"
        cfg.Database.Port = 5432
        cfg.Database.User = "order_user"
        cfg.Database.Password = "order_password"
        cfg.Database.DBName = "orders"
        cfg.Database.SSLMode = "disable"
        cfg.NATS.URL = "nats://localhost:4222"
        cfg.NATS.ClusterID = "test-cluster"
        cfg.NATS.ClientID = "order-service"
        cfg.NATS.Subject = "orders"
        return &cfg
    }
    
    if err := yaml.Unmarshal(data, &cfg); err != nil {
        panic(fmt.Sprintf("Error parsing config: %v", err))
    }
    
    return &cfg
}