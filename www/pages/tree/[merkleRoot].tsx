import CodeBlock from 'components/CodeBlock'
import PageTitle from 'components/PageTitle'
import Tutorial from 'components/Tutorial'
import {
  GetServerSideProps,
  GetServerSidePropsContext,
  InferGetServerSidePropsType,
} from 'next'
import Link from 'next/link'
import { getMerkleTree, useMerkleTree } from 'utils/api'
import { dmMintFunTwitterUrl } from 'utils/constants'
import { brandUnderlineClasses } from 'utils/theme'

type Props = {
  merkleRoot: string
}

export default function MerkleRootPage({
  merkleRoot,
}: InferGetServerSidePropsType<typeof getServerSideProps>) {
  const { data } = useMerkleTree(merkleRoot)

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

      {data !== undefined && <Tutorial addresses={data.unhashedLeaves} />}
    </div>
  )
}

export const getServerSideProps: GetServerSideProps<Props> = async (
  ctx: GetServerSidePropsContext,
) => {
  const merkleRoot = String(ctx.params?.merkleRoot)

  try {
    // will fail if tree is missing
    await getMerkleTree(merkleRoot)
    return {
      props: {
        merkleRoot,
      },
    }
  } catch {
    return {
      notFound: true,
    }
  }
}
