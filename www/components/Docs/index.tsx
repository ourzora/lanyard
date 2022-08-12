import PageTitle from 'components/PageTitle'
import CodeBlock from 'components/CodeBlock'
import { createCode, lookupCode, proofCode } from './codeSnippets'
import Section from './Section'

export default function Docs() {
  return (
    <div className="flex flex-col gap-4">
      <div className="flex flex-col mb-5 sm:mb-10">
        <PageTitle>API Documentation</PageTitle>
        <div className="h-px bg-neutral-200 w-full" />
      </div>

      <div className="flex flex-col gap-16">
        <Section
          title="Creating a Merkle tree"
          description={`If you have a list of addresses for an allow list, you can create a Merkle tree using this endpoint. This list will automatically be shared with primary marketplaces.`}
        >
          <CodeBlock code={createCode} language="javascript" />
        </Section>

        <Section title="Looking up a Merkle tree">
          <CodeBlock code={lookupCode} language="javascript" />
        </Section>
        <Section
          title="Getting proof for a value in a Merkle tree"
          description="Typically the unhashed leaf value will be an address."
        >
          <CodeBlock code={proofCode} language="javascript" />
        </Section>
      </div>
    </div>
  )
}
