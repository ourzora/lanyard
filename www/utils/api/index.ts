import useSWR from 'swr'

import { client } from './client'
import { fetcher } from './helpers'
import { ErrorResponse, isErrorResponse, CreateMerkleResponse } from './types'

export * from './helpers'
export * from './client'
export * from './types'

export function useAPIResponse<T>(
  route: string,
  initialData?: T,
  options?: {
    revalidateOnFocus?: boolean
    revalidateOnMount?: boolean
    skipFetching?: boolean
    refreshInterval?: number | ((latestData: T | undefined) => number)
  },
): {
  data?: T
  isLoading: boolean
  isError: boolean
} {
  const { data, error } = useSWR<T, Error & { status?: number }>(
    options?.skipFetching ?? false ? null : route,
    fetcher,
    {
      fallbackData: initialData,
      ...options,
    },
  )

  return {
    data: error !== undefined ? undefined : data,
    isLoading: error === undefined && data === undefined,
    isError: Boolean(error),
  }
}

export async function createMerkleRoot(
  addresses: string[],
): Promise<CreateMerkleResponse | ErrorResponse> {
  const res = await client('POST', 'v1/merkle', {
    allowedAddresses: addresses,
  })

  if (!res.ok) {
    // attempt to decode error
    const errorResponse = await res.json()
    if (!isErrorResponse(errorResponse)) {
      throw new Error('Invalid response from server')
    }
    return errorResponse
  }

  return (await res.json()) as CreateMerkleResponse
}
