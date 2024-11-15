package network

import "net"

func AnyAvailablePort() (int, error) {
	// Create a new TCP listener on port 0 (which means "any available port")
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}

	// Extract the port number from the listener address
	port := listener.Addr().(*net.TCPAddr).Port

	// Close the listener to free up the port
	err = listener.Close()

	return port, err
}
