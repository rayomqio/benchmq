package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rayomqio/benchmq/internal/bench"
	"github.com/rayomqio/benchmq/pkg/logger"
	"github.com/spf13/cobra"
)

var connCmd = &cobra.Command{
	Use:   "conn",
	Short: "Run a connection benchmark against the configured MQTT broker.",
	Long:  `Opens N concurrent MQTT connections (from config or flags) to measure connection throughput, failures, and timing.`,
	Run: func(cmd *cobra.Command, args []string) {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		// Parse flags
		host, err := cmd.Flags().GetString("host")
		if err != nil {
			logger.Error("failed to parse host flag", logger.ErrorAttr(err))
			return
		}

		port, err := cmd.Flags().GetUint16("port")
		if err != nil {
			logger.Error("failed to parse port flag", logger.ErrorAttr(err))
			return
		}

		clients, err := cmd.Flags().GetInt("clients")
		if err != nil {
			logger.Error("failed to parse clients flag", logger.ErrorAttr(err))
			return
		}

		delay, err := cmd.Flags().GetInt("delay")
		if err != nil {
			logger.Error("failed to parse delay flag", logger.ErrorAttr(err))
			return
		}

		clean, err := cmd.Flags().GetBool("clean")
		if err != nil {
			logger.Error("failed to parse clean flag", logger.ErrorAttr(err))
			return
		}

		keepalive, err := cmd.Flags().GetUint16("keepalive")
		if err != nil {
			logger.Error("failed to parse keepalive flag", logger.ErrorAttr(err))
			return
		}

		clientID, err := cmd.Flags().GetString("clientID")
		if err != nil {
			logger.Error("failed to parse clientID flag", logger.ErrorAttr(err))
			return
		}

		username, err := cmd.Flags().GetString("username")
		if err != nil {
			logger.Error("failed to parse username flag", logger.ErrorAttr(err))
			return
		}

		password, err := cmd.Flags().GetString("password")
		if err != nil {
			logger.Error("failed to parse password flag", logger.ErrorAttr(err))
			return
		}

		// Create benchmark
		b, err := bench.NewBenchmark(
			Cfg,
			bench.WithClients(clients),
			bench.WithDelay(delay),
			bench.WithCleanSession(clean),
			bench.WithKeepAlive(keepalive),
			bench.WithClientID(clientID),
			bench.WithUsername(username),
			bench.WithPassword(password),
			bench.WithHost(host),
			bench.WithPort(port),
		)
		if err != nil {
			logger.Error("failed to create benchmark", logger.ErrorAttr(err))
			return
		}

		// Run benchmark in a goroutine so we can wait for shutdown
		done := make(chan struct{})
		go func() {
			b.RunConnections()
			close(done)
		}()

		select {
		case <-sigs:
			logger.Info("received shutdown signal", logger.State("interrupted"))
			return
		case <-done:
			logger.Info("connection benchmark completed", logger.State("completed"))
		}
	},
}

func init() {
	rootCmd.AddCommand(connCmd)

	// Register flags
	connCmd.Flags().IntP("clients", "c", 100, "Number of concurrent clients to connect")
	connCmd.Flags().IntP("delay", "d", 1000, "Delay between each client connection in milliseconds")
}
