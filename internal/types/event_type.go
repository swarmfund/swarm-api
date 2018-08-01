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
	EventViewTabOverview
	EventViewTabTeam
	EventViewTabUpdates
	EventViewTabDocs
	EventViewTabValuations
	EventViewTabGovernance
	EventViewTabDefault
	EventUnsubscribe
	EventVoteInAnIo
	EventEnableTFA
	EventDisableTFA
	EventStartsKYC
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
	EventViewTabOverview:        "View Tab Overview",
	EventViewTabTeam:            "View Tab Team",
	EventViewTabUpdates:         "View Tab Updates",
	EventViewTabDocs:            "View Tab Docs",
	EventViewTabValuations:      "View Tab Valuations",
	EventViewTabGovernance:      "View Tab Governance",
	EventViewTabDefault:         "View Tab Default",
	EventUnsubscribe:            "Unsubscribe",
	EventVoteInAnIo:             "Vote in an IO",
	EventEnableTFA:              "Enable 2FA",
	EventDisableTFA:             "Disable 2FA",
	EventStartsKYC:              "Starts KYC",
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
	EventUnsubscribe:            "User Activity",
	EventVoteInAnIo:             "User Activity",
	EventEnableTFA:              "Security",
	EventDisableTFA:             "Security",
	EventStartsKYC:              "User Activity",
	EventViewTabOverview:        "User Activity",
	EventViewTabTeam:            "User Activity",
	EventViewTabUpdates:         "User Activity",
	EventViewTabDocs:            "User Activity",
	EventViewTabValuations:      "User Activity",
	EventViewTabGovernance:      "User Activity",
	EventViewTabDefault:         "User Activity",
}

func (r Event) GetSalesforceSphere() string {
	return eventToSalesforceSphere[r]
}

func (r Event) GetSalesforceActionName() string {
	return eventToSalesforceActionName[r]
}
