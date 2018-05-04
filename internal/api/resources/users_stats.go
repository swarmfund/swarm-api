package resources

type UserStats struct {
	TotalRegistrations   uint64 `json:"registrations"`
	TotalKYCApplications uint64 `json:"kyc_applications"`
	TotalKycApprovals    uint64 `json:"kyc_approvals"`
}
