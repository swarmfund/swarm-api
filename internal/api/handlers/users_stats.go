package handlers

import (
	"net/http"

	"encoding/json"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/swarmfund/api/internal/api/movetoape"
	"gitlab.com/swarmfund/api/internal/api/resources"
	"gitlab.com/tokend/go/doorman"
)

type UsersStatsResponse struct {
	Data resources.UserStats `json:"data"`
}

func UserStats(w http.ResponseWriter, r *http.Request) {
	if err := Doorman(r, doorman.SignerOf(CoreInfo(r).GetMasterAccountID())); err != nil {
		movetoape.RenderDoormanErr(w, err)
		return
	}

	totalUsers, err := UsersQ(r).TotalRegistrations()
	if err != nil {
		Log(r).WithError(err).Error("failed to get total user registrations")
		ape.Render(w, problems.InternalError())
		return
	}

	totalKYCApprovals, err := UsersQ(r).TotalKYCApprovals()
	if err != nil {
		Log(r).WithError(err).Error("failed to get total user total KYC approvals")
		ape.Render(w, problems.InternalError())
		return
	}

	totalKYCApplications, err := UsersQ(r).TotalKYCApplications()
	if err != nil {
		Log(r).WithError(err).Error("failed to get total user total KYC applications")
		ape.Render(w, problems.InternalError())
		return
	}

	response := UsersStatsResponse{
		Data: resources.UserStats{
			TotalRegistrations:   totalUsers,
			TotalKYCApplications: totalKYCApplications,
			TotalKycApprovals:    totalKYCApprovals,
		},
	}

	json.NewEncoder(w).Encode(&response)
}
