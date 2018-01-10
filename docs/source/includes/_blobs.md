# Blobs

### Blob types

| Type                | Filter |
| ------------------- | ------ |
| `asset_description` | 1      |
| `fund_overview`     | 2      |
| `fund_update`       | 4      |
| `nav_update`        | 8      |


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
## Filter

```http
GET /users/GCCJPB7QQLNEMCJ72CQJ4ODAZFFGXHET5UPOGSDWT222GXQPMTO6ZQW3/blobs HTTP/1.1
Accept: application/vnd.api+json


HTTP/1.1 200
Content-Type: application/vnd.api+json

{
    "data": [
        {
            "id": "ZCPGS6E3MH7I7WNWAHHYW3C6PZI2R4ONUMCEERN35J33MQBAGJ3Q",
            "type": "fund_overview",
            "attributes": {
                "value": "aaaa"
            }
        },
        {
            "id": "57UOND563OVFKI6BDZU7DBPDZFJJORIGFPEMHZQ44IQLWO4RYUOQ",
            "type": "nav_update",
            "attributes": {
                "value": "bbbb"
            }
        }
    ]
}
```

| Parameter | Description            |
| --------- | ---------------------- |
| `type`    | Filter by type bitmask |

Blobs index also supports filtering by user-provided relationships types.