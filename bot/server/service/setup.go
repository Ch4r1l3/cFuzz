package service

import (
	"github.com/hashicorp/go-plugin"
	"sync"
)

var running bool
var controlChan chan struct{}
var mutex sync.Mutex
var crashCheckMap map[string]bool

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "fuzz",
	MagicCookieValue: "fuzz",
}

func Setup() {
	running = false
	controlChan = make(chan struct{}, 1)
	crashCheckMap = make(map[string]bool)
}
