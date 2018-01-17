package handlers

import (
	"net/http"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/swarmfund/api/db2/api"
	"gitlab.com/swarmfund/api/utils"
	"gitlab.com/tokend/api/geoinfo"
)

func CheckIsAuthorizedDevice(r *http.Request, deviceInfo *api.DeviceInfo, email string, walletID int64) error {
	var err error
	fingerprint, err := deviceInfo.Fingerprint()
	if err != nil {
		return errors.Wrap(err, "Unable to get fingerprint of the device")
	}

	device, err := AuthorizedDeviceQ(r).ByFingerprint(fingerprint)
	if err != nil {
		return errors.Wrap(err, "Unable to get authorized device by fingerprint")
	}

	if device != nil {
		err = AuthorizedDeviceQ(r).UpdateLastLoginTime(device)
		if err != nil {
			return errors.Wrap(err, "Failed to update last login time")
		}
	}

	newDevice := api.AuthorizedDevice{
		WalletID:    walletID,
		Details:     *deviceInfo,
		Fingerprint: fingerprint,
	}

	err = AuthorizedDeviceQ(r).Create(&newDevice)
	if err != nil {
		return errors.Wrap(err, "Failed to add authorized device")
	}

	err = Notificator(r).SendNewDeviceLogin(email, newDevice)
	if err != nil {
		return errors.Wrap(err, "New device login emails sending failed")
	}
	return nil
}

func GetSenderDeviceInfo(w http.ResponseWriter, r *http.Request, userIdentifier string) (*api.DeviceInfo, error) {
	var deviceInfo api.DeviceInfo
	err := deviceInfo.InitFormRequest(r)
	if err != nil {
		return nil, err
	}

	cookie, err := r.Cookie(utils.DeviceUIDCookieName(userIdentifier))
	if err != nil {
		cookie = utils.DeviceUIDCookie(userIdentifier, Notificator(r).ClientDomain())
	} else {
		utils.UpdateCookieExpires(cookie)
	}

	http.SetCookie(w, cookie)
	deviceInfo.DeviceUID = cookie.Value

	locationInfo, err := geoinfo.GetLocationInfo(deviceInfo.IP)
	if err != nil {
		deviceInfo.Location = "Unknown"
	} else {
		deviceInfo.Location = locationInfo.FullRegion()
	}

	return &deviceInfo, nil
}
