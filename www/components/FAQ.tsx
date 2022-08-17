import Link from 'next/link'
import { twitterUrl } from 'utils/constants'
import { brandUnderlineClasses } from 'utils/theme'

export default function FAQ() {
  return (
    <div className="flex flex-col gap-y-12">
      <QandA>
        <Question>What is Lanyard?</Question>
        <Answer>
          <ol className="list-decimal pl-5">
            <li className="pl-1">
              A Merkle root generator: just paste the addresses you want to
              allow, and we generate and return your root instantly
            </li>
            <li className="pl-1">
              A{' '}
              <Link href="/tree/sample" passHref>
                <a className={brandUnderlineClasses}>step-by-step guide</a>
              </Link>{' '}
              on how to properly implement the Merkle root in your project –
              with free support via Discord if you need it
            </li>
            <li className="pl-1">
              A mechanism to ensure your allowlist works across web3
            </li>
          </ol>
        </Answer>
      </QandA>
      <QandA>
        <Question>What is a Merkle root?</Question>
        <Answer>
          A{' '}
          <a
            href="https://en.wikipedia.org/wiki/Merkle_root"
            target="_blank"
            rel="noopener noreferrer"
            className={brandUnderlineClasses}
          >
            Merkle root
          </a>{' '}
          is the most popular way to handle allowlists for projects using the
          blockchain. By using a Merkle root, you and your community pay
          significantly less in gas fees to prove a wallet is allowed to mint a
          specific project. This is done by proving an address belongs to a tree
          of hashed values leading to the root.
        </Answer>
      </QandA>
      <QandA>
        <Question>
          What do you mean when you say my allowlist will work across web3?
        </Question>
        <Answer>
          Lanyard securely stores your allowlist, Merkle root, and proofs so
          that your community can mint your project from their preferred
          interface – your website,{' '}
          <a
            href="https://mint.fun"
            target="_blank"
            rel="noopener noreferrer"
            className={brandUnderlineClasses}
          >
            mint.fun
          </a>
          , and any other platform that integrates Lanyard.
        </Answer>
      </QandA>
      <QandA>
        <Question>I build creator tools – how do I integrate Lanyard?</Question>
        <Answer>
          See our API documentation{' '}
          <Link href="/docs">
            <a className={brandUnderlineClasses}>here</a>
          </Link>{' '}
          so your creators’ allowlists work with Lanyard.
        </Answer>
      </QandA>
      <QandA>
        <Question>Who’s behind Lanyard?</Question>
        <Answer>
          Lanyard is an open source collaboration between Context, Zora, and
          others.{' '}
          <a
            href={twitterUrl}
            target="_blank"
            rel="noopener noreferrer"
            className={brandUnderlineClasses}
          >
            DM us on Twitter
          </a>{' '}
          if you want to contribute or want your integration listed on this
          page.
        </Answer>
      </QandA>
    </div>
  )
}

const Question = ({ children }: { children: React.ReactNode }) => (
  <div className="font-bold text-2xl">{children}</div>
)

const Answer = ({ children }: { children: React.ReactNode }) => (
  <div className="text-lg">{children}</div>
)

const QandA = ({ children }: { children: React.ReactNode }) => (
  <div className="flex flex-col gap-y-2">{children}</div>
)
