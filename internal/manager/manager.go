package manager

import (
	"io"
	"net/http"
	"nipple/internal/config"
	"nipple/internal/rcon"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/log"
)

type ConnectManager struct {
	rconClient        *rcon.Client
	redirect          string
	redirectTimestamp time.Time
	serverPassword    string
	lg                *log.Logger
	Status
}

func NewConnectManager(rconClient *rcon.Client, cfg config.Server, lg *log.Logger) ConnectManager {
	return ConnectManager{
		rconClient:     rconClient,
		serverPassword: cfg.Password,
		lg:             lg,
	}
}

func (c ConnectManager) Close() error {
	return c.rconClient.Close()
}

type Status struct {
	Hostname string
	SDR      string
	Connect  string
	Map      string
	Players  int
}

func (c ConnectManager) GetServerStatus() (Status, error) {
	rawStatus, err := c.rconClient.GetServerStatus()
	if err != nil {
		return Status{}, err
	}

	return c.ParseStatus(rawStatus), nil
}

func (c ConnectManager) GetRedirectIP() string {
	if time.Since(c.redirectTimestamp) < 24*time.Hour {
		return c.redirect
	}

	client := http.Client{}
	resp, err := client.Get("https://potato.tf/api/serverstatus/redirect")
	if err != nil {
		log.Errorf("can't get redirect IP: %s", err)
		return c.redirect
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("can't read response body: %s", err)
		return c.redirect
	}

	c.redirectTimestamp = time.Now()

	return string(body)
}

func (c ConnectManager) SteamConnectURL() string {
	return "steam://connect/" + c.GetRedirectIP() + "/" + c.serverPassword + "/dest=" + c.SDR
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

	c.Connect = c.SteamConnectURL()

	return res
}
