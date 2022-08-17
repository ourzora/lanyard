import CodeBlock from 'components/CodeBlock'
import PageTitle from 'components/PageTitle'
import Tutorial from 'components/Tutorial'
import { GetServerSideProps, InferGetServerSidePropsType } from 'next'
import Link from 'next/link'
import { getMerkleTree, TreeResponse } from 'utils/api'
import { brandUnderlineClasses } from 'utils/theme'

type Props = {
  merkleRoot: string
  tree: TreeResponse
}

const SAMPLE_ROOT =
  '0x9bcb34c8aba34a442d549dc3ae29995d5d1646440b80329ba55f6978a5bf23ce'

export default function MerkleRootPage({
  merkleRoot,
  tree,
}: InferGetServerSidePropsType<typeof getServerSideProps>) {
  return (
    <div className="flex flex-col gap-y-[4rem]">
      <div className="bg-brand-light outline outline-brand outline-2 rounded-xl p-4">
        <p>
          This is a sample root. Want to generate your own list? Go to the{' '}
          <Link href="/" passHref>
            <a className="underline">homepage</a>
          </Link>{' '}
          or use our{' '}
          <Link href="/docs" passHref>
            <a className="underline">API</a>
          </Link>
          .
        </p>
      </div>
      <div className="flex flex-col gap-y-4">
        <PageTitle noPadding>
          Here&rsquo;s the Merkle root for your allowlist
        </PageTitle>
        <div className="font-bold">
          <CodeBlock language="txt" code={merkleRoot} oneLiner />
        </div>
        <div>
          Wire up your Merkle root with the guide below. If you need help,{' '}
          <a
            href="https://discord.gg/context"
            target="_blank"
            rel="noreferrer noopener"
            className={brandUnderlineClasses}
          >
            message us on Discord.
          </a>
        </div>
      </div>

      <Tutorial addresses={tree.unhashedLeaves} />
    </div>
  )
}

export const getServerSideProps: GetServerSideProps<Props> = async () => {
  const merkleRoot = SAMPLE_ROOT

  try {
    const tree = await getMerkleTree(merkleRoot)
    return {
      props: {
        merkleRoot,
        tree,
      },
    }
  } catch {
    return {
      notFound: true,
    }
  }
}
