# SeederSync

[![Go Report Card](https://goreportcard.com/badge/github.com/tuzzmaniandevil/seedersync)](https://goreportcard.com/report/github.com/tuzzmaniandevil/seedersync)
[![GoDoc](https://godoc.org/github.com/tuzzmaniandevil/seedersync?status.svg)](https://godoc.org/github.com/tuzzmaniandevil/seedersync)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

`SeederSync` is a Go application designed to synchronize and manage torrent trackers for a qBittorrent client. It automates the process of fetching up-to-date tracker lists from various sources, updating qBittorrent's default trackers, applying these trackers to existing torrents, and optionally reannouncing torrents to ensure optimal peer connectivity.

## Features

*   **Tracker Management:** Fetches tracker lists from configurable remote URLs and local static lists.
*   **qBittorrent Integration:** Connects to a qBittorrent client to:
    *   Update the client's default list of trackers.
    *   Add fetched trackers to existing torrents.
    *   Optionally reannounce all torrents after tracker updates.
*   **Flexible Configuration:** Configure `SeederSync` via:
    *   Command-line flags
    *   Environment variables
    *   Configuration files (YAML, TOML, or JSON format - `config.yaml`, `config.yml`, `config.toml`, `config.json`)
*   **Tracker Contribution:** Optionally contributes newly found trackers to `newtrackon.com` to help the community.
*   **Deduplication:** Ensures that the tracker list applied is free of duplicates.

## Configuration

`SeederSync` can be configured using command-line flags, environment variables, or a configuration file. The order of precedence is: CLI flags > Environment variables > Configuration file.

**Configuration File:**

Create a configuration file named `config.yaml`, `config.yml`, `config.toml`, or `config.json` in the same directory as the `SeederSync` executable, or specify a path using the `--config` flag or `QBIT_CONFIG_FILE` environment variable.

**Example `config.yaml`:**

```yaml
host: "http://localhost:8080"  # qBittorrent host
username: "your_username"        # qBittorrent username
password: "your_password"        # qBittorrent password
tlsSkipVerify: false             # Skip TLS certificate verification (default: false)
contributeTrackers: true         # Contribute trackers to newtrackon.com (default: false)
reannounceTorrents: true         # Reannounce torrents after updating (default: false)
staticTrackers:                  # Optional: List of static trackers to always include
  - "udp://tracker.opentrackr.org:1337/announce"
trackerListURLs:                 # Optional: List of URLs to fetch tracker lists from
  - "https://newtrackon.com/api/stable?include_ipv4_only_trackers=true&include_ipv6_only_trackers=false"
  - "https://trackerslist.com/best.txt"
reannounceMaxAttempts: 3         # Max attempts for reannouncing a torrent (default: 3)
reannounceInterval: 5            # Interval in seconds between reannounce attempts (default: 5)
```

**Command-line Flags & Environment Variables:**

| Flag                    | Environment Variable        | Description                                                                 | Default (if not in config) |
| ----------------------- | --------------------------- | --------------------------------------------------------------------------- | -------------------------- |
| `--host`                | `QBIT_HOST`                 | qBittorrent host (e.g., http://localhost:8080)                              | **Required**               |
| `--username`            | `QBIT_USERNAME`             | qBittorrent username                                                        | ""                         |
| `--password`            | `QBIT_PASSWORD`             | qBittorrent password                                                        | ""                         |
| `--tlsSkipVerify`       | `QBIT_TLS_SKIP_VERIFY`      | Skip TLS certificate verification for qBittorrent client                    | `false`                    |
| `--contributeTrackers`  | `QBIT_CONTRIBUTE_TRACKERS`  | Contribute found trackers to newtrackon.com                                 | `false`                    |
| `--reannounceTorrents`  | `QBIT_REANNOUNCE_TORRENTS`  | Reannounce all torrents after updating trackers                             | `false`                    |
| `--config`              | `QBIT_CONFIG_FILE`          | Path to config file (YAML, TOML, or JSON)                                   | (tries default files)      |

*Note: The `host` is a required parameter and must be provided either via a flag, environment variable, or in the configuration file.*

## Usage

1.  **Prepare Configuration:**
    *   Ensure your qBittorrent client is running and accessible.
    *   Create a configuration file (e.g., `config.yaml`) with your qBittorrent details and preferences, or prepare to use environment variables/CLI flags.

2.  **Run `seedersync`:**
    *   If you have a pre-compiled binary:
        ```bash
        ./seedersync --host "http://your-qbittorrent-host:port" --username "user" --password "pass"
        # Or, if using a config file:
        ./seedersync --config /path/to/your/config.yaml
        # Or, if config.yaml is in the same directory:
        ./seedersync
        ```
    *   If running from source:
        ```bash
        go run main.go --host "http://your-qbittorrent-host:port" --username "user" --password "pass"
        # Or with a config file:
        go run main.go --config /path/to/your/config.yaml
        ```

The application will log its progress to the console.

## Building from Source

To build `SeederSync` from source, you need to have Go installed (version 1.24.3 or later is recommended, as per `go.mod`).

1.  **Clone the repository (if you haven't already):**
    ```bash
    git clone https://github.com/tuzzmaniandevil/seedersync.git # Replace with your actual repo URL
    cd seedersync
    ```

2.  **Build the application:**
    ```bash
    go build -o seedersync main.go
    ```
    This will create an executable file named `seedersync` in the current directory.

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request. For major changes, please open an issue first to discuss what you would like to change.

### Guidelines

- Write clear, concise commit messages.
- Ensure all tests pass before submitting a pull request.
- Follow the existing code style and format your code with `gofmt`.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.