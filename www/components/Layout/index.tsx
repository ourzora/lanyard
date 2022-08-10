import classNames from 'classnames'
import SiteNav from 'components/SiteNav'
import MetaTags from 'components/MetaTags'
import { ReactNode } from 'react'
import { siteDescription } from 'utils/constants'

type Props = {
  children: ReactNode
}

export default function Layout({ children }: Props) {
  return (
    <div className="flex flex-col">
      <MetaTags description={siteDescription} />

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
