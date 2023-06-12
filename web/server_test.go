package web

import "testing"

func TestStart(t *testing.T) {
	server := NewHTTPServer(":9090")
	server.Start()
}
