package model

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type ByteSize float64

const (
	B ByteSize = 1 << (10 * iota)
	KB
	MB
	GB
	TB
)

var byteSizeOrderedKeys = [...]string{"B", "KiB", "MiB", "GiB", "TiB"}

var byteSizeMap = map[string]ByteSize{"B": B, "KiB": KB, "MiB": MB, "GiB": GB, "TiB": TB}

func byteSizeFromString(s string) ByteSize {
	s = strings.TrimSpace(s)
	if subs := strings.Split(s, " "); len(subs) == 2 {
		if v, err := strconv.ParseFloat(subs[0], 64); err == nil {
			return ByteSize(v) * byteSizeMap[subs[1]]
		}
	}
	return 0
}

func (bs ByteSize) ToString() string {
	var k string
	for _, k = range byteSizeOrderedKeys {
		if d := bs / byteSizeMap[k]; d < (1 << 10) {
			return fmt.Sprintf("%v %v", d, k)
		}
	}
	return fmt.Sprintf("%v %v", bs/byteSizeMap[k], k)
}

type SessionStatus struct {
	Username    string
	IPAddress   string
	SessionTime string
	TimeLeft    string
	BytesDown   ByteSize
	BytesUp     ByteSize
}

func (ss *SessionStatus) IsValid() bool {
	return ss.TimeLeft != ""
}

func SessionStatusFromHTML(doc *goquery.Document) *SessionStatus {
	ss := &SessionStatus{}
	ss.Username = doc.Find("body > div.limiter > div > div > form:nth-child(2) > table > tbody > tr > td > table > tbody > tr:nth-child(1) > td:nth-child(2)").Text()
	ss.IPAddress = doc.Find("body > div.limiter > div > div > form:nth-child(2) > table > tbody > tr > td > table > tbody > tr:nth-child(2) > td:nth-child(2)").Text()
	ss.SessionTime = doc.Find("body > div.limiter > div > div > form:nth-child(2) > table > tbody > tr > td > table > tbody > tr:nth-child(3) > td:nth-child(2)").Text()
	ss.TimeLeft = doc.Find("body > div.limiter > div > div > form:nth-child(2) > table > tbody > tr > td > table > tbody > tr:nth-child(4) > td:nth-child(2)").Text()
	ss.BytesUp = byteSizeFromString(doc.Find("body > div.limiter > div > div > form:nth-child(2) > table > tbody > tr > td > table > tbody > tr:nth-child(5) > td:nth-child(2)").Text())
	ss.BytesDown = byteSizeFromString(doc.Find("body > div.limiter > div > div > form:nth-child(2) > table > tbody > tr > td > table > tbody > tr:nth-child(6) > td:nth-child(2)").Text())

	return ss
}
