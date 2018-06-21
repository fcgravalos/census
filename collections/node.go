package collections

import "net"

type Node struct {
	Hostname  string            `json:"hostname"`
	OsFamily  string            `json:"osFamily"`
	NetIfaces []net.Interface   `json:"netIfaces"`
	Metadata  map[string]string `json:"metadata"`
}
