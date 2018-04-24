# Javascript

Here are some common cases where using one of this libraries may largely simplify working with API.

## wallet.js

[View on github][wallet.js]

Probably, you will need wallet.js to make calculation of `walletId` and `walletKey` and then to decrypt keychain data using
`walletKey`. Methods reference:

### calculateWalletParams(password, email, salt, kdfParams)

```javascript
    import wallet from 'wallet.js'

    const walletParams = wallet.calculateWalletParams(
        'qwe123',
        'alice@mail.com',
        {
            "algorithm":"scrypt",
            "bits":256,
            "n":4096,
            "r":8,
            "p":1
        },
        "ivuWp73O0vkp/0vUekW6Xw=="
    )
    console.log(walletParams.walletId) // b6ea7c162a601f4c2ff00735ba1845016f98278eca202e7edbeb113014e9ae24
    console.log(walletParams.walletKey) // [998723670, -411326144, -1009165736, -798437100, 351233684, -730655273, 546018851, 1686940172]
```

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
|walletParams.walletId | string that identifies wallet |
|walletParams.walletKey | byte array that is used to decrypt keychain data |


### decryptKeychainData (keychainData, rawWalletKey)


```javascript
    import wallet from 'wallet.js'

    const walletKey = [998723670, -411326144, -1009165736, -798437100, 351233684, -730655273, 546018851, 1686940172]
    const encryptedKeychain = 'eyJJViI6IjN...I6ImdjbSJ9'
    const decryptedKeychain = wallet.decryptKeychainData(encryptedKeychain, walletKey)
    const keys = JSON.parse(decryptedKeychain)

    console.log(keys.seed) // SBCTRYOGBPDVEA7PLBVDYYCPFQGTGZS4QHVGHC6OWBLSGLAP5FUOY4CN
    console.log(keys.accountId) // GCZMY36JI2GGYJHY5MJYFOQMG4NDZTREEKGS4UL4MIGLFLKM533GN7CF
```

| Parameter | Description                |
| --------- | -------------------------- |
| keychainData    | encrypted keychain data |
| rawWalletKey    | wallet key, derived when using `calculateWalletParams`  |

#### return value

`decryptedKeychainData` - JSON string with object containing keys

| Property | Description                |
| -------- | -------------------------- |
| decryptedKeychainData.seed      | secret seed to sign requests and transactions |
| decryptedKeychainData.accountId | unique identifier of user account |


## swarm-js-base

js-base is a library for creating transactions, that you may need to send to Swarm core. There are a lot of transaction
types, and it contains builders for all of them.

### ManageAssetBuilder.assetCreationRequest(opts)

Creates operation to create asset creation request

```javascript
    import { ManageAssetBuilder } from 'swarm-js-sdk'

    const createAssetOp = ManageAssetBuilder.assetCreationRequest({
      requestID: '0',
      code: 'QTK',
      policies: 0,
      details: {
        name: 'AwesomeCoin',
        logo: {
          key: 'dpurah4infkubjhcost7fvprgkop4owgfxvfzeoxip5ni6rqh2otp2oq',
          type: 'image/png'
        },
        terms: {
          key: 'dpurah4infkubjhcost7fvprgkop4owgfxvfzeoxip5ni6rqh2otp2oq',
          type: 'application/pdf',
          name: 'AwesomeCoin terms.pdf'
        }
      },
      maxIssuanceAmount: '1000000.000000',
      initialPreissuedAmount: '1000000.000000,
      preissuedAssetSigner: 'GD2QOSZSKTSUP42AVEMFOJ7ZOBAF6GIQ4BMLVJNI3KZBGBIQIIK3L5PF'
    })
```

