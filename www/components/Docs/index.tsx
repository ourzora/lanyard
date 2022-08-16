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
          description={`If you have a list of addresses for an allow list, you can create a Merkle tree using this endpoint.  Any Merkle tree published on Lanyard will be publicly available to any user of the Lanyard’s API, including minting interfaces such as Zora or mint.fun.`}
        >
          <CodeBlock code={createCode} language="javascript" />
        </Section>

        <Section
          title="Looking up a Merkle tree"
          description="If a Merkle tree has been published to Lanyard, it’s possible to request the entire tree by providing the root. This endpoint will 404 if the tree associated with the root has not been published."
        >
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
