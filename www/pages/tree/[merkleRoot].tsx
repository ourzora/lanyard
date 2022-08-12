import CodeBlock from 'components/CodeBlock'
import PageTitle from 'components/PageTitle'
import Tutorial from 'components/Tutorial'
import {
  GetServerSideProps,
  GetServerSidePropsContext,
  InferGetServerSidePropsType,
} from 'next'
import { getMerkleTree, TreeResponse } from 'utils/api'
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
          Here&rsquo;s the Merkle root for your allow list
        </PageTitle>
        <div className="font-bold">
          <CodeBlock language="txt" code={merkleRoot} oneLiner />
        </div>
        <div>
          Need help using your Merkle root?{' '}
          <a
            href="https://discord.gg/context"
            target="_blank"
            rel="noreferrer noopener"
            className={brandUnderlineClasses}
          >
            Message us on Discord
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
