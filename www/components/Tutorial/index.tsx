import classNames from 'classnames'
import React, { useCallback } from 'react'
import CodeBlock from '../CodeBlock'
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
    <div className={classNames('flex flex-col gap-y-4', 'w-full')}>
      <h1 className="font-bold text-3xl">How to use your Merkle root</h1>

      <Heading>1. Install dependencies</Heading>
      <p>
        Install the following libraries to generate the proof for a wallet on
        your site.
      </p>
      <CodeBlock
        title="Terminal"
        code={installDependenciesCode}
        language="txt"
      />

      <Heading>2. Add Merkle tree code</Heading>
      <p>
        Use the following code, replacing the addresses with the addresses
        pasted here, to have the code ready to generate your Merkle proofs.
      </p>
      <CodeBlock
        title="merkle.ts"
        code={merkleSetupCode([])}
        codeForCopy={merkleSetupCodeForCopy}
        language="typescript"
      />

      <Heading>3. Pass Merkle proof to your contract</Heading>
      <p>
        Using the file above, import the `getMerkleProof` function to generate a
        proof for a connected wallet.
      </p>
      <CodeBlock
        title="mintpage.tsx"
        code={passMerkleProofCode}
        language="tsx"
      />

      <Heading>4. Check the Merkle proof in your contract</Heading>
      <p>
        When you have a Merkle root, you can check proofs using this helper from
        OpenZeppelin&apos;s contracts.
      </p>
      <CodeBlock
        title="NFTContract.sol"
        code={nftMerkleProofCode}
        language="sol"
      />
    </div>
  )
}

const Heading = ({ children }: { children: React.ReactNode }) => (
  <h1 className="font-bold text-2xl">{children}</h1>
)

export default React.memo(Tutorial)
