package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

func main() {
	serverMode := flag.Bool("server", false, "Whether or not to act as a server")
	clientMode := flag.Bool("client", false, "Whether or not to act as a client")
	ipv6Address := flag.String("address", "::1", "Optional. IPv6 address (default: ::1)")
	port := flag.Int("port", 8080, "Optional. Port number. Must be between 1 and 65535 (default: 8080)")

	flag.Parse()

	if *serverMode {
		err := runAsServer(*ipv6Address, *port)
		if err != nil {
			fmt.Printf("Error accepting client connection: %v\n", err)
			os.Exit(1)
		}
	}

	if *clientMode {
		fmt.Println("Not supported")
		os.Exit(1)
	}

	// Validate the IPv6 address
	if len(*ipv6Address) > 0 {
		err := error(nil)
		ip := net.ParseIP(*ipv6Address)
		if ip == nil {
			err = fmt.Errorf("invalid IPv6 address format")
		} else if ip.To4() != nil {
			err = fmt.Errorf("not an IPv6 address (IPv4 detected)")
		} else if ip.To16() == nil {
			err = fmt.Errorf("invalid IPv6 address")
		}

		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	}

	if *port != 8080 {
		if *port < 1 || *port > 65535 {
			fmt.Println("Error: Port must be between 1 and 65535")
			os.Exit(1)
		}
	}

	printIPv6Addresses()
}

func printIPv6Addresses() {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Error getting network interfaces", err)
		return
	}

	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			fmt.Printf("Error getting interface addresses for %s: %v\n", iface.Name, err)
			continue
		}

		fmt.Printf("%s\n", iface.Name)
		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok {
				if ipNet.IP.To16() != nil && ipNet.IP.To4() == nil {
					fmt.Printf("  IPv6 Address: %s\n", ipNet.IP.String())
				}
			}
		}
	}
}

func runAsClient(ipv6Addr string, port int) error {
	fullAddr := net.JoinHostPort(ipv6Addr, strconv.Itoa(port))
	conn, err := net.Dial("tcp6", fullAddr)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %v", fullAddr, err)
	}
	defer conn.Close()

	fmt.Printf("Connected to %s\n", fullAddr)

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		message := scanner.Text()
		fmt.Printf("Received from server: %s\n", message)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading from server: %v", err)
	}

	return nil
}

func runAsServer(ipv6Addr string, port int) error {
	fullAddr := net.JoinHostPort(ipv6Addr, strconv.Itoa(port))
	listener, err := net.Listen("tcp6", fullAddr)
	if err != nil {
		return fmt.Errorf("failed to bind to %s: %v", fullAddr, err)
	}
	defer listener.Close()

	fmt.Printf("IPv6 Server started on [%s]:%d\n", ipv6Addr, port)

	stopChannel := make(chan chan struct{})

	var wg sync.WaitGroup

	connectionSemaphore := make(chan struct{}, 10)

	// Handle interrupt signals
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		fmt.Println("\nReceived interrupt. Shutting down server...")
		close(stopChannel)
	}()

	go func() {
		for {
			select {
			case <-stopChannel:
				return
			default:
				// CAn we allow more connections?
				select {
				case connectionSemaphore <- struct{}{}:
					// Connection slot available
				case <-stopChannel:
					return
				}

				conn, err := listener.Accept()
				if err != nil {
					<-connectionSemaphore // Release the semaphore slot
					select {
					case <-stopChannel:
						return
					default:
						fmt.Println("Error accepting connection:", err)
						continue
					}
				}

				wg.Add(1)

				go func(conn net.Conn) {
					defer func() {
						<-connectionSemaphore // Release the semaphore slot
						wg.Done()
					}()
					handleConnection(conn)
				}(conn)
			}
		}
	}()

	<-stopChannel

	wg.Wait()
	fmt.Println("Server shutdown complete")
	return nil
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	remoteAddr := conn.RemoteAddr().String()
	fmt.Printf("Received from client [%s] (Max 10 connections)\n", remoteAddr)

	scanner := bufio.NewScanner(conn)

	_, err := conn.Write([]byte("Welcome to the IPv6 Server! (Connection limit: 10)\n"))
	if err != nil {
		fmt.Printf("Error sending welcome message to %s: %v\n", remoteAddr, err)
		return
	}

	for scanner.Scan() {
		message := scanner.Text()

		fmt.Printf("Received from client [%s]: %s\n", remoteAddr, message)

		// Send response with timestamp
		response := fmt.Sprintf("Server received your message at %s\n at address %s", time.Now(), conn.LocalAddr())
		_, err := conn.Write([]byte(response))
		if err != nil {
			fmt.Printf("Error sending response to %s: %v\n", remoteAddr, err)
			break
		}

		// Optional: exit command
		if message == "exit" {
			conn.Write([]byte("Goodbye!\n"))
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading from %s: %v\n", remoteAddr, err)
	}

	fmt.Printf("Connection from %s closed\n", remoteAddr)
}
