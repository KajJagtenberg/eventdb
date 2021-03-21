package cluster

import (
	"strings"
	"time"

	"github.com/hashicorp/memberlist"
	"github.com/kajjagtenberg/eventflowdb/env"
)

type Cluster struct {
	list *memberlist.Memberlist
}

func (c *Cluster) Join() error {
	existingNodes := env.GetEnv("EXISTING_NODES", "")

	var existing []string

	if len(existingNodes) > 0 {
		existing = strings.Split(existingNodes, ",")
	}

	if _, err := c.list.Join(existing); err != nil {
		return err
	}

	return nil
}

func (c *Cluster) Leave() error {
	return c.list.Leave(time.Second * 10) // TODO: Change to customizable timeout
}

func (c *Cluster) GetList() *memberlist.Memberlist {
	return c.list
}

func NewCluster() (*Cluster, error) {
	conf := memberlist.DefaultLocalConfig()

	list, err := memberlist.Create(conf)
	if err != nil {
		return nil, err
	}

	return &Cluster{list}, nil
}
