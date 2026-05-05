# envsync

> Diff and sync `.env` files across environments with secret masking.

---

## Installation

```bash
go install github.com/yourusername/envsync@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envsync.git
cd envsync && go build -o envsync .
```

---

## Usage

**Diff two `.env` files:**

```bash
envsync diff .env.development .env.production
```

**Sync missing keys from one file to another:**

```bash
envsync sync .env.development .env.production
```

**Mask secrets in output:**

```bash
envsync diff .env.development .env.production --mask-secrets
```

Example output:

```
~ API_URL        dev: http://localhost:3000  prod: https://api.example.com
+ NEW_FEATURE_FLAG  dev: true  prod: [missing]
- LEGACY_KEY     dev: [missing]  prod: ********
```

---

## Options

| Flag             | Description                          |
|------------------|--------------------------------------|
| `--mask-secrets` | Redact values matching secret patterns |
| `--output json`  | Output diff results as JSON           |
| `--dry-run`      | Preview sync changes without writing  |

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

---

## License

[MIT](LICENSE) © 2024 yourusername