# envoy-diff

> Compares Envoy xDS config snapshots across environments and outputs structured diffs for auditing.

---

## Installation

```bash
go install github.com/yourorg/envoy-diff@latest
```

Or build from source:

```bash
git clone https://github.com/yourorg/envoy-diff.git && cd envoy-diff && go build -o envoy-diff .
```

---

## Usage

Provide two xDS config snapshot files (JSON or YAML) and envoy-diff will output a structured diff suitable for auditing or CI pipelines.

```bash
envoy-diff --base snapshots/staging.json --target snapshots/production.json
```

**Example output:**

```
[~] clusters.payment-service.connect_timeout: 5s → 10s
[+] listeners.grpc-ingress.filter_chains[2]: <added>
[-] routes.api-gateway.virtual_hosts[1].retry_policy: <removed>
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--base` | Path to the base snapshot file | _(required)_ |
| `--target` | Path to the target snapshot file | _(required)_ |
| `--format` | Output format: `text`, `json`, `yaml` | `text` |
| `--ignore` | Comma-separated field paths to ignore | `""` |
| `--exit-code` | Exit with code 1 if differences are found | `false` |

```bash
# Output diff as JSON and fail if differences exist
envoy-diff --base base.json --target target.json --format json --exit-code
```

---

## License

MIT © yourorg