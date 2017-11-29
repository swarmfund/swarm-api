package api

// common methods on `Action` for recovery request processing

func (action *Action) PurgeUser(accountID string) error {
	err := action.App.Storage().DeleteBucket(accountID)
	if err != nil {
		return err
	}

	err = action.APIQ().Wallet().SetActive(accountID, "")
	if err != nil {
		return err
	}

	// TODO make it cascade on wallet delete

	err = action.APIQ().Users().Delete(accountID)
	if err != nil {
		return err
	}

	return nil
}
