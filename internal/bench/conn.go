package bench

import (
	"fmt"
	"time"

	"github.com/rayomqio/benchmq/internal/mqtt"
	"github.com/rayomqio/benchmq/pkg/logger"
)

func (b *Bench) RunConnections() {
	start := time.Now()
	b.logger.Info("started connection benchmark", logger.Int("time", int(start.UnixNano())))

	for i := 0; i < b.clients; i++ {
		b.wg.Add(1)
		go func(id int) {
			defer b.wg.Done()

			cfg := *b.cfg
			cfg.Client.ClientID = fmt.Sprintf("%s-%d", b.clientID, id)
			cfg.Client.CleanSession = *b.cleanSession
			cfg.Client.KeepAlive = b.keepAlive
			cfg.Client.Username = b.username
			cfg.Client.Password = b.password
			client := mqtt.NewClient(&cfg)
			defer client.Disconnect()

			if err := client.Connect(); err != nil {
				b.logger.Error("couldn't establish client", logger.ClientID(cfg.Client.ClientID), logger.State("failed"))
				return
			}
			b.logger.LogClientConnection(cfg.Client.ClientID)
		}(i)
		time.Sleep(time.Duration(b.delay) * time.Millisecond)
	}

	b.wg.Wait()
	b.logger.Info("finished connection benchmark", logger.Any("time", time.Since(start).Seconds()))
}
