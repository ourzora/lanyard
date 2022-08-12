import React, { useCallback } from 'react'
import CodeBlock from 'components/CodeBlock'
import Section from './Section'
import {
  installDependenciesCode,
  merkleSetupCode,
  nftMerkleProofCode,
  passMerkleProofCode,
} from './codeSnippets'

type Props = {
  addresses: string[]
}

function Tutorial({ addresses }: Props) {
  const merkleSetupCodeForCopy = useCallback(
    () => merkleSetupCode(addresses),
    [addresses],
  )

  return (
    <div className="flex flex-col gap-y-8 w-full">
      <h1 className="font-bold text-2xl sm:text-3xl">
        How to use the Merkle root in your contract
      </h1>
      <div className="flex flex-col gap-y-[6rem] w-full">
        <Section
          title="1. Install dependencies"
          description="Install the following libraries to generate the proof for a wallet on your site."
        >
          <CodeBlock
            title="Terminal"
            code={installDependenciesCode}
            language="txt"
          />
        </Section>

        <Section
          title="2. Add Merkle tree code"
          description="Use the following code to generate Merkle proofs for your root."
        >
          <CodeBlock
            title="merkle.ts"
            code={merkleSetupCode([])}
            codeForCopy={merkleSetupCodeForCopy}
            language="typescript"
          />
        </Section>

        <Section
          title="3. Pass Merkle proof to your contract"
          description={
            <>
              Using the file above, import the{' '}
              <code className="bg-neutral-100 font-mono px-1 py-1 rounded">
                getMerkleProof
              </code>{' '}
              function to generate a proof for a connected wallet.
            </>
          }
        >
          <CodeBlock
            title="mintpage.tsx"
            code={passMerkleProofCode}
            language="tsx"
          />
        </Section>

        <Section
          title="4. Check the Merkle proof in your contract"
          description={`With your Merkle root, you can check proofs using this helper from OpenZeppelin's contracts.`}
        >
          <CodeBlock
            title="NFTContract.sol"
            code={nftMerkleProofCode}
            language="sol"
          />
        </Section>
      </div>
    </div>
  )
}

export default React.memo(Tutorial)
