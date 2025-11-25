package cmd

import (
	"os"
	"strings"

	"github.com/rayomqio/benchmq/pkg/config"
	"github.com/rayomqio/benchmq/pkg/logger"
	"github.com/spf13/cobra"
)

var Cfg *config.Config

var rootCmd = &cobra.Command{
	Use:   "benchmq",
	Short: "BenchMQ is a simple, fast, and lightweight CLI to benchmark your MQTT broker with ease.",
	Long:  `BenchMQ is a simple, fast, and open-source CLI tool for benchmarking MQTT brokers. Measure throughput, latency, and stability of your MQTT setup with ease.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cfg, err := config.InitializeCfg()
	if err != nil {
		logger.InitGlobalLogger(logger.DevelopmentConfig())
		logger.Error("Failed to initialize config", logger.ErrorAttr(err))
		os.Exit(1)
	}
	if cfg == nil {
		logger.InitGlobalLogger(logger.DevelopmentConfig())
		logger.Error("Config is nil, this should not happen")
		os.Exit(1)
	}
	Cfg = cfg

	var lcfg logger.Config
	env := strings.ToLower(Cfg.Environment)
	switch env {
	case "production":
		lcfg = logger.ProductionConfig()
	case "development":
		lcfg = logger.DevelopmentConfig()
	default:
		lcfg = logger.DevelopmentConfig()
	}

	lcfg.Service = "benchmq"
	lcfg.Version = Cfg.Version
	lcfg.Environment = Cfg.Environment
	logger.InitGlobalLogger(lcfg)
	if Cfg.Environment != "production" && Cfg.Environment != "development" {
		logger.Warn("Invalid server environment config value, assigning default development.", logger.String("environment", Cfg.Environment))
	}

	rootCmd.PersistentFlags().StringP("host", "H", "localhost", "Hostname or IP address of the MQTT broker")
	rootCmd.PersistentFlags().Uint16P("port", "P", 1883, "Port number of the MQTT broker")
	rootCmd.PersistentFlags().StringP("clientID", "i", "benchmq-client", "Client ID for MQTT connections")
	rootCmd.PersistentFlags().BoolP("clean", "x", true, "Clean previous session when connecting")
	rootCmd.PersistentFlags().Uint16P("keepalive", "k", 60, "Keepalive interval in seconds")
	rootCmd.PersistentFlags().StringP("username", "u", "", "Username for MQTT connections")
	rootCmd.PersistentFlags().StringP("password", "p", "", "Password for MQTT connections")
}
