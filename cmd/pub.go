package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rayomqio/benchmq/internal/bench"
	"github.com/rayomqio/benchmq/pkg/logger"
	"github.com/spf13/cobra"
)

// pubCmd represents the pub command
var pubCmd = &cobra.Command{
	Use:   "pub",
	Short: "Publish messages to a topic with specified parameters",
	Long: `Publish messages to a topic with specified parameters.

Parameters:
	- host: Hostname or IP address of the broker
	- port: Port number of the broker
	- clientID: Base client ID prefix (each client appends "-<n>")
    - clients: Number of concurrent clients
    - delay: Delay between messages in milliseconds
    - count: Number of messages to publish per client
    - qos: Quality of service level (0, 1, 2)
    - message: The message payload
    - topic: Topic to publish to
    - retain: Whether to retain the last message
    - clean: Whether to use a clean session
    - keepalive: Keepalive interval in seconds`,
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

		retain, err := cmd.Flags().GetBool("retain")
		if err != nil {
			logger.Error("failed to parse retain flag", logger.ErrorAttr(err))
			return
		}

		message, err := cmd.Flags().GetString("message")
		if err != nil {
			logger.Error("failed to parse message", logger.ErrorAttr(err))
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

		keepalive, err := cmd.Flags().GetUint16("keepalive")
		if err != nil {
			logger.Error("failed to parse keepalive", logger.ErrorAttr(err))
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

		b, err := bench.NewBenchmark(
			Cfg,
			bench.WithClientID(clientID),
			bench.WithClients(clients),
			bench.WithTopic(topic),
			bench.WithQoS(qos),
			bench.WithMessageCount(count),
			bench.WithDelay(delay),
			bench.WithRetained(retain),
			bench.WithCleanSession(cleanSession),
			bench.WithKeepAlive(keepalive),
			bench.WithMessage(message),
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

		b.PublishMessages()
	},
}

func init() {
	rootCmd.AddCommand(pubCmd)

	// Register flags
	pubCmd.Flags().IntP("clients", "c", 100, "Number of concurrent clients to connect")
	pubCmd.Flags().IntP("delay", "d", 1000, "Delay between messages in milliseconds")
	pubCmd.Flags().IntP("count", "n", 1000, "Number of messages to publish per client")
	pubCmd.Flags().BoolP("retain", "r", false, "Retain the last message")
	pubCmd.Flags().Uint16P("qos", "q", 0, "Quality of service level (0, 1, 2)")
	pubCmd.Flags().StringP("message", "m", "Hello, World!", "Message to publish")
	pubCmd.Flags().StringP("topic", "t", "bench/test", "Topic to publish messages to")
}
