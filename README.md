# Space-Dos

Space-DoS is a simple HTTP stress-testing tool written in Go. It allows users to send concurrent HTTP requests to a target URL for a specified duration, simulating a denial-of-service (DoS) attack. **This tool is for educational and testing purposes only. Do not use it for malicious activities.**

## Features

- Supports both `GET` and `POST` HTTP methods.
- Configurable number of threads and attack duration.
- Randomized user agents, referrers, and paths for each request.
- Auto-retry functionality if the target server comes back online.
- Logs all sent requests and errors to a file (`attack_log.txt`).

## Usage

### Prerequisites

- Go 1.16 or later installed on your system.

### Build

To build the project, run:

```bash
go build -o space
```

### Run

To run the tool, use the following command:

```bash
./space -url <target> -threads <number_of_threads> -duration <duration_in_seconds> -method <HTTP_method> [-data <POST_payload>] [-retry <true|false>]
```

### Command-Line Arguments

| Argument       | Description                                      | Default Value            |
|----------------|--------------------------------------------------|--------------------------|
| `-url`         | Target URL (e.g., `http://localhost:8080`)       | **Required**             |
| `-threads`     | Number of concurrent threads                     | `10`                     |
| `-duration`    | Duration of the attack in seconds                | `10`                     |
| `-method`      | HTTP method: `GET` or `POST`                     | `GET`                    |
| `-data`        | POST data payload (only used with `POST` method) | `id=test&value=123`      |
| `-retry`       | Auto retry if the server comes back online       | `true`                   |

### Example

```bash
./space -url http://example.com -threads 20 -duration 30 -method POST -data "key=value" -retry true
```

This command will send `POST` requests with the payload `key=value` to `http://example.com` using 20 threads for 30 seconds. If the server goes offline, it will retry automatically.

## Output

- Logs are saved to `attack_log.txt` in the same directory as the executable.
- Example log entries:
  ```
  2023-03-01T12:00:00Z [+] Sent: GET /index
  2023-03-01T12:00:01Z [-] Error: connection refused
  ```

## Disclaimer

This tool is intended for educational purposes and stress-testing your own servers. **Do not use this tool to attack servers without explicit permission.** Unauthorized use may violate laws and result in severe consequences.
