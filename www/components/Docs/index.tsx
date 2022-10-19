import Link from 'next/link'
import { brandUnderlineClasses } from 'utils/theme'
import PageTitle from 'components/PageTitle'
import CodeBlock from 'components/CodeBlock'
import { createCode, lookupCode, proofCode, rootsCode } from './codeSnippets'
import Section from './Section'

export default function Docs() {
  return (
    <div className="flex flex-col gap-4">
      <div className="flex flex-col mb-5 sm:mb-10">
        <PageTitle noPadding>API Documentation</PageTitle>
        <p className="font-bold mt-6">Client libraries</p>
        <p className="mt-4">
          The easiest way to get started with integrating Lanyard are our client
          libraries. We have a{' '}
          <a
            target="_blank"
            rel="noopener noreferrer"
            href="https://www.npmjs.com/package/lanyard"
            className={brandUnderlineClasses}
          >
            npm package
          </a>{' '}
          (works on both Node.js and the browser) and a{' '}
          <a
            target="_blank"
            rel="noopener noreferrer"
            href="https://pkg.go.dev/github.com/contextwtf/lanyard/api"
            className={brandUnderlineClasses}
          >
            Go client
          </a>
          .
        </p>
        <p className="font-bold mt-4">Solidity</p>

        <p className="mt-4 mb-10">
          Looking for an example on how to implement Merkle roots in your
          contract?{' '}
          <Link href="/tree/sample" passHref>
            <a className={brandUnderlineClasses}>
              Here’s a guide with our sample root.
            </a>
          </Link>
        </p>
        <div className="h-px bg-neutral-200 w-full" />
      </div>

      <div className="flex flex-col gap-16">
        <Section
          title="Creating a Merkle tree"
          description={`If you have a list of addresses for an allowlist, you can create a Merkle tree using this endpoint.  Any Merkle tree published on Lanyard will be publicly available to any user of Lanyard’s API, including minting interfaces such as Zora or mint.fun.`}
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
        <Section
          title="Looking up potential roots for a given proof"
          description="A more advanced use but helpful if you just have a proof."
        >
          <CodeBlock code={rootsCode} language="javascript" />
        </Section>
      </div>
    </div>
  )
}
