import { useState } from 'react'

export function useQuery(initialQuery: string | null = '') {
  const [query, setQuery] = useState<string>(initialQuery ?? '')
  const trimmedQuery = query.trim()
  const isDisabled = trimmedQuery === ''

  return { query, trimmedQuery, setQuery, isDisabled }
}
