package backend

import (
	"os"
	"time"

	"github.com/fcgravalos/census/collections"
	log "github.com/sirupsen/logrus"
)

var supportedBackends map[string]bool

var logger *log.Logger

func init() {
	logger = log.New()
	supportedBackends = map[string]bool{
		"etcdv3":    true,
		"zookeeper": false,
		"redis":     false,
		"s3":        false,
	}
}

type event struct {
	ID          string
	timestamp   time.Time
	eventType   string
	description string
}

type Backend interface {
	GetNodes() ([]collections.Node, error)
	RegisterNode(node string, nodeData string) error
	//Watch(prefix string, eventChan chan event) error
	RegisterService() error
	GetServices() error
	Close()
}

func NewBackend(selectedBackend string) Backend {
	var backend Backend
	if supportedBackends[selectedBackend] {
		switch selectedBackend {
		case "etcdv3":
			backend := newEtcdv3Backend(logger)
			return backend
		}
	} else {
		logger.Errorf("unknown/disabled backend %s", selectedBackend)
		os.Exit(1)
	}
	return backend
}
