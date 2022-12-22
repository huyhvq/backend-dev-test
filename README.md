## Getting started

Make sure that you're in the root of the project directory, fetch the dependencies with `go mod tidy`, then run the
application using `go run main.go serve`:

```
$ go mod tidy
$ go run main.go serve
```

## Project structure

Everything in the codebase is designed to be editable. Feel free to change and adapt it to meet your needs.

|     |     |
| --- | --- |
| **`assets`** | Contains the non-code assets for the application. |
| `↳ assets/migrations/` | Contains SQL migrations. |
| `↳ assets/efs.go` | Declares an embedded filesystem containing all the assets. |

|     |     |
| --- | --- |
| **`cmd/api`** | Application-specific code (handlers, routing, middleware) for dealing with HTTP requests and responses. |

|                             |                                                                                                                                  |
|-----------------------------|----------------------------------------------------------------------------------------------------------------------------------|
| **`internal`**              | Contains private application and library code. |
| `↳ internal/database/`      | Contains database-related code. |
| `↳ internal/server/`        | Contains a helper function for starting and gracefully shutting down the server. |
| `↳ internal/organizations/` | Contains the ecommerce domain specific code. |

|                          |     |
|--------------------------| --- |
| **`pkg`**                | Contains various helper packages used by the application. |
| `↳ internal/leveledlog/` | Contains a leveled logger implementation. |
| `↳ internal/request/`    | Contains helper functions for decoding JSON requests. |
| `↳ internal/response/`   | Contains helper functions for sending JSON responses. |
| `↳ internal/validator/`  | Contains validation helpers. |
| `↳ internal/version/`    | Contains the application version number definition. |

## Configuration settings

Configuration settings are managed via yaml file or os environment.
os enviroment > config yaml > default config in `cmd/api/config.go`

## How to start at local environment
Make sure that you're in the root of the project directory, start docker instance using:
`docker-compose up --build`

If you make a request to the `GET /_/status` endpoint using `curl` you should get a response like this:
```
$ curl -i localhost:3000/_/status
HTTP/1.1 200 OK
Content-Type: application/json
Date: Tue, 29 Nov 2022 17:46:40 GMT
Content-Length: 20

{
	"Status": "OK"
}
```
Application APIs endpoints:
##### Create a hub :
POST `/v1/hubs`
```
$ curl -i --location --request POST 'http://localhost:3000/v1/hubs' \
    --header 'Content-Type: application/json' \
    --data-raw '{
    "name": "Test User",
    "location": "Ho Chi Minh"
}'
HTTP/1.1 200 OK
Content-Type: application/json
Date: Thu, 22 Dec 2022 18:12:24 GMT
Content-Length: 142

{
	"id": 11,
	"name": "Test Hub",
	"location": "Ho Chi Minh",
	"created_at": "0001-01-01T00:00:00Z",
	"updated_at": "0001-01-01T00:00:00Z"
}
```
##### Create a team :
POST `/v1/teams`
```
$ curl -i --location --request POST 'http://localhost:3000/v1/teams' \
    --header 'Content-Type: application/json' \
    --data-raw '{
    "name": "Team 1",
    "type": "backend",
    "hub_id": 1
}'
HTTP/1.1 200 OK
Content-Type: application/json
Date: Thu, 22 Dec 2022 18:14:28 GMT
Content-Length: 145

{
	"id": 10,
	"name": "Team 1",
	"type": "backend",
	"hub_id": 1,
	"created_at": "0001-01-01T00:00:00Z",
	"updated_at": "0001-01-01T00:00:00Z"
}
```

##### Create a user :
POST `/v1/users`
```
$ curl -i --location --request POST 'http://localhost:3000/v1/users' \
    --header 'Content-Type: application/json' \
    --data-raw '{
    "name": "John Doe",
    "title": "SRE",
    "team_id": 1
}'
HTTP/1.1 200 OK
Content-Type: application/json
Date: Thu, 22 Dec 2022 18:15:53 GMT
Content-Length: 143

{
	"id": 2,
	"name": "John Doe",
	"type": "SRE",
	"team_id": 1,
	"created_at": "0001-01-01T00:00:00Z",
	"updated_at": "0001-01-01T00:00:00Z"
}
```
##### get hubs & teams with search term:
GET `/v1/organizations/search?term=Ho`
For example: I want to get product with: search term is `Ho` with curl:
```
$ curl -i --location --request GET 'http://localhost:3000/v1/organizations/search?term=Ho'
HTTP/1.1 200 OK
Content-Type: application/json
Date: Thu, 22 Dec 2022 18:17:49 GMT
Content-Length: 1054

{
	"hubs": [
		{
			"id": 2,
			"name": "Shopee non tech",
			"location": "Ho Chi Minh",
			"teams": [
				{
					"id": 2,
					"name": "Manual Quality Assurance",
					"type": "qa"
				},
				{
					"id": 3,
					"name": "Devops",
					"type": "sre"
				},
				{
					"id": 7,
					"name": "Platform",
					"type": "backend"
				}
			]
		},
		{
			"id": 3,
			"name": "Lazada tech hub",
			"location": "Ho Chi Minh",
			"teams": [
				{
					"id": 4,
					"name": "ABC Swad",
					"type": "backend"
				},
				{
					"id": 5,
					"name": "Risk Management",
					"type": "nontech"
				}
			]
		},
		{
			"id": 11,
			"name": "Test User",
			"location": "Ho Chi Minh",
			"teams": []
		},
		{
			"id": 10,
			"name": "Test User",
			"location": "Ho Chi Minh",
			"teams": []
		},
		{
			"id": 8,
			"name": "Test User",
			"location": "Ho Chi Minh",
			"teams": []
		},
		{
			"id": 9,
			"name": "Test User",
			"location": "Ho Chi Minh",
			"teams": []
		},
		{
			"id": 7,
			"name": "Test User",
			"location": "Ho Chi Minh",
			"teams": []
		}
	]
}
```
##### user join to a team:
PATCH `/v1/users/1/join-team`
```
$ curl -i --location --request PATCH 'http://localhost:3000/v1/users/1/join-team' \
    --header 'Content-Type: application/json' \
    --data-raw '{
    "team_id": 2
}'
HTTP/1.1 200 OK
Content-Type: application/json
Date: Thu, 22 Dec 2022 18:20:57 GMT
Content-Length: 157

{
	"id": 2,
	"name": "Manual Quality Assurance",
	"type": "qa",
	"hub_id": 2,
	"created_at": "2022-12-17T14:16:14Z",
	"updated_at": "0001-01-01T00:00:00Z"
}
```

##### team join to a hub:
PATCH `/v1/teams/1/join-hub`
```
$ curl -i --location --request PATCH 'http://localhost:3000/v1/teams/2/join-hub' \
    --header 'Content-Type: application/json' \
    --data-raw '{
    "hub_id": 2
}'
HTTP/1.1 200 OK
Content-Type: application/json
Date: Thu, 22 Dec 2022 18:21:54 GMT
Content-Length: 157

{
	"id": 2,
	"name": "Manual Quality Assurance",
	"type": "qa",
	"hub_id": 2,
	"created_at": "2022-12-17T14:16:14Z",
	"updated_at": "0001-01-01T00:00:00Z"
}
```