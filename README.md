# CNLookup - CNAME Resolution Chain Lookup Tool

CNLookup is a command-line tool for querying domain CNAME resolution chains and their final resolved IP addresses. This tool is particularly useful for network troubleshooting, DNS resolution analysis, and website hosting migration analysis.

## Features

- Tracks complete CNAME resolution chains
- Supports custom DNS servers
- Automatically detects and displays Unicode representations of Punycode domains
- Detects CNAME circular references
- Distinguishes between IPv4 and IPv6 addresses
- Provides detailed resolution statistics

## Installation

### Pre-built Binaries

You can download pre-compiled binaries for your platform from the [GitHub Releases](https://github.com/Mxmilu666/cnlookup/releases) page.

### Building from Source

Ensure you have Go installed (Go 1.16 or higher recommended), then run:

```bash
git clone https://github.com/Mxmilu666/cnlookup.git
cd cnlookup
go build
```

## Usage

Basic usage:

```bash
cnlookup example.com
```

Using a specific DNS server:

```bash
cnlookup -s 8.8.8.8 example.com
```

Or:

```bash
cnlookup -d example.com -s 8.8.8.8:53
```

### Parameters

- `-d <domain>`: Domain to query
- `[-s <server>]`: Custom DNS server in `IP:PORT` format (e.g., `8.8.8.8:53`)

## Output Example

```
CNAME Lookup: www.example.com
Using DNS server: 8.8.8.8:53
--------------------------------------------------
CNAME resolution chain:
  [1] www.example.com → example.com.cdn.cloudflare.net
  [2] example.com.cdn.cloudflare.net → example.com.cdn.cloudflare.net (No CNAME record)
--------------------------------------------------
IP addresses:
  IPv4: 104.21.30.175
  IPv4: 172.67.215.240
  IPv6: 2606:4700:3035::6815:1eaf
  IPv6: 2606:4700:3036::ac43:d7f0

Found 2 IPv4 address(es) and 2 IPv6 address(es)
--------------------------------------------------
Resolution summary: www.example.com -> example.com.cdn.cloudflare.net (2 level CNAME chain)
```

## Special Features

### Punycode Domain Support

For Internationalized Domain Names (IDN) like `xn--qei482atfla8286c994a.love`, the tool automatically displays their Unicode form `❤山田リョウ.love`.

### CNAME Loop Detection

When a circular reference in CNAME records is detected, the tool displays a warning and interrupts the query.

## Dependencies

- golang.org/x/net/idna - For handling Punycode domain conversions

## License

[MIT License](LICENSE)
