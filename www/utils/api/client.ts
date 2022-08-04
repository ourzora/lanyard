import { GetServerSidePropsContext, NextPageContext } from 'next'
import pRetry from 'p-retry'
import { apiURL } from './helpers'

export async function client(
  method: string,
  url: string,
  body?: string | object | undefined | null,
  ctx?: GetServerSidePropsContext | NextPageContext,
) {
  const headers: {
    ['X-Forwarded-For']?: string
    Cookie?: string
    ['Content-Type']?: string
  } = {}

  if (ctx) {
    if (ctx.req?.headers.forwarded !== undefined) {
      headers['X-Forwarded-For'] = ctx.req.headers.forwarded
    }
    if (ctx.req?.headers.cookie !== undefined) {
      headers['Cookie'] = ctx.req.headers.cookie
    }
  }

  let jsonBody: string | undefined
  if (body !== null && body !== undefined) {
    jsonBody = JSON.stringify(body)
    headers['Content-Type'] = 'application/json'
  }

  const response = await fetch(`${apiURL()}/${encodeURI(url)}`, {
    headers: {
      ...headers,
    },
    method,
    body: jsonBody,
    credentials: 'same-origin',
  })

  if (response.status >= 500 && response.status <= 599) {
    const error = new Error(`Request Failed: status ${response.status} `)
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-expect-error
    error.url = url
    throw error
  }

  return response
}

export function clientRetry(
  method: string,
  url: string,
  body?: string | object | undefined | null,
  ctx?: GetServerSidePropsContext | NextPageContext,
) {
  return pRetry(() => client(method, url, body, ctx), { retries: 5 })
}
