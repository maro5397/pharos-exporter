# pharos-exporter
pharos exporter for validators

## Usage

Run the `start` command with your RPC endpoint, BLS key, and log path:

```bash
go run . start \
  -rpc https://YOUR_RPC \
  -my-bls-key 0xYOUR_BLS_KEY \
  -log-path /data/pharos-node/domain/light/log/consensus.log
```

### Notes
- `-log-from-start` reads the log from the beginning; omit it to tail only new lines.
- The log tailer follows file rotation (e.g. `consensus.log` renamed to `consensus.log.x`).
- `-check-block-proof`, `-check-validator-set`, `-check-propose` and `-check-endorse` are enabled by default.

### Options
Use `-h` to see all available flags and defaults:

```bash
go run . start -h
```

Example output:
```text
Usage of start:
  -check-block-proof
        check signedBlsKeys metrics (default true)
  -check-endorse
        check endorse metrics (default true)
  -check-propose
        check propose metrics (default true)
  -check-validator-set
        check validator set metrics (default true)
  -exporter-port string
        metrics listen port (default "9123")
  -log-from-start
        start reading log from beginning (default: false)
  -log-path string
        path to log file to tail
  -log-poll-interval duration
        poll interval for log tailing (default 1s)
  -my-bls-key string
        my BLS pubkey (0x...)
  -rpc string
        JSON-RPC endpoint (default "https://atlantic-rpc.dplabs-internal.com/")
  -rpc-poll-interval duration
        poll interval for latest block (default 1s)
```

### Exported Metrics
The `/metrics` endpoint includes default Go/process/promhttp metrics. Custom metrics exposed by this exporter:

- `validator_active_timestamp` (gauge): Unix timestamp when validator active status was last observed.
- `validator_active_total` (counter): Total number of blocks where the validator was active in the validator set.
- `validator_endorse_total` (counter): Total number of endorse events observed in logs.
- `validator_last_endorse_timestamp` (gauge): Unix timestamp of the last endorse event observed in logs.
- `validator_last_propose_timestamp` (gauge): Unix timestamp of the last propose event observed in logs.
- `validator_propose_total` (counter): Total number of propose attempts observed in logs.
- `validator_vote_inclusion_timestamp` (gauge): Unix timestamp when the validator vote was last included.
- `validator_vote_inclusion_total` (counter): Total number of blocks where the validator vote was included.

## Systemd Setup

Build the binary and install it to `/usr/local/bin`:

```bash
go build -o pharos-exporter .
sudo mv pharos-exporter /usr/local/bin/
```

Copy the example service unit and update the placeholders:

```bash
sudo cp pharos-exporter.service.example /etc/systemd/system/pharos-exporter.service
sudo sed -i 's|<https://YOUR_RPC>|https://YOUR_RPC|g' /etc/systemd/system/pharos-exporter.service
sudo sed -i 's|<0xYOUR_BLS_KEY>|0xYOUR_BLS_KEY|g' /etc/systemd/system/pharos-exporter.service
sudo sed -i 's|<YOUR_LOG_PATH>|YOUR_LOG_PATH|g' /etc/systemd/system/pharos-exporter.service
```

Reload systemd and start the service:

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now pharos-exporter
```
