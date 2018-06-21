package backend

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/fcgravalos/census/collections"
	log "github.com/sirupsen/logrus"
)

const (
	backendIdEtcdv3 = "etcdv3"
	etcdReqTimeout  = 5  // FIXME: should be config
	maxLeaseTime    = 30 // FIXME: should be config
	minTimeToWait   = 2 * time.Second
	maxTimeToWait   = 60 * time.Second
	nodesNamespace  = "nodes"
)

type etcdv3 struct {
	id     string
	logger *log.Logger
	client *clientv3.Client
}

// FIXME: RegisterNode should probably receive a context
func (e *etcdv3) RegisterNode(node string, data string) error {
	c := e.client
	lease, err := c.Grant(context.TODO(), maxLeaseTime)

	if err != nil {
		logger.Errorf("could'n grant lease for key: %v")
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	_, err = c.Put(ctx, fmt.Sprintf("%s/%s", nodesNamespace, node), data, clientv3.WithLease(lease.ID))
	cancel()

	if err != nil {
		logger.Errorf("Could not register node: %v", err)
		return err
	}
	return nil
}

func (e *etcdv3) GetNodes() ([]collections.Node, error) {
	var nodes []collections.Node
	c := e.client
	ctx, cancel := context.WithTimeout(context.Background(), etcdReqTimeout*time.Second)
	resp, err := c.Get(ctx, nodesNamespace, clientv3.WithPrefix())
	cancel()
	if err != nil {
		logger.Errorf("could not get nodes: %v", err)
		return nil, err
	}
	tmpNode := new(collections.Node)
	for _, kv := range resp.Kvs {
		if err = json.Unmarshal(kv.Value, &tmpNode); err != nil {
			logger.Errorf("error decoding json, skipping node: %v", err)
			continue
		}
		nodes = append(nodes, *tmpNode)
	}
	return nodes, nil
}

/*func (e *etcdv3) Watch(prefix string, eventChan ) {
	c := e.client
	go func() {
		for {
			ech := c.Watch("nodes")

		}
	}()


}
*/
func (e *etcdv3) RegisterService() error {
	fmt.Println("Registering service")
	return nil
}

func (e *etcdv3) GetServices() error {
	fmt.Println("mysql, redis, nginx")
	return nil
}

func (e *etcdv3) Close() {
	e.client.Close()
}

// We have followed the approach of waiting for the backend over and over
// We assume, there'd be an alert somewhere telling us that the backend is down
// That way, we don't couple census deployment to backend deployment
func newEtcdv3Backend(logger *log.Logger) *etcdv3 {
	var c *clientv3.Client
	cfg := clientv3.Config{
		Endpoints:   []string{"http://127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}
	retries := 0
	timeToWait := minTimeToWait
	for {
		if timeToWait > maxTimeToWait {
			timeToWait = minTimeToWait
		}
		cli, err := clientv3.New(cfg)
		if err != nil {
			logger.Errorf("couldn't connect to etcd: %v", err)
			logger.Warningf("retrying connection in %v", timeToWait)
			<-time.After(timeToWait)
			if retries > 0 {
				timeToWait = time.Duration(retries) * timeToWait
			}
		} else {
			c = cli
			break
		}
		retries++
	}
	return &etcdv3{id: backendIdEtcdv3, logger: logger, client: c}
}
