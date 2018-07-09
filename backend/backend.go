package backend

import (
	"context"
	"os"

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

type Backend interface {
	GetNodes(ctx context.Context) ([]collections.Node, error)
	RegisterNode(ctx context.Context, node string, nodeData string) error
	Watch(ctx context.Context, eventChan chan int)
	RegisterService(ctx context.Context) error
	GetServices(ctx context.Context) error
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
