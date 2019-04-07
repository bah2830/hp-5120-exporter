package hp

import (
	"net"
	"regexp"
	"strconv"
	"strings"

	"github.com/bah2830/switch-exporter/pkg/networkswitch"
	"golang.org/x/crypto/ssh"
)

type Switch struct {
	clientConfig *ssh.ClientConfig
	client       *ssh.Client
	ip           string
	port         uint16
}

func NewWithPassword(ip string, port uint16, username, password string) (*Switch, error) {
	hpSwitch := &Switch{
		ip:   ip,
		port: port,
		clientConfig: &ssh.ClientConfig{
			User: username,
			Auth: []ssh.AuthMethod{
				ssh.Password(password),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Config: ssh.Config{
				Ciphers: []string{"aes128-cbc"},
			},
		},
	}

	client, err := ssh.Dial("tcp", net.JoinHostPort(ip, strconv.Itoa(int(port))), hpSwitch.clientConfig)
	if err != nil {
		return nil, err
	}

	hpSwitch.client = client

	return hpSwitch, nil
}

func (s *Switch) GetEnvironmentDetails() (*networkswitch.EnvironmentDetails, error) {
	output, err := s.getSSHResponse("dis env")
	if err != nil {
		return nil, err
	}

	r, err := regexp.Compile("hotspot (?P<name>\\d+)\\s+(?P<temp>\\S+)\\s+(?P<lower>\\S+)\\s+(?P<warning>\\S+)\\s+(?P<alarm>\\S+)\\s+(?P<alarm>\\S+)")
	if err != nil {
		return nil, err
	}

	details := &networkswitch.EnvironmentDetails{
		Sensors: make([]networkswitch.Sensor, 0),
	}

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if !strings.HasPrefix(line, "hotspot") {
			continue
		}

		matches := r.FindStringSubmatch(line)
		if len(matches) != 7 {
			continue
		}

		var indexToInt = func(match string) (int, error) {
			if match == "NA" {
				return 0, nil
			}

			output, err := strconv.Atoi(match)
			if err != nil {
				return 0, err
			}

			return output, nil

		}

		temp, err := indexToInt(matches[2])
		if err != nil {
			return nil, err
		}

		lower, err := indexToInt(matches[3])
		if err != nil {
			return nil, err
		}

		warning, err := indexToInt(matches[4])
		if err != nil {
			return nil, err
		}

		alarm, err := indexToInt(matches[5])
		if err != nil {
			return nil, err
		}

		critical, err := indexToInt(matches[6])
		if err != nil {
			return nil, err
		}

		details.Sensors = append(details.Sensors, networkswitch.Sensor{
			Name:        "hotspot " + matches[1],
			TempCelsius: temp,
			Limits: networkswitch.Limits{
				Lower:    lower,
				Warning:  warning,
				Alarm:    alarm,
				Critical: critical,
			},
		})

	}

	return details, nil
}

func (s *Switch) Disconnect() {
	s.client.Close()
}
