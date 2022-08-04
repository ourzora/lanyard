import classNames from 'classnames'
import Link from 'next/link'
import { useRouter } from 'next/router'

export default function SiteNav() {
  return (
    <div className="flex mt-4 mb-4 flex-col gap-y-4">
      <Link href="/">
        <a className="font-bold text-3xl">allowlist</a>
      </Link>
      <div className="flex gap-x-4">
        <NavTab href="/" title="Creators" />
        <NavTab href="/docs" title="API Documentation" />{' '}
      </div>
    </div>
  )
}

const NavTab = ({ href, title }: { href: string; title: string }) => {
  const { asPath } = useRouter()
  const isActive = asPath === href
  return (
    <Link href={href}>
      <a className={classNames(isActive && 'font-semibold')}>{title}</a>
    </Link>
  )
}
