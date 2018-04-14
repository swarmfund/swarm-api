# Libraries

Here are some common cases where using one of this libraries may largely simplify working with API.

## wallet.js

[View on github][wallet.js]

Probably, will need wallet.js to make calculation of `walletId` and `walletKey` and then to decrypt keychain data using
`walletKey`. Methods reference:

### calculalteWalletParams(password, email, salt, kdfParams)

#### parameters

| Parameter | Description                |
| -------- | -------------------------- |
| password | User provided password |
| email    | User provided email |
| kdfParams| Key derive function params, see [kdf] (#get-kdf-params) |
| salt     | Salt, comes with `kdfParams` |

#### return value

`walletParams`

| Property | Description                |
| --------- | -------------------------- |
|walletParams.walletId | wallet id that identifies wallet |
|walletParams.walletKey | wallet key that is used to decrypt keychain data |


### decryptKeychainData (keychainData, rawWalletKey)

| Parameter | Description                |
| --------- | -------------------------- |
| keychainData    | encrypted keychain data |
| rawWalletKey    | wallet key, derived when using `calculateWalletParams`  |

#### return value

`decryptedKeychainData`

| Property | Description                |
| -------- | -------------------------- |
| decryptedKeychainData.seed      | secret seed to sign requests and transactions |
| decryptedKeychainData.accountId | unique identifier of user account |


## swarm-js-sdk

[View on github][swarm-js-sdk]

[wallet.js]: https://github.com/swarmfund/wallet-js
[swarm-js-sdk]: https://github.com/swarmfund/swarm-js-sdk
