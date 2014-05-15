package storage

import (
	zklib "github.com/samuel/go-zookeeper/zk"

	"github.com/zenoss/serviced/coordinator/client"
	"github.com/zenoss/serviced/coordinator/client/zookeeper"
	"github.com/zenoss/serviced/domain/host"
	"github.com/zenoss/glog"

	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestClient(t *testing.T) {

	zookeeper.EnsureZkFatjar()
	basePath := "/basePath"
	tc, err := zklib.StartTestCluster(1)
	if err != nil {
		t.Fatalf("could not start test zk cluster: %s", err)
	}
	defer os.RemoveAll(tc.Path)
	defer tc.Stop()
	time.Sleep(time.Second)

	servers := []string{fmt.Sprintf("127.0.0.1:%d", tc.Servers[0].Port)}

	dsnBytes, err := json.Marshal(zookeeper.DSN{Servers: servers, Timeout: time.Second * 15})
	if err != nil {
		t.Fatal("unexpected error creating zk DSN: %s", err)
	}
	dsn := string(dsnBytes)
        zclient, err := client.New("zookeeper", dsn, basePath, nil)


	conn, err := zclient.GetConnection()
	if err != nil {
		t.Fatal("unexpected error getting connection")
	}

	h := host.New()
	h.ID = "nodeID"
	h.IPAddr = "192.168.1.5"
	c := NewClient(h, zclient)
	defer c.Close()
	time.Sleep(time.Second * 5)

	nodePath := fmt.Sprintf("/storage/clients/%s", h.IPAddr)
	glog.Infof("about to check for %s", nodePath)
	if exists, err := conn.Exists(nodePath); err != nil {
		t.Fatalf("did not expect error checking for existence of %s: %s", nodePath, err)
	} else {
		if !exists {
			t.Fatalf("could not find %s", nodePath)
		}
	}
}
