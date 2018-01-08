package actions

// JSON implementors can respond to a request whose response type was negotiated
// to be MimeHal or MimeJSON.
type JSON interface {
	JSON()
}

// Raw implementors can respond to a request whose response type was negotiated
// to be MimeRaw.
type Raw interface {
	Raw()
}
