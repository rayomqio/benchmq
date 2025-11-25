package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rayomqio/benchmq/internal/bench"
	"github.com/rayomqio/benchmq/pkg/logger"
	"github.com/spf13/cobra"
)

var subCmd = &cobra.Command{
	Use:   "sub",
	Short: "Subscribe to a topic with specified parameters",
	Long: `Subscribe to a topic with specified parameters.

Parameters:
	- host: Hostname or IP address of the broker
	- port: Port number of the broker
	- clientID: Base client ID prefix (each client appends "-<n>")
    - clients: Number of concurrent subscribers
    - qos: Quality of service level (0, 1, 2)
    - topic: Topic to subscribe to
    - clean: Whether to use a clean session
    - keepalive: Keepalive interval in seconds
    - delay: Optional sleep between subscription lifetime checks
    - count: Expected number of messages (used to determine how long to wait)`,
	Run: func(cmd *cobra.Command, args []string) {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(sigs)

		// Parse flags
		host, err := cmd.Flags().GetString("host")
		if err != nil {
			logger.Error("failed to parse host", logger.ErrorAttr(err))
			return
		}

		port, err := cmd.Flags().GetUint16("port")
		if err != nil {
			logger.Error("failed to parse port", logger.ErrorAttr(err))
			return
		}

		clientID, err := cmd.Flags().GetString("clientID")
		if err != nil {
			logger.Error("failed to parse client ID", logger.ErrorAttr(err))
			return
		}

		clients, err := cmd.Flags().GetInt("clients")
		if err != nil {
			logger.Error("failed to parse number of clients", logger.ErrorAttr(err))
			return
		}

		delay, err := cmd.Flags().GetInt("delay")
		if err != nil {
			logger.Error("failed to parse delay", logger.ErrorAttr(err))
			return
		}

		count, err := cmd.Flags().GetInt("count")
		if err != nil {
			logger.Error("failed to parse message count", logger.ErrorAttr(err))
			return
		}

		topic, err := cmd.Flags().GetString("topic")
		if err != nil {
			logger.Error("failed to parse topic", logger.ErrorAttr(err))
			return
		}

		qos, err := cmd.Flags().GetUint16("qos")
		if err != nil {
			logger.Error("failed to parse QoS", logger.ErrorAttr(err))
			return
		}

		cleanSession, err := cmd.Flags().GetBool("clean")
		if err != nil {
			logger.Error("failed to parse clean session flag", logger.ErrorAttr(err))
			return
		}

		username, err := cmd.Flags().GetString("username")
		if err != nil {
			logger.Error("failed to parse username", logger.ErrorAttr(err))
			return
		}

		password, err := cmd.Flags().GetString("password")
		if err != nil {
			logger.Error("failed to parse password", logger.ErrorAttr(err))
			return
		}

		keepalive, err := cmd.Flags().GetUint16("keepalive")
		if err != nil {
			logger.Error("failed to parse keepalive", logger.ErrorAttr(err))
			return
		}

		b, err := bench.NewBenchmark(
			Cfg,
			bench.WithClientID(clientID),
			bench.WithClients(clients),
			bench.WithTopic(topic),
			bench.WithQoS(qos),
			bench.WithMessageCount(count),
			bench.WithDelay(delay),
			bench.WithCleanSession(cleanSession),
			bench.WithKeepAlive(keepalive),
			bench.WithUsername(username),
			bench.WithPassword(password),
			bench.WithHost(host),
			bench.WithPort(port),
		)
		if err != nil {
			logger.Error("failed to create benchmark", logger.State("failed"), logger.ErrorAttr(err))
			return
		}

		go func() {
			<-sigs
			logger.Info("received shutdown signal", logger.State("interrupted"))
			os.Exit(0)
		}()

		b.Subscribe()
	},
}

func init() {
	rootCmd.AddCommand(subCmd)

	// Register flags
	subCmd.Flags().IntP("clients", "c", 100, "Number of concurrent subscriber clients")
	subCmd.Flags().IntP("delay", "d", 1000, "Delay between subscription lifetime checks (ms)")
	subCmd.Flags().IntP("count", "n", 1000, "Expected number of messages per client")
	subCmd.Flags().Uint16P("qos", "q", 0, "Quality of service level (0, 1, 2)")
	subCmd.Flags().StringP("topic", "t", "bench/test", "Topic to subscribe to")
}
