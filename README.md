# procwatch

Lightweight process monitor that sends alerts via webhook when a service crashes or exceeds resource thresholds.

## Installation

```bash
go install github.com/yourusername/procwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/procwatch.git && cd procwatch && go build -o procwatch .
```

## Usage

Create a `config.yaml` file:

```yaml
webhook_url: "https://hooks.slack.com/services/your/webhook/url"
interval: 10s
processes:
  - name: "nginx"
    pid_file: "/var/run/nginx.pid"
    max_cpu: 80.0
    max_memory_mb: 512
  - name: "myapp"
    pattern: "myapp --serve"
    max_cpu: 60.0
    max_memory_mb: 256
```

Run procwatch:

```bash
procwatch --config config.yaml
```

procwatch will poll the specified processes at the configured interval and POST a JSON alert payload to your webhook URL if a process is not found or exceeds the defined resource thresholds.

### Alert Payload Example

```json
{
  "process": "nginx",
  "event": "crash",
  "message": "Process 'nginx' not found",
  "timestamp": "2024-05-10T14:32:00Z"
}
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--config` | `config.yaml` | Path to configuration file |
| `--log-level` | `info` | Log verbosity (debug, info, warn, error) |

## License

MIT © 2024 yourusername