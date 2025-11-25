package bench

import (
	"sync"

	"github.com/rayomqio/benchmq/pkg/config"
	"github.com/rayomqio/benchmq/pkg/er"
	"github.com/rayomqio/benchmq/pkg/logger"
)

type QoSLevel uint8

// Bench represents the benchmark fields
type Bench struct {
	delay        int
	clients      int
	clientID     string
	topic        string
	message      string
	messageCount int
	retained     bool
	cleanSession *bool
	qos          QoSLevel
	keepAlive    uint16
	host         string
	port         uint16
	username     string
	password     string
	wg           sync.WaitGroup // Wait Group
	cfg          *config.Config // Config
	logger       *logger.Logger // Logger
}

type Option func(*Bench)

const (
	QoS0 QoSLevel = 0 // QoS At Most Once
	QoS1 QoSLevel = 1 // QoS At Least Once
	QoS2 QoSLevel = 2 // QoS Exactly Once
)

const (
	DefaultDelay        = 1000             // Default delay between connection (ms)
	DefaultClients      = 100              // Default clients to connect
	DefaultClientID     = "benchmq-client" // Default client id
	DefaultTopic        = "bench/test"     // Default publish/subscribe topic
	DefaultCleanSession = true             // Default clean session state
	DefaultQoS          = QoS0             // Default QoS level
	DefaultKeepAlive    = 60               // Default connection keep alive
	DefaultMessageCount = 100              // Default message count
	DefaultMessage      = "Hello, World!"  // Default message
	DefaultRetained     = false            // Default retained message state
)

// NewBenchmark constructor initializes the bench struct
func NewBenchmark(cfg *config.Config, options ...Option) (*Bench, error) {
	if cfg == nil {
		return nil, &er.Error{
			Package: "Bench",
			Func:    "NewBenchmark",
			Message: er.ErrNilConfig,
			Raw:     er.ErrNilConfig,
		}
	}

	bench := Bench{
		delay:        DefaultDelay,
		clients:      DefaultClients,
		clientID:     DefaultClientID,
		topic:        DefaultTopic,
		message:      DefaultMessage,
		messageCount: DefaultMessageCount,
		retained:     DefaultRetained,
		cleanSession: &cfg.Client.CleanSession,
		qos:          DefaultQoS,
		keepAlive:    cfg.Client.KeepAlive,
		host:         cfg.Server.Host,
		port:         cfg.Server.Port,
		cfg:          cfg,
		logger:       logger.NewBenchmarkLogger("bench"),
	}

	for _, option := range options {
		if option != nil {
			option(&bench)
		}
	}

	if err := bench.validate(); err != nil {
		return nil, err
	}

	return &bench, nil
}

// Validate checks semantic correctness of the benchmark configuration
func (b *Bench) validate() error {
	if b.clients <= 0 {
		return &er.Error{
			Package: "Bench",
			Func:    "Validate",
			Message: er.ErrInvalidClients,
			Raw:     er.ErrInvalidClients,
		}
	}
	if b.delay < 0 {
		return &er.Error{
			Package: "Bench",
			Func:    "Validate",
			Message: er.ErrInvalidDelay,
			Raw:     er.ErrInvalidDelay,
		}
	}
	if b.host == "" {
		return &er.Error{
			Package: "Bench",
			Func:    "Validate",
			Message: er.ErrEmptyHost,
			Raw:     er.ErrEmptyHost,
		}
	}
	if b.topic == "" {
		return &er.Error{
			Package: "Bench",
			Func:    "Validate",
			Message: er.ErrEmptyTopic,
			Raw:     er.ErrEmptyTopic,
		}
	}
	if b.port == 0 {
		return &er.Error{
			Package: "Bench",
			Func:    "Validate",
			Message: er.ErrInvalidPort,
			Raw:     er.ErrInvalidPort,
		}
	}
	if b.qos > QoS2 {
		return &er.Error{
			Package: "Bench",
			Func:    "Validate",
			Message: er.ErrInvalidQoS,
			Raw:     er.ErrInvalidQoS,
		}
	}
	// Set default clientID
	if b.clientID == "" {
		b.clientID = DefaultClientID
	}
	if b.keepAlive == 0 {
		b.keepAlive = DefaultKeepAlive
	}
	if b.cleanSession == nil {
		cs := true
		b.cleanSession = &cs
	}
	return nil
}

func WithDelay(delay int) Option {
	return func(b *Bench) {
		b.delay = delay
	}
}

func WithClients(clients int) Option {
	return func(b *Bench) {
		b.clients = clients
	}
}

func WithClientID(clientID string) Option {
	return func(b *Bench) {
		b.clientID = clientID
	}
}

func WithTopic(topic string) Option {
	return func(b *Bench) {
		b.topic = topic
	}
}

func WithCleanSession(cleanSession bool) Option {
	return func(b *Bench) {
		b.cleanSession = &cleanSession
		if b.cfg != nil {
			b.cfg.Client.CleanSession = cleanSession
		}
	}
}

func WithQoS(qos uint16) Option {
	return func(b *Bench) {
		b.qos = QoSLevel(qos)
	}
}

func WithKeepAlive(keepAlive uint16) Option {
	return func(b *Bench) {
		b.keepAlive = keepAlive
		if b.cfg != nil {
			b.cfg.Client.KeepAlive = keepAlive
		}
	}
}

func WithHost(host string) Option {
	return func(b *Bench) {
		b.host = host
		if b.cfg != nil {
			b.cfg.Server.Host = host
		}
	}
}

func WithPort(port uint16) Option {
	return func(b *Bench) {
		b.port = port
		if b.cfg != nil {
			b.cfg.Server.Port = port
		}
	}
}

func WithMessage(message string) Option {
	return func(b *Bench) {
		b.message = message
	}
}

func WithMessageCount(count int) Option {
	return func(b *Bench) {
		b.messageCount = count
	}
}

func WithRetained(retained bool) Option {
	return func(b *Bench) {
		b.retained = retained
	}
}

func WithUsername(username string) Option {
	return func(b *Bench) {
		b.username = username
	}
}

func WithPassword(password string) Option {
	return func(b *Bench) {
		b.password = password
	}
}
