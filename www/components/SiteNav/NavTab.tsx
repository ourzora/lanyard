import React from 'react'
import Link from 'next/link'
import { useRouter } from 'next/router'
import classNames from 'classnames'

type Props = {
  href: string
  title: string
  selectedOverride?: boolean
}

function NavTab({ href, title, selectedOverride = false }: Props) {
  const { asPath } = useRouter()
  const isActive = asPath === href
  return (
    <Link href={href}>
      <a
        className={classNames(
          'text-md',
          (isActive || selectedOverride) && 'font-bold border-b-4 border-brand',
        )}
      >
        {title}
      </a>
    </Link>
  )
}

export default React.memo(NavTab)
