import CodeBlock from 'components/CodeBlock'

const createCode = `
POST /api/v1/tree

// Request body
{
  "unhashedLeaves": [
    "0x0000000000000000000000000000000000000001",
    "0x0000000000000000000000000000000000000002"
  ],
  // abiSig defaults to \`address\`, but you can pass something else
  // if you have different data requirements
  "abiSig": "address,uint256,uint256"
}

// Response body
{
  "merkleRoot": "0x000000000000000000000000000000000000000000000000000000000000000f",
  "abiSig": "address,uint256,uint256"
}
`.trim()

const lookupCode = `
GET /api/v1/tree?root={root}&cursor={cursor}

// Response body
{
  "unhashedLeaves": [
    "0x0000000000000000000000000000000000000001",
    "0x0000000000000000000000000000000000000002"
  ],
  "cursor": "2", // or null if there are no more results
  "totalLeafCount": 400,
  "abiSig": "address" // defaults to address
}
`.trim()

const proofCode = `
GET /api/v1/proof?root={root}&unhashedLeaf={unhashedLeaf}

// Response body
{
  "proof": [ // or empty if the unhashed leaf is not in the merkle tree
    "0x0000000000000000000000000000000000000001",
    "0x0000000000000000000000000000000000000002"
  ]
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
