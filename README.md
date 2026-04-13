# portwatch

A lightweight CLI daemon that monitors open ports and alerts on unexpected changes in real time.

---

## Installation

```bash
go install github.com/yourusername/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/portwatch.git
cd portwatch && go build -o portwatch .
```

---

## Usage

Start the daemon with default settings (polls every 30 seconds):

```bash
portwatch start
```

Specify a custom polling interval and log output:

```bash
portwatch start --interval 10s --log /var/log/portwatch.log
```

Run a one-time snapshot of currently open ports:

```bash
portwatch scan
```

**Example alert output:**

```
[2024-11-15 14:32:01] ALERT: New port opened → TCP :8080
[2024-11-15 14:35:44] ALERT: Port closed    → TCP :3306
```

### Flags

| Flag         | Default | Description                        |
|--------------|---------|------------------------------------|
| `--interval` | `30s`   | Polling interval                   |
| `--log`      | stdout  | Path to log file                   |
| `--baseline` | —       | File to load a known-good snapshot |

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

---

## License

[MIT](LICENSE) © 2024 yourusername