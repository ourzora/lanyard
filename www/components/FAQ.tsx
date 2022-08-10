import Link from 'next/link'
import { twitterUrl } from 'utils/constants'

export default function FAQ() {
  return (
    <div className="flex flex-col gap-y-12">
      <QandA>
        <Question>What is Lanyard?</Question>
        <Answer>
          Lanyard is a free and open source tool for quickly and easily
          generating a Merkle root for your NFT project allow list.
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
            className="underline"
          >
            Merkle root
          </a>{' '}
          is the most popular way to handle allow lists for projects using the
          blockchain. By using a Merkle root, you and your community pay
          significantly less in gas fees to prove a wallet is allowed to mint a
          specific project. This is done by proving an address belongs to a tree
          of hashed values leading to the root.
        </Answer>
      </QandA>
      <QandA>
        <Question>
          What do you mean when you say my allow list will work across web3?
        </Question>
        <Answer>
          Lanyard securely stores your allow list and Merkle root so that your
          community can mint your project from their preferred interface – your
          website, mint.fun, and any other platform that integrates Lanyard.
        </Answer>
      </QandA>
      <QandA>
        <Question>
          I’m building a minting platform – how do I integrate Lanyard?
        </Question>
        <Answer>
          See our{' '}
          <Link href="/docs">
            <a className="underline">API documentation here</a>
          </Link>
          .
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
            className="underline"
          >
            DM us on Twitter
          </a>{' '}
          if you want to contribute.
        </Answer>
      </QandA>
    </div>
  )
}

const Question = ({ children }: { children: React.ReactNode }) => (
  <div className="font-bold text-xl">{children}</div>
)

const Answer = ({ children }: { children: React.ReactNode }) => (
  <div className="text-xl">{children}</div>
)

const QandA = ({ children }: { children: React.ReactNode }) => (
  <div className="flex flex-col gap-y-2">{children}</div>
)
