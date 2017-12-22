# KYC

## Entity types

### Individual

```json
{
  "first_name": "John",
  "last_name": "Doe"
}
```

Personal details for individual account type

#### Required fields

* `first_name` - string
* `last_name` - string



## Create entity

```http
POST /users/GCCJPB7QQLNEMCJ72CQJ4ODAZFFGXHET5UPOGSDWT222GXQPMTO6ZQW3/entity HTTP/1.1
Content-Type: application/vnd.api+json

{
	"data": {
		"type": "individual",
		"attributes": {
			"first_name": "John",
			"last_name": "Doe"
		}
	}
}

HTTP/1.1 204
```



## Get entities

```http
GET /users/GCCJPB7QQLNEMCJ72CQJ4ODAZFFGXHET5UPOGSDWT222GXQPMTO6ZQW3/entities HTTP/1.1
Accept: application/vnd.api+json
Content-Type: application/vnd.api+json


HTTP/1.1 200
Content-Type: application/vnd.api+json

{
    "data": [
        {
            "id": "01C1Z041SK6QB59B7CX5RS8RGR",
            "type": "individual",
            "attributes": {
                "first_name": "John",
                "last_name": "Doe"
            }
        }
    ]
}
```



## Update entity

```http
PUT /users/GCCJPB7QQLNEMCJ72CQJ4ODAZFFGXHET5UPOGSDWT222GXQPMTO6ZQW3/entities/01C1Z041SK6QB59B7CX5RS8RGR HTTP/1.1
Content-Type: application/vnd.api+json

{
	"data": {
		"type": "individual",
		"attributes": {
			"first_name": "Samuel",
			"last_name": "Jackson"
		}
	}
}

HTTP/1.1 204

```

KYC entity does not support `PATCH` semantics, each request replaces attributes completely.