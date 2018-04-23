package api

import (
	"database/sql/driver"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/mssola/user_agent"
	"gitlab.com/tokend/go/hash"
)

type DeviceInfo struct {
	IP          string `json:"ip"`
	Location    string `json:"location"`
	Browser     string `json:"browser"`
	BrowserFull string `json:"browser"`
	OS          string `json:"os"`
	OSFull      string `json:"os"`
	DeviceUID   string `json:"duid"`
}

type AuthorizedDevice struct {
	Id          int64      `db:"id"`
	WalletID    int64      `db:"wallet_id"`
	Fingerprint string     `db:"fingerprint"`
	Details     DeviceInfo `db:"details"`
	CreatedAt   time.Time  `db:"created_at"`
	LastLoginAt time.Time  `db:"last_login_at"`
}

func (di *DeviceInfo) Fingerprint() (fingerprint string, err error) {
	fpBytes, err := json.Marshal(di)
	if err != nil {
		return "", err
	}

	fpHash := hash.Hash(fpBytes)
	fingerprint = base64.StdEncoding.EncodeToString(fpHash[:])

	return fingerprint, nil
}

func (d DeviceInfo) Value() (driver.Value, error) {
	return json.Marshal(d)
}

func (d *DeviceInfo) Scan(src interface{}) error {
	v := reflect.ValueOf(src)
	if !v.IsValid() {
		return nil
	}
	if data, ok := src.([]byte); ok {
		err := json.Unmarshal(data, d)
		return err
	}
	return fmt.Errorf("could not not decode type %T -> %T", src, d)
}

const AppUAId = "BLCAPP"

func (d *DeviceInfo) InitFormRequest(r *http.Request) error {
	var err error
	d.IP, _, err = net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return err
	}

	uaStr := r.Header.Get("User-Agent")
	if strings.Contains(uaStr, AppUAId) {
		d.parseBlcAppUA(uaStr)
	} else {
		d.parseUA(uaStr)
	}

	if d.Browser == "" {
		d.Browser = "Unknown"
	}
	if d.BrowserFull == "" {
		d.BrowserFull = "Unknown"
	}
	if d.OS == "" {
		d.OS = "Unknown"
	}
	if d.OSFull == "" {
		d.OSFull = "Unknown"
	}

	return nil
}

func (d *DeviceInfo) parseBlcAppUA(uaStr string) {
	// schema of custom UA from ours apps
	// BLCAPP|<app_name>|<app_version>|<device_name>|<device_os>
	separated := strings.Split(uaStr, "|")
	if len(separated) != 5 {
		return
	}
	d.Browser = separated[1]
	d.BrowserFull = separated[1] + " " + separated[2]
	d.OS = separated[3]
	d.OSFull = separated[3] + " " + separated[4]
}

func (d *DeviceInfo) parseUA(uaStr string) {
	ua := user_agent.New(uaStr)
	osInfo := ua.OS()
	name, version := ua.Browser()

	d.Browser = name
	d.BrowserFull = name + " " + version
	d.OS = osInfo
	d.OSFull = strings.Replace(osInfo, "_", ".", -1)
}
