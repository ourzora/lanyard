import PageTitle from 'components/PageTitle'
import CodeBlock from 'components/CodeBlock'
import { createCode, lookupCode, proofCode } from './codeSnippets'

export default function Docs() {
  return (
    <div className="flex flex-col gap-4">
      <PageTitle>API Documentation</PageTitle>

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