#### parameters
Type | Parameter                | Description                |
---- | ------------------------ | -------------------------- |
string | requestID              | request ID, if '0' - creates new, updates otherwise
string | code                   | Asset code
string | preissuedAssetSigner   | AccountID of keypair which will sign request for asset to be authrorized to be issued
string | maxIssuanceAmount      | Max amount can be issued of that asset
number | policies               | Asset policies
string | initialPreissuedAmount | Amount of pre issued tokens available after creation of the asset
object | details                | Additional details about asset
string | details.name           | Name of the asset
string | details.logo           | Details of the logo (see [documents] (#documents), here you need to previously upload 'asset_logo')
string | details.logo.key       | Key to compose asset picture url
string | details.logo.type      | Content type of asset logo
string | details.terms          | Asset terms
string | details.terms.type     | Content type of terms document
string | details.terms.name     | Name of terms document
string | source                 | The source account for the operation. Defaults to the transaction's source account.

### SaleRequestBuilder.createSaleCreationRequest(opts)

Creates operation to create sale creation request

```javascript
    import { SaleRequestBuilder } from 'swarm-js-sdk'

    const createSaleOp = SaleRequestBuilder.createSaleCreationRequest({
        defaultQuoteAsset: 'SUN',
        price: '1.231223',
        startTime: 1544572800
        endTime: 1602288000,
        softCap: '1000.000000',
        hardCap: '2000.000000',
        requestID: '0',
        baseAsset: 'QTK',
        quoteAssets: [
          {
              asset: 'BTC',
              price: '0.0021'
          },
          {
              asset: 'ETH',
              price: '0.0123'
          }
        ],
        details: {
          name: 'IO name',
          short_description: 'IO description, max 180 symbols',
          description: 'ZP5B3NPOKCG7OIHO4HMWCGT4QHMJVC3XTVERBDEYRRLY2P5THSDQ',
          logo: {
            key: 'dpurah4infkubjhcost7fvprgkop4owgfxvfzeoxip5ni6rqh2otp2oq',
            type: 'image/png'
          },
          youtube_video_id: 'GuReK0GsLFA'
        },
        isCrowdfunding: true
    })
```

#### parameters
Type | Parameter                   | Description                |
---- | --------------------------- | -------------------------- |
string | requestID                 |  ID of the request. 0 - to create new;
string | baseAsset                 |  asset for which sale will be performed
string | defaultQuoteAsset         |  asset in which hardcap/soft cap will be calculated
string | startTime                 |  start time of the sale
string | endTime                   |  close time of the sale
string | softCap                   |  minimum amount of quote asset to be received at which sale will be considered a successful
string | hardCap                   |  max amount of quote asset to be received
object | details                   |  sale specific details
object | details.name              |  name of the sale
object | details.short_description |  short description of the sale
object | details.description       |  blob ID of the description (see [blobs] (#blobs), here you need to previously upload `fund_overview`)
object | details.logo              |  details of the logo (see [documents] (#documents), here you need to previously upload `fund_logo`)
array  | quoteAssets               |  accepted assets
object | quoteAsset.price          |  price for 1 baseAsset in terms of quote asset
object | quoteAsset.asset          |  asset code of the quote asset
object | isCrowdfunding            |  states if sale type is crowd funding
string | source                    |  The source account for the operation. Defaults to the transaction's source account.

## swarm-js-sdk

[View on github][swarm-js-sdk]


swarm-js-sdk is a helper module to simplify the workflow with transactions backend modules. It also uses
[swarm-js-base] (#swarm-js-base) internally, but you don't need to import js-base itself, you can import js-base methods
from js-sdk itself

### server
js-sdk `server` is a factory of builders for making responses to horizon server. it contains lots of builders for making
`GET` requests for tx records, and also provides an interface to simplify crafting and sending a transaction.

### submitOperationGroup(operations, source, keypair)

```javascript
    import { Keypair, submitOperationGroup } from 'swarm-js-sdk'

    async function signAndSubmitTransaction () {
        const seed = SBCTRYOGBPDVEA7PLBVDYYCPFQGTGZS4QHVGHC6OWBLSGLAP5FUOY4CN
        const source = GCZMY36JI2GGYJHY5MJYFOQMG4NDZTREEKGS4UL4MIGLFLKM533GN7CF

        await submitOperationGroup ([createAssetOp, createSaleOp], source, Keypair.fromSecret(seed))
        console.log('done')
    }
```

| Parameter | Description                |
| --------- | -------------------------- |
| operations    | an array `Operation` objects, each representing an operation. (Every operation here will be included in transaction. If there is an error in any of operation, all transaction will be rejected) |
| source        | account ID that is the performs a transaction |
| keypair       | keypair derived from seed to sign the transaction |


[wallet.js]: https://github.com/swarmfund/wallet-js
[swarm-js-sdk]: https://github.com/swarmfund/swarm-js-sdk

### call builder

```javascript
    import { server, Keypair } from 'swarm-js-sdk'

    async function getBalances () {
        const respnse = await
            server
             .accounts()
             .accountId('GCZMY36JI2GGYJHY5MJYFOQMG4NDZTREEKGS4UL4MIGLFLKM533GN7CF')
             .callWithSignature(Keypair.fromSecret('SBCTRYOGBPDVEA7PLBVDYYCPFQGTGZS4QHVGHC6OWBLSGLAP5FUOY4CN'))
    }
```

Call builder is a helper module to make get requests. Call builder has simple structure for building a requests. Also it gives you a possibility
to sign request. All you need, is to call the `callWithSignature` at the and of chaining builders methods