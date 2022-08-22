export const createCode = `
POST https://lanyard.build/api/v1/tree

// Request body
{
  "unhashedLeaves": [
    "0x0000000000000000000000000000000000000001",
    "0x0000000000000000000000000000000000000002"
  ],

  // in general you can omit the two following fields, but if you have specific data
  // requirements you can include them to help with indexing
  "leafTypeDescriptor": ["address"], // defaults to \`["address"]\`, can pass other solidity types
  "packedEncoding": true // defaults to \`true\`
}

// Response body
{
  "merkleRoot": "0x000000000000000000000000000000000000000000000000000000000000000f"
}
`.trim()

export const lookupCode = `
GET https://lanyard.build/api/v1/tree?root={root}

// Response body
{
  "unhashedLeaves": [
    "0x0000000000000000000000000000000000000001",
    "0x0000000000000000000000000000000000000002"
  ],
  "leafCount": 2,

  // in general you can ignore the two following fields
  "leafTypeDescriptor": null, // or an array of solidity types
  "packedEncoding": null // or a boolean value
}
`.trim()

export const proofCode = `
GET https://lanyard.build/api/v1/proof?root={root}&unhashedLeaf={unhashedLeaf}

// Response body
{
  "proof": [ // or empty if the unhashed leaf is not in the Merkle tree
    "0x0000000000000000000000000000000000000001",
    "0x0000000000000000000000000000000000000002"
  ],
  "unhashedLeaf": "0x0000000000000000000000000000000000000003" // or null if not in the tree
}
`.trim()

export const rootCode = `
// proof is 0x prefixed, comma separated values
GET https://lanyard.build/api/v1/root?proof={proof}

// Response body
{
  "root": "0x0000000000000000000000000000000000000003" // returns error if not found
}
`.trim()
