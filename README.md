# mtls-bridge

Wrap any stdin/stdout process with mutual TLS. Designed for bridging [MCP](https://modelcontextprotocol.io/) servers across machines, but works with any stdio-based protocol.

## Install

Download the latest release from [GitHub Releases](https://github.com/jdfalk/mtls-bridge/releases), or build from source:

```bash
go install github.com/jdfalk/mtls-bridge/cmd/mtls-bridge@latest
```

## Quick Start

### 1. Generate a pre-shared key

```bash
mtls-bridge provision --generate-psk
```

Copy `.mtls/psk.txt` to the other machine (or use a shared filesystem).

### 2. Start the server

```bash
mtls-bridge serve --powershell "/path/to/your-script.ps1"
```

### 3. Connect from the client

```bash
mtls-bridge connect
```

The first connection exchanges the PSK for mTLS certificates. All subsequent connections use mutual TLS.

## How It Works

```
Client (stdio) <-> mtls-bridge connect <-mTLS/TCP-> mtls-bridge serve <-> Subprocess (stdio)
```

- **Provisioning:** A pre-shared key (PSK) bootstraps certificate generation. The server generates a CA + server cert + client cert, sends the client its credentials over a TLS-encrypted channel.
- **Normal operation:** Both sides present certificates signed by the shared CA. TLS 1.3 minimum.
- **Auto-update:** The `serve` command auto-updates on startup. The `connect` command notifies when updates are available.
- **Reconnect:** The `connect` command automatically reconnects with exponential backoff if the connection drops.

## Commands

| Command | Description |
|---------|-------------|
| `serve --powershell <path>` | Start mTLS server wrapping a subprocess |
| `connect` | Connect to server, bridge to local stdio |
| `provision --generate-psk` | Generate a new pre-shared key |
| `provision --renew` | Regenerate certs from existing CA |
| `provision --reset` | Delete all certs and config |
| `update` | Self-update to latest release |
| `version` | Print version, commit, and build date |

## License

MIT
