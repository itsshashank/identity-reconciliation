# identity-reconciliation

## How to Run in Locally

create a .env file with following secrets
```
DB_USER=username
DB_PASSWORD=somestrongpassword
DB_NAME=dbname
```

Docker compose up will bring both app and db 

Run locally
```bash
docker-compose up
```

Can access on localhost:8000 or at containerip:8080
```bash
docker inspect 480cf1308d0d | grep "IPAddress"
```
```bash
➜  identity-reconciliation git:(main) ✗ docker inspect c2499030d285 | grep "IPAddress"
            "SecondaryIPAddresses": null,
            "IPAddress": "172.17.0.3",
                    "IPAddress": "172.17.0.3",

```


## How to check app

Create orders:

```bash
curl -X POST http://localhost:8080/order \
  -H 'Content-Type: application/json' \
  -d '{
        "email": "lorraine@hillvalley.edu",
	    "phoneNumber": "123456"
    }'
```

The response:

```json
{
  "message": "Order created successfully"
}
```

Get identity reconciled

```bash
curl -X GET localhost:8080/identify
```

The response:
```json
{
  "contact": {
    "primaryContatctId": 11,
    "emails": [
      "lorraine@hillvalley.edu",
      "mcfly@hillvalley.edu"
    ],
    "phoneNumbers": [
      "123456"
    ],
    "secondaryContactIds": [
      12
    ]
  }
}
```