# BenchMQ

**BenchMQ** is a simple, fast, and open-source CLI tool for benchmarking MQTT brokers. Measure throughput, latency, and stability of your MQTT setup with ease.

## Features

- 🚀 **Zero Dependencies**: Single binary with no external config file required
- 📊 **Multiple Benchmark Types**: Connection, publish, and subscribe benchmarks
- 🔧 **Flexible Configuration**: Use command-line flags or optional config file
- 📈 **Concurrent Testing**: Support for multiple concurrent clients
- 🎯 **Quality of Service**: Full QoS 0, 1, and 2 support
- 🔐 **Authentication**: Username/password authentication support
- 📝 **Detailed Logging**: Comprehensive logging with different levels

## Quick Start

### Download

#### Using Homebrew
```bash
# Add the tap
brew tap rayomqio/tap

# Install benchmq
brew install benchmq
```

#### Using Go Install
```bash
go install github.com/rayomqio/benchmq@latest
```

#### From Releases
Download the latest binary from the [releases page](https://github.com/rayomqio/benchmq/releases) for your platform.

#### Build from Source
```bash
git clone https://github.com/rayomqio/benchmq.git
cd benchmq
go build -o benchmq .
```

#### Using Docker
```bash
# Pull from GitHub Container Registry (GHCR)
docker pull ghcr.io/rayomqio/benchmq:latest

# Pull from Docker Hub (when available)
docker pull me3di/benchmq:latest

# Or build locally
docker build -t benchmq .

# Run with Docker
docker run --rm ghcr.io/rayomqio/benchmq:latest --help
```

### Basic Usage

BenchMQ works out of the box with sensible defaults. No configuration file needed!

**Prerequisites**: You need an MQTT broker running. For testing, you can quickly set up:

**Option 1: Mosquitto (Local)**
```bash
# macOS
brew install mosquitto
mosquitto

# Ubuntu/Debian
sudo apt install mosquitto
mosquitto

# Windows
# Download from https://mosquitto.org/download/
```

**Option 2: Docker (Recommended for testing)**
```bash
# Quick start with default config
docker run -it -p 1883:1883 eclipse-mosquitto

# With persistence
docker run -it -p 1883:1883 -v mosquitto-data:/mosquitto/data eclipse-mosquitto
```

**Option 3: Cloud Services**
- [HiveMQ Cloud](https://www.hivemq.com/mqtt-cloud-broker/) (Free tier available)
- [AWS IoT Core](https://aws.amazon.com/iot-core/)
- [Azure IoT Hub](https://azure.microsoft.com/en-us/services/iot-hub/)

```bash
# Test MQTT broker connections (localhost:1883)
benchmq conn

# Publish messages to a topic
benchmq pub -t test/topic -m "Hello World" -c 10

# Subscribe to a topic
benchmq sub -t test/topic -c 5
```

### Docker Usage

When using Docker, you can run BenchMQ commands by passing them as arguments:

```bash
# Test connections using Docker
docker run --rm ghcr.io/rayomqio/benchmq:latest conn -c 50

# Publish messages using Docker
docker run --rm ghcr.io/rayomqio/benchmq:latest pub -t test/topic -m "Hello Docker" -c 10

# Subscribe to messages using Docker
docker run --rm ghcr.io/rayomqio/benchmq:latest sub -t test/topic -c 5

# Connect to external MQTT broker (replace with your broker's IP/hostname)
docker run --rm ghcr.io/rayomqio/benchmq:latest conn --host broker.hivemq.com --port 1883 -c 10
```

**Note**: When running in Docker, `localhost` refers to the container itself. To connect to an MQTT broker on your host machine or external services, you'll need to:
- Use `host.docker.internal` (macOS/Windows) to connect to host machine
- Use the actual IP address or hostname of your MQTT broker
- Use Docker networks for container-to-container communication

## Commands

### Connection Benchmark (`conn`)

Test connection throughput and stability by opening multiple concurrent MQTT connections.

```bash
benchmq conn [flags]
```

**Examples:**
```bash
# Test 100 concurrent connections with 1 second delay between connections
benchmq conn -c 100 -d 1000

# Test connections with authentication
benchmq conn -u myuser -p mypass -c 50

# Test with custom client ID prefix
benchmq conn -i "load-test" -c 200
```

**Flags:**
- `-H, --host string`: Hostname or IP address of the broker
- `-P, --port uint16`: Port number of the broker (default: 1883)
- `-c, --clients int`: Number of concurrent clients (default: 100)
- `-d, --delay int`: Delay between connections in milliseconds (default: 1000)
- `-i, --clientID string`: Client ID prefix (default: "benchmq-client")
- `-u, --username string`: MQTT username
- `-p, --password string`: MQTT password
- `-k, --keepalive uint16`: Keepalive interval in seconds (default: 60)
- `-x, --clean`: Clean session flag (default: true)

### Publish Benchmark (`pub`)

Benchmark message publishing with multiple concurrent publishers.

```bash
benchmq pub [flags]
```

**Examples:**
```bash
# Publish 1000 messages per client with 10 concurrent publishers
benchmq pub -t sensors/data -c 10 -n 1000 -m '{"temp":23.5}'

# High-frequency publishing (no delay between messages)
benchmq pub -t test/performance -d 0 -c 5 -n 5000

# QoS 2 publishing with message retention
benchmq pub -t important/data -q 2 -r -n 100
```

**Flags:**
- `-H, --host string`: Hostname or IP address of the broker
- `-P, --port uint16`: Port number of the broker (default: 1883)
- `-t, --topic string`: Topic to publish to (default: "benchmq")
- `-m, --message string`: Message payload (default: "Hello, World!")
- `-c, --clients int`: Number of concurrent publishers (default: 100)
- `-n, --count int`: Messages per client (default: 1000)
- `-d, --delay int`: Delay between messages in milliseconds (default: 1000)
- `-q, --qos uint16`: Quality of service (0, 1, or 2) (default: 0)
- `-r, --retain`: Retain messages
- `-i, --clientID string`: Client ID prefix (default: "benchmq-client")
- `-u, --username string`: MQTT username
- `-p, --password string`: MQTT password
- `-k, --keepalive uint16`: Keepalive interval in seconds (default: 60)
- `-x, --clean`: Clean session flag (default: true)

### Subscribe Benchmark (`sub`)

Benchmark message subscription with multiple concurrent subscribers.

```bash
benchmq sub [flags]
```

**Examples:**
```bash
# Subscribe with 5 clients expecting 1000 messages each
benchmq sub -t sensors/# -c 5 -n 1000

# QoS 1 subscription with authentication
benchmq sub -t data/stream -q 1 -u subscriber -p secret

# Long-running subscription test
benchmq sub -t test/topic -c 10 -n 10000 -d 5000
```

**Flags:**
- `-H, --host string`: Hostname or IP address of the broker (default: "localhost")
- `-P, --port uint16`: Port number of the broker (default: 1883)
- `-t, --topic string`: Topic to subscribe to (default: "benchmq")
- `-c, --clients int`: Number of concurrent subscribers (default: 100)
- `-n, --count int`: Expected messages per client (default: 1000)
- `-d, --delay int`: Delay between checks in milliseconds (default: 1000)
- `-q, --qos uint16`: Quality of service (0, 1, or 2) (default: 0)
- `-i, --clientID string`: Client ID prefix (default: "benchmq-subscriber")
- `-u, --username string`: MQTT username
- `-p, --password string`: MQTT password
- `-k, --keepalive uint16`: Keepalive interval in seconds (default: 60)
- `-x, --clean`: Clean session flag (default: true)

## Configuration

### Command Line Only (Recommended)

BenchMQ works perfectly with just command-line flags. All broker connection details can be specified via flags:

```bash
# Connect to remote broker with authentication
benchmq conn \
  --clients 50 \
  --username myuser \
  --password mypass
```

**Note**: By default, BenchMQ connects to `localhost:1883`. To connect to a different broker, you'll need to use a configuration file (see below).

### Optional Configuration File

For advanced scenarios or to avoid repeating flags, you can create an optional `config.yml` file:

```yaml
name: BenchMQ
version: 1.0.0
environment: development  # or production

server:
  host: mqtt.example.com  # Change this for remote brokers
  port: 1883              # Standard MQTT port (8883 for TLS)

client:
  client_id: benchmq-client
  keep_alive: 60
  clean_session: true
  username: ""            # Set if broker requires auth
  password: ""            # Set if broker requires auth
```

Place this file in the same directory as the binary. If no config file exists, BenchMQ will use sensible defaults.

## Common Use Cases

### Testing Broker Capacity

```bash
# Test maximum concurrent connections
benchmq conn -c 1000 -d 100

# Test sustained message throughput
benchmq pub -c 20 -n 10000 -d 0 -t load/test
```

### Load Testing Before Production

```bash
# Simulate IoT device connections
benchmq conn -c 500 -i "device" -d 200

# Test sensor data publishing
benchmq pub -t sensors/temperature -c 100 -n 1440 -d 60000 -m '{"temp":22.5,"unit":"C"}'
```

### Performance Benchmarking

```bash
# High-frequency publishing test
benchmq pub -c 10 -n 10000 -d 0 -q 0

# QoS comparison test
benchmq pub -q 2 -c 5 -n 1000  # Test QoS 2 performance
```

### Authentication Testing

```bash
# Test authenticated connections
benchmq conn -u testuser -p testpass -c 100

# Test with different client IDs
benchmq pub -i "auth-test" -u user -p pass -t secure/data
```

## Output and Monitoring

BenchMQ provides detailed logging output including:
- Connection success/failure rates
- Message publishing statistics
- Timing information
- Error details
- Progress indicators

Logs are written to stdout and include timestamps, log levels, and structured information for easy parsing.

## Troubleshooting

### Connection Refused Errors
If you see "Couldn't establish client" errors, ensure:
1. MQTT broker is running and accessible
2. Correct host/port in config (default: localhost:1883)
3. Firewall allows MQTT traffic
4. Authentication credentials are correct (if required)

### Performance Issues
- Start with fewer clients and increase gradually
- Monitor system resources (CPU, memory, network)
- Check broker logs for errors or limits
- Consider broker connection limits and message rate limits

## Best Practices

1. **Start Small**: Begin with low client counts and increase gradually
2. **Monitor Resources**: Watch CPU and memory usage during large tests
3. **Use Realistic Payloads**: Test with message sizes similar to your use case
4. **Test Different QoS Levels**: Each QoS level has different performance characteristics
5. **Consider Network Latency**: Factor in network conditions when interpreting results
6. **Have a Broker Ready**: Ensure your MQTT broker is properly configured and running

## License

This project is licensed under the Apache 2.0 License.

## Support

- 🐛 [Report Issues](https://github.com/rayomqio/benchmq/issues)
