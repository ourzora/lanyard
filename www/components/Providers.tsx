import { ReactNode } from 'react'
import { SWRConfig } from 'swr'

type Props = {
  children: ReactNode
}

export default function Providers({ children }: Props) {
  return <SWRConfig value={{ shouldRetryOnError: false }}>{children}</SWRConfig>
}
