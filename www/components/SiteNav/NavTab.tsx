import React from 'react'
import Link from 'next/link'
import { useRouter } from 'next/router'
import classNames from 'classnames'
import { brandUnderlineClasses } from 'utils/theme'

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
          (isActive || selectedOverride) && [
            'font-bold',
            brandUnderlineClasses,
          ],
        )}
      >
        {title}
      </a>
    </Link>
  )
}

export default React.memo(NavTab)
