package api

type UserLimitReviewState int

const (
	UserLimitReviewNil UserLimitReviewState = iota
	UserLimitReviewPending
)
