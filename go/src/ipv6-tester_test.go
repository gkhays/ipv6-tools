package main

import (
	"bufio"
	"net"
	"strings"
	"testing"
	"time"
)

func TestRunAsServer(t *testing.T) {
	// Start server in a goroutine
	go func() {
		err := runAsServer("::1", 8081)
		if err != nil {
			t.Errorf("Server failed: %v", err)
		}
	}()

	// Give the server time to start
	time.Sleep(100 * time.Millisecond)

	// Connect as a client
	conn, err := net.Dial("tcp6", "[::1]:8081")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	// Read welcome message
	reader := bufio.NewReader(conn)
	welcome, err := reader.ReadString('\n')
	if err != nil {
		t.Fatalf("Failed to read welcome message: %v", err)
	}

	if !strings.Contains(welcome, "Welcome to the IPv6 Server!") {
		t.Errorf("Unexpected welcome message: %s", welcome)
	}

	// Send a test message
	testMessage := "Hello, server!\n"
	_, err = conn.Write([]byte(testMessage))
	if err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}

	// Read response
	response, err := reader.ReadString('\n')
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}

	if !strings.Contains(response, "Server received your message") {
		t.Errorf("Unexpected response: %s", response)
	}

	// Test connection limit
	connections := make([]net.Conn, 0)
	defer func() {
		for _, c := range connections {
			c.Close()
		}
	}()

	// Try to create 11 connections (more than the 10 limit)
	for i := 0; i < 11; i++ {
		conn, err := net.Dial("tcp6", "[::1]:8081")
		if err != nil && i < 10 {
			t.Errorf("Failed to create connection %d: %v", i+1, err)
		} else if err == nil && i < 10 {
			connections = append(connections, conn)
		} else if err == nil && i >= 10 {
			t.Error("Server accepted more than 10 connections")
			conn.Close()
		}
	}
}
