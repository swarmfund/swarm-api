package api

import (
	"time"

	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"

	"gitlab.com/swarmfund/api/geoinfo"
)

type LogInfo struct {
	GeoInfo *geoinfo.LocationInfo `json:"geoinfo"`
}

type AuditLogAction struct {
	ID          int64     `db:"id"`
	ActionType  int32     `db:"action"`
	UserAddress string    `db:"user_address"`
	Details     LogInfo   `db:"details"`
	PerformedAt time.Time `db:"performed_at"`
}

func (d *LogInfo) Scan(src interface{}) error {
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

func (d LogInfo) Value() (driver.Value, error) {
	return json.Marshal(d)
}

func (a AuditLogAction) GetIP() string {
	if a.Details.GeoInfo != nil {
		return a.Details.GeoInfo.Ip
	}
	return ""
}
