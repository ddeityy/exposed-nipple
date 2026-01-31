package provider

import (
	"nipple/internal/config"
	"nipple/internal/rcon"

	"github.com/charmbracelet/log"
)

type Provider struct {
	lg         *log.Logger
	rconClient *rcon.Client
	cfg        *config.Config
}

func New(cfg *config.Config, rconClient *rcon.Client, lg *log.Logger) Provider {
	return Provider{
		lg:         lg,
		rconClient: rconClient,
		cfg:        cfg,
	}
}

func (p *Provider) Close() {
	p.rconClient.Close()
}

func (p *Provider) Logger() *log.Logger {
	return p.lg
}

func (p *Provider) RconClient() *rcon.Client {
	return p.rconClient
}

func (p *Provider) Config() *config.Config {
	return p.cfg
}
