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
	"gitlab.com/swarmfund/go/hash"
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

func (di DeviceInfo) Value() (driver.Value, error) {
	return json.Marshal(di)
}

func (di *DeviceInfo) Scan(src interface{}) error {
	v := reflect.ValueOf(src)
	if !v.IsValid() {
		return nil
	}
	if data, ok := src.([]byte); ok {
		err := json.Unmarshal(data, di)
		return err
	}
	return fmt.Errorf("could not not decode type %T -> %T", src, di)
}

const AppUAId = "SWARM_APP"

func (di *DeviceInfo) InitFormRequest(r *http.Request) error {
	var err error
	di.IP, _, err = net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return err
	}

	uaStr := r.Header.Get("User-Agent")
	if strings.Contains(uaStr, AppUAId) {
		di.parseSwarmAppUA(uaStr)
	} else {
		di.parseUA(uaStr)
	}

	if di.Browser == "" {
		di.Browser = "Unknown"
	}
	if di.BrowserFull == "" {
		di.BrowserFull = "Unknown"
	}
	if di.OS == "" {
		di.OS = "Unknown"
	}
	if di.OSFull == "" {
		di.OSFull = "Unknown"
	}

	return nil
}

func (di *DeviceInfo) parseSwarmAppUA(uaStr string) {
	// schema of custom UA from ours apps
	// SWARM_APP|<app_name>|<app_version>|<device_name>|<device_os>
	separated := strings.Split(uaStr, "|")
	if len(separated) != 5 {
		return
	}
	di.Browser = separated[1]
	di.BrowserFull = separated[1] + " " + separated[2]
	di.OS = separated[3]
	di.OSFull = separated[3] + " " + separated[4]
}

func (di *DeviceInfo) parseUA(uaStr string) {
	ua := user_agent.New(uaStr)
	osInfo := ua.OS()
	name, version := ua.Browser()

	di.Browser = name
	di.BrowserFull = name + " " + version
	di.OS = osInfo
	di.OSFull = strings.Replace(osInfo, "_", ".", -1)
}
