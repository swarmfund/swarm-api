# Balances

### Account 🔒

All user balances are placed in account structure

```ruby
GET /accounts/GCZA6YVQBH46NYQ5PL2MXVA6U4PEFMNNHOB27N3GD6T5AU2YDFOP3ESJ

{
  "_links": {
    "self": {
      "href": "https://invest-dev.swarm.fund/accounts/GA43MKTMTE7NKDVLORPIZ6UVIL6DUH6L6WAVXB4JAHNAYH72K6HE5VEI"
    },
    "transactions": {
      "href": "https://invest-dev.swarm.fund/accounts/GA43MKTMTE7NKDVLORPIZ6UVIL6DUH6L6WAVXB4JAHNAYH72K6HE5VEI/transactions{?cursor,limit,order}",
      "templated": true
    },
    "operations": {
      "href": "https://invest-dev.swarm.fund/accounts/GA43MKTMTE7NKDVLORPIZ6UVIL6DUH6L6WAVXB4JAHNAYH72K6HE5VEI/operations{?cursor,limit,order}",
      "templated": true
    },
    "payments": {
      "href": "https://invest-dev.swarm.fund/accounts/GA43MKTMTE7NKDVLORPIZ6UVIL6DUH6L6WAVXB4JAHNAYH72K6HE5VEI/payments{?cursor,limit,order}",
      "templated": true
    }
  },
  "id": "GA43MKTMTE7NKDVLORPIZ6UVIL6DUH6L6WAVXB4JAHNAYH72K6HE5VEI",
  "account_id": "GA43MKTMTE7NKDVLORPIZ6UVIL6DUH6L6WAVXB4JAHNAYH72K6HE5VEI",
  "is_blocked": false,
  "block_reasons_i": 0,
  "block_reasons": [],
  "account_type_i": 5,
  "account_type": "AccountTypeNotVerified",
  "referrer": "",
  "thresholds": {
    "low_threshold": 0,
    "med_threshold": 0,
    "high_threshold": 0
  },
  "balances": [
    {
      "balance_id": "BBJ6SSJCEPOYL5OHDI2J4TZ44DSBXI2FIT2ZSNSM3WQSB2FOUBDV3SXK",
      "account_id": "GA43MKTMTE7NKDVLORPIZ6UVIL6DUH6L6WAVXB4JAHNAYH72K6HE5VEI",
      "asset": "BTC",
      "balance": "0.000000",
      "locked": "0.000000",
      "require_review": false
    },
    {
      "balance_id": "BBVI4C7EA54D2R3JLDHA3I4FGGB3MJX2A4RU6EDKJYWXI7KQLGKN2WEB",
      "account_id": "GA43MKTMTE7NKDVLORPIZ6UVIL6DUH6L6WAVXB4JAHNAYH72K6HE5VEI",
      "asset": "ETH",
      "balance": "0.000000",
      "locked": "0.000000",
      "require_review": false
    }
  ],
  "signers": [
    ...
  ],
  "limits": {
    "daily_out": "9223372036854.775807",
    "weekly_out": "9223372036854.775807",
    "monthly_out": "9223372036854.775807",
    "annual_out": "9223372036854.775807"
  },
  "statistics": {
    "daily_outcome": "0.000000",
    "weekly_outcome": "0.000000",
    "monthly_outcome": "0.000000",
    "annual_outcome": "0.000000"
  },
  "policies": {
    "account_policies_type_i": 0,
    "account_policies_types": null
  },
  "account_kyc": {
    "KYCData": {
      "blob_id": "MA36J5T2B4WPEQ3QUGMEA2RWHTDH7VL5PH4S3NHUALL5OLNMIJUA"
    }
  },
  "external_system_accounts": [
    {
      "type": {
        "name": "bitcoin",
        "value": 1
      },
      "data": "1LnCYq8zqJiEaii8QGyjBkUWjt4oK34SBC"
    },
    {
      "type": {
        "name": "ethereum",
        "value": 2
      },
      "data": "0xA7F7f6C40f05c2E5D8751E4b3cf7b5f680B820e3"
    }
  ],
  "referrals": []
}


HTTP/1.1 204
```
