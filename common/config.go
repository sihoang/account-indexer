package common

import (
	"fmt"
	"sync"
	"time"
)

// Configuration for the indexer
type Configuration struct {
	CleanInterval   time.Duration
	BlockTTL        time.Duration
	WatcherInterval time.Duration
	// OOSThreshold threshold
	OOSThreshold time.Duration
	Port         int
	NumBatch     int
}

func (con *Configuration) String() string {
	return fmt.Sprintf("CleanInterval=%v BlockTTL=%v WatcherInterval=%v OOSThreshold=%v Port=%v NumBatch=%v",
		con.CleanInterval, con.BlockTTL, con.WatcherInterval, con.OOSThreshold, con.Port, con.NumBatch)
}

var config *Configuration
var once sync.Once

// GetConfig Singleton
func GetConfig() *Configuration {
	once.Do(func() {
		config = &Configuration{}
	})
	return config
}
