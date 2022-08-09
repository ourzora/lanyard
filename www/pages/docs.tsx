import CodeBlock from 'components/CodeBlock'

const createCode = `
POST /api/v1/tree

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

const lookupCode = `
GET /api/v1/tree?root={root}

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

const proofCode = `
GET /api/v1/proof?root={root}&unhashedLeaf={unhashedLeaf}

// Response body
{
  "proof": [ // or empty if the unhashed leaf is not in the merkle tree
    "0x0000000000000000000000000000000000000001",
    "0x0000000000000000000000000000000000000002"
  ],
  "unhashedLeaf": "0x0000000000000000000000000000000000000003"
}
`.trim()

export default function Docs() {
  return (
    <div className="flex flex-col gap-4">
      <h1 className="font-bold text-3xl">API Documentation</h1>

      <Heading>Creating a Merkle tree</Heading>
      <p>
        If you have a list of addresses for an allow list, you can create a
        Merkle tree using this endpoint. This list will automatically be shared
        with primary marketplaces.
      </p>
      <CodeBlock code={createCode} language="javascript" />

      <Heading>Looking up a Merkle tree</Heading>
      <CodeBlock code={lookupCode} language="javascript" />

      <Heading>Getting proof for a value in a Merkle tree</Heading>
      <p>Typically the unhashed leaf value will be an address.</p>
      <CodeBlock code={proofCode} language="javascript" />
    </div>
  )
}

const Heading = ({ children }: { children: React.ReactNode }) => (
  <h1 className="font-bold text-2xl">{children}</h1>
)
