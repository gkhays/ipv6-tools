# IPv6 Tools

[![Java](https://img.shields.io/badge/Java-17-red.svg)](https://www.oracle.com/java/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

A collection of tools for testing and experimenting with IPv6 connectivity. This project provides a simple client-server implementation to verify IPv6 communication.

## 🌟 Features

- IPv6 server implementation
- IPv6 client implementation
- Real-time message exchange
- Timestamp-based logging
- Configurable port and IPv6 address
- Support for both local and remote IPv6 connections
- Automatic listing of available IPv6 addresses on the host
- Multi-client support (up to 10 simultaneous connections)

## 📋 Prerequisites

- Java 17 or higher
- IPv6-enabled network environment
- Basic understanding of IPv6 addressing

## 🚀 Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/ipv6-tools.git
cd ipv6-tools
```

## 💻 Usage

The `IPv6Tester` class can be run in either server or client mode. Here are the basic usage patterns:

### Running as Server

```bash
java java/src/IPv6Tester.java server [ipv6_address] [port]
```

### Running as Client

```bash
java java/src/IPv6Tester.java client [ipv6_address] [port]
```

### Viewing Available IPv6 Addresses

To see all available IPv6 addresses on your system, simply run the tool without any arguments:

```bash
java java/src/IPv6Tester.java
```

This will display:
- Usage instructions
- List of all available IPv6 addresses on your system, showing the interface name and address
- Example commands

## 📝 Examples

1. View available IPv6 addresses and usage help:
```bash
java java/src/IPv6Tester.java
```

2. Start a server on the default IPv6 address (::1) and port (8080):
```bash
java java/src/IPv6Tester.java server
```

3. Start a server on a specific IPv6 address and port:
```bash
java java/src/IPv6Tester.java server 2001:db8:1234:5678::1 8888
```

4. Connect a client to the server:
```bash
java java/src/IPv6Tester.java client 2001:db8:1234:5678::1 8888
```

## 🔍 How it Works

- The server listens for incoming IPv6 connections on the specified address and port
- The client connects to the server and sends messages
- The server responds with timestamps and acknowledgment messages
- Both sides log all communication for debugging purposes
- The server can handle up to 10 simultaneous client connections
- When no arguments are provided, the tool displays available IPv6 addresses on the system

## 📊 Output Examples

### Server Mode Output

When running in server mode, you'll see output like this:
```
IPv6 Server started on [2001:db8:1234:5678::1]:8080
Maximum number of simultaneous clients: 10
Client connected from: [2001:db8:1234:5678::2]
Received from client [2001:db8:1234:5678::2]: Hello from IPv6 client at 2024-03-21 14:30:45
```

The server will:
1. Display the IPv6 address and port it's listening on
2. Show the maximum number of simultaneous clients allowed
3. Show when a client connects, including the client's IPv6 address
4. Log each message received from the client
5. Send back acknowledgment messages with timestamps

### Client Mode Output

When running in client mode, you'll see output like this:
```
Connected to server at [2001:db8:1234:5678::1]:8080
Sent to server: Hello from IPv6 client at 2024-03-21 14:30:45
Server response: Server received your message at 2024-03-21 14:30:45 at address 2001:db8:1234:5678::1
```

The client will:
1. Confirm successful connection to the server
2. Show each message sent to the server
3. Display the server's response with timestamps
4. Continue this pattern for 20 iterations with 1-second delays between messages

### Available IPv6 Addresses Output

When running without arguments, you'll see output like this:
```
Usage: java IPv6Tester <server|client> [ipv6_address] [port]
  server|client    - Required. Run as server or client
  ipv6_address     - Optional. IPv6 address (default: ::1)
  port             - Optional. Port number (default: 8080)

Available IPv6 addresses on this host:
  eth0: 2001:db8:1234:5678::1
  wlan0: 2001:db8:abcd:efgh::1
  lo: ::1

Java IPv6 preferences:
  java.net.preferIPv4Stack: false
  java.net.preferIPv6Addresses: false

Examples:
  java IPv6Tester server
  java IPv6Tester server 2001:db8:1234:5678::1
  java IPv6Tester client 2001:db8:1234:5678::1 8888
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by the need for simple IPv6 testing tools
- Built with Java's built-in networking capabilities
