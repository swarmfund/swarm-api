package listener

import "gitlab.com/tokend/horizon-connector/internal/resources"

type ReviewableRequestEvent struct {
	body *resources.Request
	err  error
}

func (e *ReviewableRequestEvent) Unwrap() (*resources.Request, error) {
	return e.body, e.err
}

func (e *ReviewableRequestEvent) GetLoganFields() map[string]interface{} {
	return map[string]interface{} {
		"body": e.body,
		"err":  e.err,
	}
}
