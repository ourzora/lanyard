import classNames from 'classnames'
import SiteNav from 'components/SiteNav'
import { ReactNode } from 'react'

type Props = {
  children: ReactNode
}

export default function Layout({ children }: Props) {
  return (
    <div className="flex flex-col">
      <div
        className={classNames(
          'overflow-x-hidden w-full max-w-screen-xl',
          'px-3 md:px-8 pb-8',
          'self-center',
        )}
      >
        <SiteNav />
        <main>{children}</main>
      </div>
    </div>
  )
}
