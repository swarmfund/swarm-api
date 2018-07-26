package types

//go:generate jsonenums -tprefix=false -transform=snake -type=Event

type Event int32

const (
	EventSignUpForNewsletter Event = iota + 1
	EventReadSwarmBasics
	EventViewDownloadWhitepaper
	EventWatchedYtVideo
	EventVisitedFaq
	EventClickedLoginButton
	EventRegister
	EventVerifyEmail
	EventFirstLogin
	EventLogin
	EventLogout
	EventBrowseIo
	EventViewYtVideo
	EventViewedASpecificTab
	EventUnsubscribe
	EventVoteInAnIo
	EventEnableTFA
	EventDisableTFA
)

var eventToSalesforceActionName = map[Event]string{
	EventSignUpForNewsletter:    "Sign up for Newsletter",
	EventReadSwarmBasics:        "Read Swarm Basics",
	EventViewDownloadWhitepaper: "View/Download Whitepaper",
	EventWatchedYtVideo:         "Watched YT video",
	EventVisitedFaq:             "Visited FAQ",
	EventClickedLoginButton:     "Clicked Login Button",
	EventRegister:               "Register",
	EventVerifyEmail:            "Verify Email",
	EventFirstLogin:             "First Login",
	EventLogin:                  "Login",
	EventLogout:                 "Logout",
	EventBrowseIo:               "Browse IO",
	EventViewYtVideo:            "View YT Video",
	EventViewedASpecificTab:     "Viewed a specific tab",
	EventUnsubscribe:            "Unsubscribe",
	EventVoteInAnIo:             "Vote in an IO",
	EventEnableTFA:              "Enable 2FA",
	EventDisableTFA:             "Disable 2FA",
}

var eventToSalesforceSphere = map[Event]string{
	EventSignUpForNewsletter:    "User Activity",
	EventReadSwarmBasics:        "Education",
	EventViewDownloadWhitepaper: "Education",
	EventWatchedYtVideo:         "Education",
	EventVisitedFaq:             "Education",
	EventClickedLoginButton:     "User Activity",
	EventRegister:               "Registration",
	EventVerifyEmail:            "Registration",
	EventFirstLogin:             "User Activity",
	EventLogin:                  "User Activity",
	EventLogout:                 "User Activity",
	EventBrowseIo:               "Education",
	EventViewYtVideo:            "Education",
	EventViewedASpecificTab:     "Education",
	EventUnsubscribe:            "User Activity",
	EventVoteInAnIo:             "User Activity",
	EventEnableTFA:              "Security",
	EventDisableTFA:             "Security",
}

func (r Event) GetSalesforceSphere() string {
	return eventToSalesforceSphere[r]
}

func (r Event) GetSalesforceActionName() string {
	return eventToSalesforceActionName[r]
}
