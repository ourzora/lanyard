import CodeBlock from 'components/CodeBlock'

const createCode = `
POST /api/v1/merkle

// Request body
{
  "allowedAddresses": [
    "0x0000000000000000000000000000000000000001",
    "0x0000000000000000000000000000000000000002"
  ]
}

// Response body
{
  "merkleRoot": "0x000000000000000000000000000000000000000000000000000000000000000f"
}
`.trim()

const lookupCode = `
GET /api/v1/merkle/{merkleRoot}?cursor={cursor}

// Response body
{
  "allowedAddresses": [
    "0x0000000000000000000000000000000000000001",
    "0x0000000000000000000000000000000000000002"
  ],
  "cursor": "2", // or null if there are no more results
  "totalAddressCount": 400
}
`.trim()

const proofCode = `
GET /api/v1/merkle/{merkleRoot}/proof/{address}

// Response body
{
  "proof": [ // or empty if the address is not in the merkle tree
    "0x0000000000000000000000000000000000000001",
    "0x0000000000000000000000000000000000000002"
  ]
}
`.trim()

export default function PrimaryMarketplacesPage() {
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

      <Heading>Getting proof for an address in a Merkle tree</Heading>
      <CodeBlock code={proofCode} language="javascript" />
    </div>
  )
}

const Heading = ({ children }: { children: React.ReactNode }) => (
  <h1 className="font-bold text-2xl">{children}</h1>
)
