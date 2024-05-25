package core

import (
	"github.com/shashimalcse/is-cli/internal/config"
	"go.uber.org/zap"
)

type CLI struct {
	Config config.Config
	Logger zap.Logger
}
