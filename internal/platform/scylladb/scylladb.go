package scylladb

import (
	"time"

	"github.com/gocql/gocql"
	"github.com/phamduytien1805/package/config"
)

type ClusterManager struct {
	Cluster *gocql.ClusterConfig
}

func NewClusterManager(config *config.ScyllaConfig) *ClusterManager {
	cluster := createCluster(gocql.Quorum, config)

	return &ClusterManager{Cluster: cluster}
}

func createCluster(consistency gocql.Consistency, config *config.ScyllaConfig) *gocql.ClusterConfig {
	retryPolicy := &gocql.ExponentialBackoffRetryPolicy{
		Min:        time.Second,
		Max:        10 * time.Second,
		NumRetries: 5,
	}
	cluster := gocql.NewCluster(config.Hosts...)
	cluster.Keyspace = config.Keyspace
	cluster.Timeout = 5 * time.Second
	cluster.RetryPolicy = retryPolicy
	cluster.Consistency = consistency
	cluster.PoolConfig.HostSelectionPolicy = gocql.TokenAwareHostPolicy(gocql.RoundRobinHostPolicy())
	return cluster
}
