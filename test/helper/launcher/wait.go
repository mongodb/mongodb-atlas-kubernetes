package launcher

import (
	"fmt"
	"time"
)

var (
	NoWait = (*WaitConfig)(nil)
)

type WaitConfig struct {
	Condition string
	Target    string
	Timeout   time.Duration
}

func WaitReady(target string, timeout time.Duration) *WaitConfig {
	return &WaitConfig{Condition: "condition=Ready", Target: target, Timeout: timeout}
}

func (cfg *WaitConfig) waitArgs() []string {
	args := []string{"wait", fmt.Sprintf("--for=%s", cfg.Condition), cfg.Target}
	if cfg.Timeout == 0 {
		args = append(args, fmt.Sprintf("--timeout=%v", cfg.Timeout))
	}
	return args
}
