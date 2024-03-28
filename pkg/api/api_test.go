package api

import (
	"bufio"
	"net"
	"testing"
)

func TestEndpoint(t *testing.T) {

	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		t.Fatalf("Error connecting to server: %s", err)
	}
	defer conn.Close()

	_, err = conn.Write([]byte("get_client_list"))
	if err != nil {
		t.Fatalf("Error sending request: %s", err)
	}

	resp, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		t.Fatalf("Error reading response: %s", err)
	}

	if resp != "[]" {
		t.Fatalf("Expected empty list, got: %s", resp)
	}

}
