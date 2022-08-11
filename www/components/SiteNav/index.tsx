import Link from 'next/link'
import { useRouter } from 'next/router'
import { useMemo } from 'react'
import { twitterUrl, githubUrl } from 'utils/constants'
import NavTab from './NavTab'

export default function SiteNav() {
  const { pathname } = useRouter()

  const createTabSelectedOverride = useMemo(
    () => (pathname === '/tree/[merkleRoot]' ? true : undefined),
    [pathname],
  )

  return (
    <div className="flex flex-col sm:flex-row items-center justify-between mt-8 mb-16 sm:mb-24 lg:mb-32 gap-4">
      <Link href="/">
        <a className="font-bold text-3xl">Lanyard</a>
      </Link>
      <div className="flex gap-x-6">
        <NavTab
          href="/"
          title="Create"
          selectedOverride={createTabSelectedOverride}
        />
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
