import classNames from 'classnames'
import React from 'react'

type Props = {
  children: React.ReactNode
  noPadding?: boolean
}

export default function PageTitle({ children, noPadding = false }: Props) {
  return (
    <div
      className={classNames(
        'font-bold text-2xl sm:text-3xl lg:text-4xl',
        'text-center sm:text-left',
        !noPadding && 'mb-10',
      )}
    >
      {children}
    </div>
  )
}
