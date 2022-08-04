import { GetServerSidePropsContext, NextPageContext } from 'next'
import { clientRetry } from './client'

export const apiURL = (): string => {
  // use the API_URL if we're on the backend
  // use the next-prefixed url's if available otherwise
  const apiUrlBackend = process.env.API_URL

  if (apiUrlBackend !== undefined) {
    return `${apiUrlBackend}/api`
  }

  if (typeof window === 'undefined') {
    throw new Error('missing API_URL env on backend')
  }

  return `/api`
}

export async function fetcher(
  url: string,
  ctx?: GetServerSidePropsContext | NextPageContext,
) {
  const r = await clientRetry('GET', url, null, ctx)
  if (!r.ok) {
    throw new Error(`Failed Response: status=${r.status}`)
  }
  return r.json()
}
