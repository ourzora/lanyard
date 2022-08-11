import CodeBlock from 'components/CodeBlock'
import Tutorial from 'components/Tutorial'
import {
  GetServerSideProps,
  GetServerSidePropsContext,
  InferGetServerSidePropsType,
} from 'next'
import { getMerkleTree, TreeResponse } from 'utils/api'

type Props = {
  merkleRoot: string
  tree: TreeResponse
}

export default function MerkleRootPage({
  merkleRoot,
  tree,
}: InferGetServerSidePropsType<typeof getServerSideProps>) {
  return (
    <div className="flex flex-col gap-y-10">
      <div className="flex flex-col gap-y-4">
        <div className="font-bold text-2xl">Here&rsquo;s your Merkle root!</div>
        <CodeBlock language="txt" code={merkleRoot} title="Merkle root" />
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
