package manager

import (
	"errors"
	"fmt"
	"nipple/internal/config"
	"nipple/internal/rcon"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/charmbracelet/log"
)

type ConnectManager struct {
	cfg               config.Config
	redirect          string
	redirectTimestamp time.Time
	serverPassword    string
	lg                *log.Logger
	Status
}

func NewConnectManager(cfg config.Config, lg *log.Logger) ConnectManager {
	return ConnectManager{
		cfg: cfg,
		lg:  lg,
	}
}

type Status struct {
	Status   string
	Hostname string
	Direct   string
	SDR      string
	Map      string
	Players  int
}

func (c ConnectManager) GetServerStatus() (Status, error) {
	rconClient, err := rcon.NewClient(c.cfg.RCON, c.lg)
	if err != nil {
		if errors.Is(err, syscall.ECONNREFUSED) || errors.Is(err, syscall.EHOSTUNREACH) {
			c.lg.Errorf("server %s unreachable, returning default status", c.cfg.RCON.Host)
			return Status{
				Status:   "Unreachable",
				Hostname: "N/A",
				Direct:   c.cfg.RCON.Host,
				SDR:      "N/A",
				Map:      "N/A",
				Players:  0,
			}, nil
		}
		return Status{}, fmt.Errorf("could not establish rcon connection: %w", err)
	}

	defer rconClient.Close()

	rawStatus, err := rconClient.GetServerStatus()
	if err != nil {
		return Status{}, err
	}

	return c.ParseStatus(rawStatus), nil
}

func (c ConnectManager) ParseStatus(s string) Status {
	res := Status{}

	for line := range strings.SplitSeq(s, "\n") {
		switch {
		case strings.HasPrefix(line, "hostname"):
			res.Hostname = strings.TrimSpace(strings.Split(line, ":")[1])
			continue
		case strings.HasPrefix(line, "udp/ip"):
			r := regexp.MustCompile(`\b(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}):(\d{1,5})\b`)
			match := r.FindAllString(line, -1)
			if len(match) > 0 {
				res.SDR = match[0]
			}
			continue
		case strings.HasPrefix(line, "map"):
			res.Map = strings.Split(strings.TrimSpace(strings.Split(line, ":")[1]), " ")[0]
			continue
		case strings.HasPrefix(line, "players"):
			players := strings.Split(strings.TrimSpace(strings.Split(strings.Split(line, ":")[1], ",")[0]), " ")
			res.Players, _ = strconv.Atoi(players[0])
			continue
		default:
			continue
		}
	}

	res.Status = "Online"
	res.Direct = c.cfg.RCON.Host

	return res
}
