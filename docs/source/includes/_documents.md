# Documents

### Document types

| Type            | Visibility |
| --------------- | ---------- |
| `asset_logo`    | public     |
| `fund_logo`     | public     |
| `fund_document` | public     |
| `nav_report`    | public     |

### Allowed content types

* `application/pdf`
* `image/jpeg`
* `image/tiff`
* `image/png`
* `image/gif`

## Upload 🔒

```http
POST /users/GBT3XFWQUHUTKZMI22TVTWRA7UHV2LIO2BIFNRCH3CXWPYVYPTMXMDGC/documents HTTP/1.1
Content-Type: application/vnd.api+json
Accept: application/vnd.api+json

{
	"data": {
     	"type": "asset_logo",
     	"attributes": {
        	"content_type": "image/png"
     	}
	}
}

HTTP/1.1 200
Content-Type: application/vnd.api+json

{
    "data": {
        "type": "upload_policy",
        "attributes": {
            "bucket": "api",
            "key": "dpurgh4infjubjhcost7fvmaij6y2s3as33sueqfjdmnyrxig3qwm3uc",
            "policy": "eyJleHBpcmF0aW9...ZXF1ZXN0Il1dfQ==",
            "url": "http://localhost:9000/api/",
            "x-amz-algorithm": "AWS4-HMAC-SHA256",
            "x-amz-credential": "2SRMRAST49JEIMUUKWKH/20171225/us-east-1/s3/aws4_request",
            "x-amz-date": "20171225T001154Z",
            "x-amz-signature": "4881f2cada23e19df7f5a92fd1d82fe4c8e4b33ea69e61f05ca2a8860d9a5f3b"
        }
    }
}
```

