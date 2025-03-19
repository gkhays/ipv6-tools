package com.ipv6tools;

import java.net.Inet6Address;
import java.net.InetSocketAddress;
import java.net.ServerSocket;
import java.net.Socket;
import java.io.*;
import java.time.LocalDateTime;
import java.time.format.DateTimeFormatter;

public class IPv6Tester {
    private static final int DEFAULT_PORT = 8080;
    private static final String DEFAULT_IPV6_ADDRESS = "::1";
    private static final DateTimeFormatter formatter = DateTimeFormatter.ofPattern("yyyy-MM-dd HH:mm:ss");

    public static void main(String[] args) {
        if (args.length < 1 || args.length > 3) {
            printUsage();
            System.exit(1);
        }

        String mode = args[0];
        String ipv6Address = args.length > 1 ? args[1] : DEFAULT_IPV6_ADDRESS;
        int port = args.length > 2 ? parsePort(args[2]) : DEFAULT_PORT;

        if (!mode.equals("server") && !mode.equals("client")) {
            printUsage();
            System.exit(1);
        }

        try {
            if (mode.equals("server")) {
                runServer(ipv6Address, port);
            } else {
                runClient(ipv6Address, port);
            }
        } catch (IOException e) {
            System.err.println("Error: " + e.getMessage());
            e.printStackTrace();
        }
    }

    private static void printUsage() {
        System.out.println("Usage: java IPv6Tester <server|client> [ipv6_address] [port]");
        System.out.println("  server|client    - Required. Run as server or client");
        System.out.println("  ipv6_address     - Optional. IPv6 address (default: ::1)");
        System.out.println("  port             - Optional. Port number (default: 8080)");
        System.out.println("\nExamples:");
        System.out.println("  java IPv6Tester server");
        System.out.println("  java IPv6Tester server 2001:db8::1");
        System.out.println("  java IPv6Tester client 2001:db8::1 8888");
    }

    private static int parsePort(String portStr) {
        try {
            int port = Integer.parseInt(portStr);
            if (port < 1 || port > 65535) {
                System.err.println("Error: Port must be between 1 and 65535");
                System.exit(1);
            }
            return port;
        } catch (NumberFormatException e) {
            System.err.println("Error: Invalid port number");
            System.exit(1);
        }
        return DEFAULT_PORT; // Will never reach here due to System.exit
    }

    private static void runServer(String ipv6Address, int port) throws IOException {
        try (ServerSocket serverSocket = new ServerSocket()) {
            // Bind to specified IPv6 address
            serverSocket.bind(new InetSocketAddress(ipv6Address, port));
            System.out.println("IPv6 Server started on [" + ipv6Address + "]:" + port);


                try (Socket clientSocket = serverSocket.accept();
                     PrintWriter out = new PrintWriter(clientSocket.getOutputStream(), true);
                     BufferedReader in = new BufferedReader(new InputStreamReader(clientSocket.getInputStream()))) {

                    String clientAddress = clientSocket.getInetAddress().getHostAddress();
                    System.out.println("Client connected from: [" + clientAddress + "]");

		    while (true) {
                       // Read client message
                       String message = in.readLine();
                       System.out.println("Received from client: " + message);

                       // Send response with timestamp
                       String response = "Server received your message at " + LocalDateTime.now().format(formatter) + "address " + serverSocket.getInetAddress().getHostAddress();
                       out.println(response);
                       
                       // Add a delay of 1 second
                       Thread.sleep(1000);
		    }
                } catch (InterruptedException e) {
                    Thread.currentThread().interrupt();
                    System.err.println("Sleep interrupted: " + e.getMessage());
                }

        }
    }

    private static void runClient(String ipv6Address, int port) throws IOException {
        try (Socket socket = new Socket()) {
            // Connect to specified IPv6 address
            socket.connect(new InetSocketAddress(ipv6Address, port));
            System.out.println("Connected to server at [" + ipv6Address + "]:" + port);

            try (PrintWriter out = new PrintWriter(socket.getOutputStream(), true);
                 BufferedReader in = new BufferedReader(new InputStreamReader(socket.getInputStream()))) {

                for (int i = 0; i < 20; i++) {
                    // Send message to server
                    String message = "Hello from IPv6 client at " + LocalDateTime.now().format(formatter);
                    out.println(message);
                    System.out.println("Sent to server: " + message);

                    // Read server response
                    String response = in.readLine();
                    System.out.println("Server response: " + response);

                    // Wait 1 second before next iteration
                    if (i < 19) {
                        Thread.sleep(1000);
                    }
                }
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
                System.err.println("Sleep interrupted: " + e.getMessage());
            }
        }
    }
} 
