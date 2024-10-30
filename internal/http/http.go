package http

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"context"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-sharif/net/internal/model"
	"github.com/spf13/viper"
	"github.com/tatsushid/go-fastping"
)

func PingTicker(ctx context.Context, cancelFunc context.CancelFunc, pingChan chan<- string, ip string) {
	pinger := fastping.NewPinger()
	rIP, err := net.ResolveIPAddr("ip4:icmp", ip)
	if err != nil {
		log.Println("Couldn't resolve ip address for the PingChecker ...")
	}
	pinger.AddIPAddr(rIP)

	retry := 5

	pinger.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		retry = 5
		pingChan <- fmt.Sprintf("Pinged %v in RTT %v", addr.IP, rtt)
	}
	pinger.OnIdle = func() {
		pingChan <- fmt.Sprintf("Failed to ping address %v", ip)
		if retry == 0 {
			cancelFunc()
		}
		retry--
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := pinger.Run(); err != nil {
				// handle when we don't have the icmp permissions
				pingChan <- err.Error()
			}
		}
	}

}

type SessionStatusHandler struct {
	URL     string
	UseIP   bool
	history []model.SessionStatus
}

func (sh *SessionStatusHandler) addToHistory(sessionStatus *model.SessionStatus) {
	sh.history = append(sh.history, *sessionStatus)
}

func (sh *SessionStatusHandler) GetHistory() []model.SessionStatus {
	return sh.history
}

func (sh *SessionStatusHandler) ResetHistory() {
	sh.history = make([]model.SessionStatus, 100)
}

func (sh *SessionStatusHandler) Init() error {

	if sh.URL == "" {
		hostIP := viper.GetString("hostIP")
		hostDomain := viper.GetString("hostDomain")
		statusEndpoint := viper.GetString("statusEndpoint")
		if sh.UseIP {
			sh.URL = hostIP + statusEndpoint
		} else {
			sh.URL = hostDomain + statusEndpoint
		}
	}

	if sh.URL == "" {
		return fmt.Errorf("No session status url found (not configed nor specified)")
	}

	return nil
}

type SessionStatusDiff struct {
	BytesDown model.ByteSize
	BytesUp   model.ByteSize
}

func (sh *SessionStatusHandler) Diff() *SessionStatusDiff {
	l := len(sh.history)

	if l < 2 {
		return nil
	}

	h1 := sh.history[l-1]
	h2 := sh.history[l-2]

	return &SessionStatusDiff{
		BytesDown: h1.BytesDown - h2.BytesDown,
		BytesUp:   h1.BytesUp - h2.BytesUp,
	}
}

func (sh *SessionStatusHandler) GetSessionStatus(saveHistory bool) (int, *model.SessionStatus, error) {

	var tr *http.Transport
	if sh.UseIP {
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	client := &http.Client{Transport: tr}

	res, err := client.Get(sh.URL)

	if err != nil {
		return -1, nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return -1, nil, fmt.Errorf("Can not get the status %v", res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return 200, nil, err
	}

	ss := model.SessionStatusFromHTML(doc)
	if !ss.IsValid() {
		return 200, nil, fmt.Errorf("Failed to parse the session status from HTML")
	}

	if saveHistory {
		sh.addToHistory(ss)
	}

	return 200, ss, nil
}

type LoginHandler struct {
	URL      string
	Username string
	Password string
	UseIP    bool
}

func (lh *LoginHandler) Init() error {
	if lh.Username == "" {
		lh.Username = viper.GetString("username")
	}
	if lh.Password == "" {
		lh.Password = viper.GetString("password")
	}
	if lh.URL == "" {
		hostIP := viper.GetString("hostIP")
		hostDomain := viper.GetString("hostDomain")
		loginEndpoint := viper.GetString("loginEndpoint")
		if lh.UseIP {
			lh.URL = hostIP + loginEndpoint
		} else {
			lh.URL = hostDomain + loginEndpoint
		}
	}

	if lh.Username == "" || lh.Password == "" || lh.URL == "" {
		return fmt.Errorf("One of the required fields are missing(not specified nor found in the config)")
	}

	return nil
}

func (lh *LoginHandler) Login() (int, error) {
	body := bytes.NewBuffer([]byte(fmt.Sprintf("username=%v&password=%v", lh.Username, lh.Password)))

	r, err := http.NewRequest("POST", lh.URL, body)
	if err != nil {
		return 0, err
	}

	var tr *http.Transport
	if lh.UseIP {
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	client := &http.Client{Transport: tr}
	res, err := client.Do(r)
	if err != nil {
		return 0, err
	}

	defer res.Body.Close()

	return res.StatusCode, err
}
