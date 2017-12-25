# Blobs

### Blob types

* `asset_description`



## Create

```http
POST /users/GCCJPB7QQLNEMCJ72CQJ4ODAZFFGXHET5UPOGSDWT222GXQPMTO6ZQW3/blobs HTTP/1.1
Content-Type: application/vnd.api+json
Accept: application/vnd.api+json

{
	"data": {
		"type": "asset_description",
		"attributes": {
			"value": "foobar"
		}
	}
}


HTTP/1.1 201
Content-Type: application/vnd.api+json

{
    "data": {
        "id": "CMKFL7K43ZG6E4S3QFEYFYM2F7GS2QJHQPLI25PIWFYMT5V4ZOKQ",
        "type": "asset_description",
        "attributes": {
            "value": "foobar"
        }
    }
}
```

## Get

```http
GET /users/GCCJPB7QQLNEMCJ72CQJ4ODAZFFGXHET5UPOGSDWT222GXQPMTO6ZQW3/blobs/CMKFL7K43ZG6E4S3QFEYFYM2F7GS2QJHQPLI25PIWFYMT5V4ZOKQ HTTP/1.1
Accept: application/vnd.api+json


HTTP/1.1 200
Content-Type: application/vnd.api+json

{
    "data": {
        "id": "CMKFL7K43ZG6E4S3QFEYFYM2F7GS2QJHQPLI25PIWFYMT5V4ZOKQ",
        "type": "asset_description",
        "attributes": {
            "value": "foobar"
        }
    }
}
```

