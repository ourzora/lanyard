import classNames from 'classnames'
import Link from 'next/link'
import { useRouter } from 'next/router'
import { twitterUrl, githubUrl } from 'utils/constants'

export default function SiteNav() {
  return (
    <div className="flex flex-col sm:flex-row items-center justify-between my-8 gap-4">
      <Link href="/">
        <a className="font-bold text-3xl">Lanyard</a>
      </Link>
      <div className="flex gap-x-6">
        <NavTab href="/" title="Create" />
        <NavTab href="/docs" title="API" />
        <a
          href={twitterUrl}
          target="_blank"
          rel="noopener noreferrer"
          className="text-md"
        >
          Twitter
        </a>
        <a
          href={githubUrl}
          target="_blank"
          rel="noopener noreferrer"
          className="text-md"
        >
          Github
        </a>
      </div>
    </div>
  )
}

const NavTab = ({ href, title }: { href: string; title: string }) => {
  const { asPath } = useRouter()
  const isActive = asPath === href
  return (
    <Link href={href}>
      <a
        className={classNames(
          'text-md',
          isActive && 'font-bold border-b-4 border-brand',
        )}
      >
        {title}
      </a>
    </Link>
  )
}
