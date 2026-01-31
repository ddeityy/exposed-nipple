package rcon

import (
	"nipple/internal/config"

	"github.com/charmbracelet/log"

	"github.com/gorcon/rcon"
)

type Client struct {
	conn *rcon.Conn
	lg   *log.Logger
}

func NewClient(cfg config.RCON, lg *log.Logger) (*Client, error) {
	conn, err := rcon.Dial(cfg.Host, cfg.Password)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn,
		lg,
	}, nil
}

func (c *Client) GetServerStatus() (string, error) {
	out, err := c.conn.Execute("status")
	if err != nil {
		return "", err
	}

	return out, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
