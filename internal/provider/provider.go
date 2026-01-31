package provider

import (
	"nipple/internal/config"
	"nipple/internal/manager"

	"github.com/charmbracelet/log"
)

type Provider struct {
	lg          *log.Logger
	connManager manager.ConnectManager
	cfg         *config.Config
}

func New(cfg *config.Config, connManager manager.ConnectManager, lg *log.Logger) Provider {
	return Provider{
		lg:          lg,
		connManager: connManager,
		cfg:         cfg,
	}
}

func (p *Provider) Close() error {
	return p.connManager.Close()
}

func (p *Provider) ConnManager() manager.ConnectManager {
	return p.connManager
}

func (p *Provider) Logger() *log.Logger {
	return p.lg
}

func (p *Provider) Config() *config.Config {
	return p.cfg
}
