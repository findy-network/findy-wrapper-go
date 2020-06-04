// Package config includes configurations to setup libindy properly.
package config

import (
	"github.com/optechlab/findy-go/dto"
	"github.com/optechlab/findy-go/internal/c2go"
)

// Set is wrapper function for libindy to set runtime configuration. Please see
// more information from indy
// SDK documentation.
func Set(config SystemConfig) int {
	return c2go.FindySetRuntimeConfig(dto.ToJSON(config))
}

// SystemConfig is wrapper struct for libindy's corresponding JSON type.
type SystemConfig struct {
	CryptoThreadPoolSize int  `json:"crypto_thread_pool_size,omitempty"`
	CollectBacktrace     bool `json:"collect_backtrace,omitempty"`
}
