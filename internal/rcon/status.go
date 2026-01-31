package rcon

import (
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/log"
)

type command string

const (
	status string = "status"
)

type Status struct {
	Hostname          string
	SDR               string
	sdrPassword       string
	Connect           string
	Map               string
	Players           int
	redirect          string
	redirectTimestamp time.Time
}

func (s Status) GetRedirectIP() string {
	if time.Since(s.redirectTimestamp) < 24*time.Hour {
		return s.redirect
	}

	client := http.Client{}
	resp, err := client.Get("https://potato.tf/api/serverstatus/redirect")
	if err != nil {
		log.Errorf("can't get redirect IP: %s", err)
		return s.redirect
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("can't read response body: %s", err)
		return s.redirect
	}

	s.redirectTimestamp = time.Now()

	return string(body)
}

func (s Status) SteamConnectURL() string {
	return "steam://connect/" + s.GetRedirectIP() + "/" + s.sdrPassword + "/dest=" + s.SDR
}

func parseStatus(s string) Status {
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

	res.sdrPassword = "password"
	res.Connect = res.SteamConnectURL()

	return res
}
