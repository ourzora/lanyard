import CodeBlock from 'components/CodeBlock'
import PageTitle from 'components/PageTitle'
import Tutorial from 'components/Tutorial'
import {
  GetServerSideProps,
  GetServerSidePropsContext,
  InferGetServerSidePropsType,
} from 'next'
import Link from 'next/link'
import { getMerkleTree, TreeResponse } from 'utils/api'
import { dmMintFunTwitterUrl } from 'utils/constants'
import { brandUnderlineClasses } from 'utils/theme'

type Props = {
  merkleRoot: string
  tree: TreeResponse
}

export default function MerkleRootPage({
  merkleRoot,
  tree,
}: InferGetServerSidePropsType<typeof getServerSideProps>) {
  return (
    <div className="flex flex-col gap-y-[6rem]">
      <div className="flex flex-col gap-y-4">
        <PageTitle noPadding>
          Here&rsquo;s the Merkle root for your allowlist
        </PageTitle>
        <div className="font-bold">
          <CodeBlock language="txt" code={merkleRoot} oneLiner />
        </div>
        <div>
          The entire list of addresses in the allowlist can be found on the{' '}
          <Link href={`/membership/${merkleRoot}`} passHref>
            <a className={brandUnderlineClasses}>membership page</a>
          </Link>{' '}
          for your Merkle root.
        </div>
        <div>
          Wire up your Merkle root with the guide below. If you need help,{' '}
          <a
            href={dmMintFunTwitterUrl}
            target="_blank"
            rel="noreferrer noopener"
            className={brandUnderlineClasses}
          >
            DM us on Twitter.
          </a>
        </div>
      </div>

      <Tutorial addresses={tree.unhashedLeaves} />
    </div>
  )
}

export const getServerSideProps: GetServerSideProps<Props> = async (
  ctx: GetServerSidePropsContext,
) => {
  const merkleRoot = String(ctx.params?.merkleRoot)

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
