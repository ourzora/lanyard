# api

## Endpoints

```
POST /api/v1/tree

Request Body:
{
  "allowedAddresses": [
    "0x0000000000000000000000000000000000000001",
    "0x0000000000000000000000000000000000000002"
  ]
}

Response Body:
{
  "merkleRoot": "0x0000000000000000000000000000000000000000000000000000000000000001",
}
```

```
GET /api/v1/tree/?root={root}&cursor={cursor}

Response Body:
{
  "allowedAddresses": [
    "0x0000000000000000000000000000000000000001",
    "0x0000000000000000000000000000000000000002"
  ],
  "cursor": "2", // or null if there are no more results
  "totalAddressCount": 400
}
```

```
GET /api/v1/proof?root={root}&address={address}

Response Body:
{
  "proof": [ // or empty if the address is not in the merkle tree
    "0x0000000000000000000000000000000000000001",
    "0x0000000000000000000000000000000000000002"
  ]
}
```
