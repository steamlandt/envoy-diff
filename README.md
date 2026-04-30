# envoy-diff

> CLI tool to diff environment variable sets across `.env` files and running processes

---

## Installation

```bash
go install github.com/yourusername/envoy-diff@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envoy-diff.git
cd envoy-diff
go build -o envoy-diff .
```

---

## Usage

Compare two `.env` files:

```bash
envoy-diff .env.staging .env.production
```

Compare a `.env` file against a running process:

```bash
envoy-diff .env --pid 12345
```

Compare the environment of two running processes:

```bash
envoy-diff --pid 12345 --pid 67890
```

### Example Output

```
+ DB_HOST=prod.db.example.com
- DB_HOST=staging.db.example.com
  DB_PORT=5432
+ NEW_FEATURE_FLAG=true
- STAGING_ONLY_VAR=foo
```

Keys present in one source but missing in the other are flagged, and differing values are shown inline.

---

## Flags

| Flag | Description |
|------|-------------|
| `--pid` | Target a running process by PID |
| `--ignore` | Comma-separated list of keys to ignore |
| `--only-missing` | Show only keys missing from one side |

---

## License

MIT © 2024 yourusername