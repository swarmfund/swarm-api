package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	. "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

type ActionName string

const SignUpForNewsletter ActionName = "sign_up_for_newsletter"
const ReadSwarmBasics ActionName = "read_swarm_basics"
const ViewDownloadWhitepaper ActionName = "view_download_whitepaper"
const WatchedYtVideo ActionName = "watched_yt_video"
const VisitedFaq ActionName = "visited_faq"
const ClickedLoginButton ActionName = "clicked_login_button"
const Register ActionName = "register"
const VerifyEmail ActionName = "verify_email"
const FirstLogin ActionName = "first_login"
const Login ActionName = "login"
const Logout ActionName = "logout"
const BrowseIo ActionName = "browse_io"
const ViewYtVideo ActionName = "view_yt_video"
const ViewedASpecificTab ActionName = "viewed_a_specific_tab"
const Unsubscribe ActionName = "unsubscribe"
const VoteInAnIo ActionName = "vote_in_an_io"

var actionExists = map[ActionName]bool{
	SignUpForNewsletter:    true,
	ReadSwarmBasics:        true,
	ViewDownloadWhitepaper: true,
	WatchedYtVideo:         true,
	VisitedFaq:             true,
	ClickedLoginButton:     true,
	Register:               true,
	VerifyEmail:            true,
	FirstLogin:             true,
	Login:                  true,
	Logout:                 true,
	BrowseIo:               true,
	ViewYtVideo:            true,
	ViewedASpecificTab:     true,
	Unsubscribe:            true,
	VoteInAnIo:             true,
}

var actionNameToSalesforceSphere = map[ActionName]string{
	SignUpForNewsletter:    "Sign up for Newsletter",
	ReadSwarmBasics:        "Read Swarm Basics",
	ViewDownloadWhitepaper: "View/Download Whitepaper",
	WatchedYtVideo:         "Watched YT video",
	VisitedFaq:             "Visited FAQ",
	ClickedLoginButton:     "Clicked Login Button",
	Register:               "Register",
	VerifyEmail:            "Verify Email",
	FirstLogin:             "First Login",
	Login:                  "Login",
	Logout:                 "Logout",
	BrowseIo:               "Browse IO",
	ViewYtVideo:            "View YT Video",
	ViewedASpecificTab:     "Viewed a specific tab",
	Unsubscribe:            "Unsubscribe",
	VoteInAnIo:             "Vote in an IO",
}

var actionNameToSalesforceActionName = map[ActionName]string{
	SignUpForNewsletter:    "User Activity",
	ReadSwarmBasics:        "Education",
	ViewDownloadWhitepaper: "Education",
	WatchedYtVideo:         "Education",
	VisitedFaq:             "Visited FAQ",
	ClickedLoginButton:     "User Activity",
	Register:               "Registration",
	VerifyEmail:            "Registration",
	FirstLogin:             "User Activity",
	Login:                  "User Activity",
	Logout:                 "User Activity",
	BrowseIo:               "Education",
	ViewYtVideo:            "Education",
	ViewedASpecificTab:     "Education",
	Unsubscribe:            "User Activity",
	VoteInAnIo:             "User Activity",
}

var Email = &emailRule{message: "must be a valid email"}

type emailRule struct {
	message string
}

func (v *emailRule) Validate(value interface{}) error {
	stringValue, ok := value.(*string)
	if !ok {
		return errors.New(v.message)
	}
	if !strings.ContainsRune(*stringValue, '@') || !strings.ContainsRune(*stringValue, '.') {
		return errors.New(v.message)
	}
	return nil
}

var ActionExists = &actionExistsRule{message: "must be a valid action", exists: actionExists}

type actionExistsRule struct {
	message string
	exists  map[ActionName]bool
}

func (v *actionExistsRule) Validate(value interface{}) error {
	name, ok := value.(*ActionName)
	if !ok {
		return errors.New(v.message)
	}
	if !v.exists[*name] {
		return errors.New(v.message)
	}
	return nil
}

type (
	ProxyActionRequest struct {
		Name       ActionName `json:"action_name"`
		ActorName  string     `json:"actor_name"`
		ActorEmail string     `json:"actor_email"`
		Time       *time.Time `json:"action_time"`
	}
)

func NewProxyActionRequest(r *http.Request) (ProxyActionRequest, error) {
	request := ProxyActionRequest{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}
	return request, request.Validate()
}

func (r ProxyActionRequest) Validate() error {
	err := Errors{
		"action_name": Validate(&r.Name, Required, ActionExists),
		"actor_name":  Validate(&r.ActorName, Required),
		"actor_email": Validate(&r.ActorEmail, Required, Email),
		"action_time": Validate(&r.Time, Required),
	}.Filter()
	return err
}

/*
Example POST-request body:
{
	"action_name": "verify_email",
	"action_time": "2002-10-02T10:00:00-05:00",
	"actor_name": "Ivan Ianov",
	"actor_email": "ivan.ivanov@gmail.com"
}
*/

func ProxyAction(w http.ResponseWriter, r *http.Request) {
	request, err := NewProxyActionRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	Salesforce(r).SendEvent(actionNameToSalesforceSphere[ActionName(request.Name)], actionNameToSalesforceActionName[ActionName(request.Name)], request.Time, request.ActorName, string(request.ActorEmail), 0, "")

	w.WriteHeader(204)

}
