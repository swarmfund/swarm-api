package notificator

type Letter struct {
	Header string
	Body   string
	Link   string
}

type LoginNoticeLetter struct {
	Header       string
	BrowserFull  string
	BrowserShort string
	Date         string
	DeviceFull   string
	DeviceShort  string
	Ip           string
	Location     string
}
