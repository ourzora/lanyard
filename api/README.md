# api

## Endpoints

```
POST /api/v1/tree

Request Body:
{
    "unhashedLeaves": [
        "0x0000000000000000000000000000000000000001",
        "0x0000000000000000000000000000000000000002"
    ],
    "leafTypeDescriptor": "address",
    "packedEncoding": true
}

Response Body:
{
  "merkleRoot": "0x0000000000000000000000000000000000000000000000000000000000000001",
}
```

```
GET /api/v1/tree?root={root}

Response Body:
{
  "unhashedLeaves": [
    "0x0000000000000000000000000000000000000001",
    "0x0000000000000000000000000000000000000002"
  ],
  "leafCount": 2
}
```

```
GET /api/v1/proof?root={root}&unhashedLeaf={unhashedLeaf}

Response Body:
{
  "proof": [ // or empty if the address is not in the merkle tree
    "0x0000000000000000000000000000000000000001",
    "0x0000000000000000000000000000000000000002"
  ],
  "unhashedLeaf": "0x0000000000000000000000000000000000000003" // or null if not in the tree
}
```
