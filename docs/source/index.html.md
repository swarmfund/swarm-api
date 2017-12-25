---
title: Swarm API Reference

search: true

includes:
  - documents
  - kyc
---

# Overview

The Swarm API tries its best to follow [JSONAPI](http://jsonapi.org/format/1.0/).
Most important parts of protocol will be included here, but to get better feel of what's going on it's advised to get yourself familiar with JSONAPI before continuing.

| Legend | Description                |
| ------ | -------------------------- |
| 🔒     | Request requires signature |



## Making requests

Each request must include correct content negotiation headers.

# Errors

API uses conventional HTTP response codes to indicate the success or failure of a request.
In general, codes in the `2xx` range indicate success, codes in the `4xx` range indicate an error that 
failed given the information provided and codes in `5xx` range indicate an problem with Swarm servers.

## General HTTP response status codes

If not stated otherwise client should expect one of the following status codes:

| Status code      | Description                              |
| ---------------- | ---------------------------------------- |
| 200 OK           |                                          |
| 201 Created      |                                          |
| 204 No Content   | The server successfully processed the request and is not returning any content |
| 400 Bad Request  | Request was invalid in some way, response should contain details |
| 401 Unauthorized | Request signature was invalid or you are not authorized to access resource |

## Application-specific error codes

| Error code              | Description                              |
| ----------------------- | ---------------------------------------- |
| `tfa_required`          | See [docs section](#two-factor-authentication) for details |
| `verification_required` | Wallet with verified email is required before proceeding |

# Wallets

Used to store encrypted user keys used for signing requests.

## Get KDF params

```http
GET /kdf HTTP/1.1
Content-Type: application/vnd.api+json
Accept: application/vnd.api+json

HTTP/1.1 200 OK
Content-Type: application/vnd.api+json

{
    "data": {
        "attributes": {
            "algorithm": "scrypt", 
            "bits": 256, 
            "n": 4096, 
            "p": 1, 
            "r": 8
        }, 
        "id": "1", 
        "type": "kdf"
    }
}
```

Returns current default derivation parameters. Should be used for all new wallets.

| Parameter | Description                              |
| --------- | ---------------------------------------- |
| `email`   | will return KDF parameters for wallet with provided email, default if wallet does not exist |



## Create wallet

```http
POST /wallets HTTP/1.1
Accept: application/vnd.api+json
Content-Type: application/vnd.api+json

{
	"data": {
		"type": "wallet",
		"id": "388108095960430b80554ac3efb6807a9f286854033aca47f6f466094ab50876",
		"attributes": {
			"account_id": "GD6PPS6VCAN5AN52N2BSUJQTKW2T22AERA6HHI33VBW67T5GCDFWTVET",
			"email": "fo@ba.ar",
			"salt": "wAWBYjv5eSDVZsjY1suFFA==",
			"keychain_data": "eyJJViI6IjN...I6ImdjbSJ9"
		},
		"relationships": {
            "kdf": {
                "data": {
                    "type": "kdf",
                    "id": "1"
                }
            },
          	"factor": {
            	"data": {
            		"type": "password",
            		"attributes": {
            			"account_id": "GDI54FYDBF2S6GEQJHBLS3HMIEYYKDLVT7YCCI33K5J6B4JTGNP77DEK",
            			"keychain_data": "foo..bar",
            			"salt": "salt"
            		}
            	}
            }
        }
	}
}


HTTP/1.1 201
Content-Type: application/vnd.api+json

{
    "data": {
        "type": "wallet",
        "id": "388108095960430b80554ac3efb6807a9f286854033aca47f6f466094ab50876",
        "attributes": {
            "verified": false
        },
        "relationships": {
            "factor": {
                "data": {
                    "type": "password",
                    "id": "43"
                }
            }
        }
    },
    "included": [
        {
            "type": "password",
            "id": "43",
            "attributes": {
                "account_id": "GDI54FYDBF2S6GEQJHBLS3HMIEYYKDLVT7YCCI33K5J6B4JTGNP77DEK",
                "keychain_data": "foo...bar",
                "salt": "salt"
            }
        }
    ]
}
```

Create wallet requests should contain following resources:

#### Wallet

| Field                 | Description                              |
| --------------------- | ---------------------------------------- |
| `id`                  | lowercase, hex-encoded string with key derived from password, email and salt using KDF parameters |
| `account_idDocuments` | address of derived from wallet secret key |
| `salt`                | client generated salt                    |
| `keychain_data`       | secret key encrypted with password, email and salt |
| `email`               | wallet email address                     |

#### KDF

| Field | Description                              |
| ----- | ---------------------------------------- |
| `id`  | version of KDF parameters used to derive wallet data |

#### Factor

Second wallet with different seed encrypted by same email/password used as second factor to confirm password possession during [wallet update](#update-wallet) flow. 

### Response

Succeeded request will have response with current wallet state as well related factor `id`.

## Email verification

```http
PUT /wallets/388108095960430b80554ac3efb6807a9f286854033aca47f6f466094ab50876/verification HTTP/1.1
Content-Type: application/vnd.api+json

{
	"data": {
		"attributes": {
			"token": "JOqIgfCNSjnGWDrTPWbW"
		}	
	}
}

HTTP/1.1 204
```

Once signed up user should receive email with verification link with [client router payload](#client-redirects).

| Field       | Description              |
| ----------- | ------------------------ |
| `token`     | Email verification token |
| `wallet_id` | Related wallet id        |

Client then should send request to `PUT /wallets/{wallet-id}/verification` to confirm email address.


## Requesting email resend

```ruby
POST /wallets/5xkxn3nfjmfx9p58geo5stt9/verification HTTP/1.1

HTTP/1.1 204
```



## Get Wallet

```ruby
GET /wallets/isheebk4wi962pjpnvc1pp66r HTTP/1.1
Accept: application/vnd.api+json


HTTP/1.1 200
Content-Type: application/vnd.api+json

{
    "data": {
        "type": "wallet",
        "id": "isheebk4wi962pjpnvc1pp66r",
        "attributes": {
            "account_id": "GAAAAAAA",
            "email": "aff49g266gv5fpz3iyj6h6w29@test.com",
            "keychain_data": "foo",
            "verified": true
        }
    }
}
```

<aside class="notice">Previously knows as <code>POST /wallets/show</code></aside>



## Update wallet 🔒 

```ruby
PUT /wallets/5xkxn3nfjmfx9p58geo5stt9
Content-Type: application/vnd.api+json
Accept: application/vnd.api+json

{
	"data": {
		"type": "wallet",
		"id": "isheebk4wi962pjpnvc1pp66r",
		"attributes": {
			"account_id": "GAAAAAAA",
			"salt": "salt==",
			"keychain_data": "foo..bar"
		},
		"relationships": {
			"transaction": {
				"data": {
					"attributes": {
						"envelope": "AAA...AAAA"
					}
				}	
			},
            "kdf": {
                "data": {
                    "type": "kdf",
                    "id": "1"
                }
            },
            "factor": {
            	"data": {
            		"type": "password",
            		"attributes": {
            			"account_id": "GDI54FYDBF2S6GEQJHBLS3HMIEYYKDLVT7YCCI33K5J6B4JTGNP77DEK",
            			"keychain_data": "foo...bar",
            			"salt": "salt"
            		}
            	}
            }
        }
	}
}

HTTP/1.1 204
```

<aside class="notice">Previously known as <code>POST /wallets/update</code></aside>

## Index

```http
GET /wallets?page=1 HTTP/1.1
Content-Type: application/vnd.api+json
Accept: application/vnd.api+json

HTTP/1.1 200
Content-Type: application/vnd.api+json

{
    "data": [
        {
            "attributes": {
                "account_id": "GCCJPB7QQLNEMCJ72CQJ4ODAZFFGXHET5UPOGSDWT222GXQPMTO6ZQW3", 
                "email": "yr0a3ke29d78skm2030kep14i@test.com", 
                "keychain_data": "foo", 
                "verified": true
            }, 
            "id": "912svl5uj6hxtypnonryhkt9", 
            "type": "wallet"
        }
    ], 
    "links": {
        "next": "/wallets?page=3", 
        "prev": "/wallets?page=1", 
        "self": "/wallets?page=2"
    }
}


```

| Parameters | Description                     |
| ---------- | ------------------------------- |
| `page`     |                                 |
| `state`    | Mask to filter wallets by state |

### Wallet states

| State | Description                  |
| ----- | ---------------------------- |
| `1`   | email has not been confirmed |
| `2`   | email has been confirmed     |



# Client redirects

> link format

```shell
http://client.com/r/eyAic3RhdHVzIjoyMDAsImF...G4zbmZqbWZ4OXA1OGdlbzVzdHQ5In19
```

> decoded value

```json
{
	"status": 200,
	"type": 1,
	"meta": {
		"token": "JOqIgfCNSjnGWDrTPWbW",
		"wallet_id": "388108095960430b80554ac3efb6807a9f286854033aca47f6f466094ab50876"
	}
}
```

| Field    | Description                              |
| -------- | ---------------------------------------- |
| `status` | Action result code following HTTP status code semantics, might be omitted if 200 |
| `type`   | [Redirect type](#redirect-types)         |
| `meta`   | Types specific meta-information, might be omitted if empty |

## Redirect types

| Type | Description                              |
| ---- | ---------------------------------------- |
| `1`  | Email confirmation token, meta will include verification `token` and `wallet_id` |

# Users

#### Types

| Value          | Filter |
| -------------- | ------ |
| `not_verified` | 1      |
| `syndicate`    | 2      |

#### States

| Value                  | Filter | Description                      |
| ---------------------- | ------ | -------------------------------- |
| `nil`                  | 1      | Initial user state               |
| `waiting_for_approval` | 2      | User is waiting for KYC approval |
| `approved`             | 4      | User has approved KYC            |
| `rejected`             | 8      | User has rejected KYC            |



## Create 🔒

```http
PUT /users/GBT3XFWQUHUTKZMI22TVTWRA7UHV2LIO2BIFNRCH3CXWPYVYPTMXMDGC HTTP/1.1
Content-Type: application/vnd.api+json

{
	"data": {
		"attributes": {}
	}
}

HTTP/1.1 204
```

## Get 🔒

```http
GET /users/GBT3XFWQUHUTKZMI22TVTWRA7UHV2LIO2BIFNRCH3CXWPYVYPTMXMDGC HTTP/1.1
Content-Type: application/vnd.api+json


HTTP/1.1 200
Content-Type: application/vnd.api+json

{
    "data": {
        "type": "syndicate",
        "id": "GCCJPB7QQLNEMCJ72CQJ4ODAZFFGXHET5UPOGSDWT222GXQPMTO6ZQW3",
        "attributes": {
            "email": "test@test.com",
            "state": "nil"
        }
    }
}
```



## Index 🔒 

```http
GET /users HTTP/1.1
Content-Type: application/vnd.api+json

HTTP/1.1 200
Content-Type: application/vnd.api+json

{
    "data": [
        {
            "type": "syndicate",
            "id": "GCCJPB7QQLNEMCJ72CQJ4ODAZFFGXHET5UPOGSDWT222GXQPMTO6ZQW3",
            "attributes": {
                "email": "yr0a3ke29d78skm2030kep14i@test.com",
                "state": "waiting_for_approval"
            }
        }
    ],
    "links": {
        "self": "/users?page=1",
        "next": "/users?page=2"
    }
}
```



| Parameter | Description                              |
| --------- | ---------------------------------------- |
| `page`    | Pagination cursor                        |
| `state`   | Bit mask to filter users by state        |
| `type`    | Bit mask to filter users by type         |
| `email`   | Substring to match against user emails   |
| `address` | Substring to match against account addresses |



## Update 🔒

> set user type

```http
PATCH /users/GBT3XFWQUHUTKZMI22TVTWRA7UHV2LIO2BIFNRCH3CXWPYVYPTMXMDGC HTTP/1.1
Content-Type: application/vnd.api+json

{
    "data": {
        "type": "syndicate"
    }
}

HTTP/1.1 204
```

> approve user

```http
PATCH /users/ HTTP/1.1
Content-Type: application/vnd.api+json

{
	"data": {
		"attributes": {
			"state": "approved"
		},
		"relationships": {
			"transaction": {
				"data": {
					"attributes": {
						"envelope": "AAA...AAA"
					}
				}
			}	
		}
	}
}

HTTP/1.1 204
```



| Field                    | User | Admin | Description                              |
| ------------------------ | ---- | ----- | ---------------------------------------- |
| `/data/type`             | +    | -     | Updating user type is allowed only if current type is `not_verified` |
| `/data/attributes/state` | +    | +     | Used by admin to update state from `waiting_for_approval` to `rejected` or `approved`. User could try to set state to `waiting_for_approval` why KYC is ready for review. |
|                          |      |       |                                          |
|                          |      |       |                                          |



# Two-factor authentication

> example second factor error

```json
{
    "errors": [
        {
            "title": "Forbidden",
            "detail": "Additional factor required",
            "status": "403",
          	"meta": {
                "factor_id": 26,
              	"factor_type": "password",
                "token": "c64c45125a5a846dbfa8cbbaf1b7f3dcc8c3ce9d",
              	"wallet_id": "388108095960430b80554ac3efb6807a9f286854033aca47f6f466094ab50876"
            }
        }
    ]
}
```

Second factor verification is required when accessing some resources. Clients should expect `tfa_required` error code for any request and prompt user to provide OTP.

Subsequent requests will have same `token` . 

After successful OTP verification user should be able to access requested resource.

### Response

| Field         | Description                              |
| ------------- | ---------------------------------------- |
| `factor_id`   | ID of a factor which triggered verification |
| `token`       | Hash based on request                    |
| `wallet_id`   | ID of a wallet for which verification was triggered |
| `factor_type` | Type of the factor which triggered verification |



### Factor types

| Type       | Description                              |
| ---------- | ---------------------------------------- |
| `password` | Created during signup, updated only with change password flow |
| `totp`     | Optional Google Authenticator token      |



## Get factors

```ruby
GET /wallets/388108095960430b80554ac3efb6807a9f286854033aca47f6f466094ab50876/factors
Content-Type: application/vnd.api+json
Accept: application/vnd.api+json

HTTP/1.1 200
Content-Type: application/vnd.api+json

{
    "data": [
        {
            "type": "totp",
            "id": 24,
            "attributes": {
                "priority": 0
            }
        }
    ]
}
```

<aside class="notice">Previously known as <code>GET /tfa</code></aside>

## Create TOTP factor 🔒

```http
POST /wallets/388108095960430b80554ac3efb6807a9f286854033aca47f6f466094ab50876/factors HTTP/1.1
Content-Type: application/vnd.api+json
Accept: application/vnd.api+json

{
	"data": {
    	"type": "totp"
    }  
}


HTTP/1.1 201
Content-Type: application/vnd.api+json

{
    "data": {
        "id": 15,
        "type": "totp",
        "attributes": {
            "secret": "4FW4QD2NGGE45NOA",
            "seed": "otpauth://totp/%7B%7B%20.Project%20%7D%7D:fo@ob.ar?algorithm=SHA1&digits=6&issuer=%7B%7B+.Project+%7D%7D&period=30&secret=4FW4QD2NGGE45NOA"
        }
    }
}
```

`POST /wallets/{wallet-id}/factors`

Single TOTP factor is allowed per wallet, subsequent requests will result in `409 Conflict`

### Response

| Field    | Description                              |
| -------- | ---------------------------------------- |
| `secret` | HMAC secret                              |
| `seed`   | Google Authenticator URI for generating QR-code |



## Factor verification

```http
PUT /wallets/82dc746c10b28f87eb6eae695fc3c5ab91c31988ea43b692c8ae50f52f24b9d6/factors/31/verification HTTP/1.1
Content-Type: application/vnd.api+json

{
	"data": {
		"attributes": {
			"token": "f86e176643374de6d13690e499440229d43e2cb7",
			"otp": "697313"
		}
	}
}

HTTP/1.1 204
```

`PUT /wallets/{wallet-id}/factors/{factor-id}/verification`

| Field   | Description                              |
| ------- | ---------------------------------------- |
| `token` | Request based hash return by `tfa_required` error |
| `otp`   | Factor specific one-time password        |

### TOTP

For `totp` factor `otp` is current time based token provide by user's application.

### Password

For `password` factor `otp` is `token` signed with seed from factor keychain.

## Update factor

```ruby
PATCH /wallets/82dc746c10b28f87eb6eae695fc3c5ab91c31988ea43b692c8ae50f52f24b9d6/factors/26
Content-Type: application/vnd.api+json

{
	"data": {
		"attributes": {
			"priority": 1
		}
	}
}

HTTP/1.1 204
```

<aside class="notice">Previously known as <code>POST /tfa/{id}</code></aside>



Only factors with `priority` greater than 0 will be considered enabled.

## Delete factor 🔒 

```ruby
DELETE /wallets/388108095960430b80554ac3efb6807a9f286854033aca47f6f466094ab50876/factors/24

HTTP/1.1 204
```
`DELETE /wallets/{wallet-id}/factors/{factor}`
