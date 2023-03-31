const ensSubgraphUrl = 'https://api.thegraph.com/subgraphs/name/ensdomains/ens'

const query = `
query DomainsQuery($names: [String!]!) {
  domains(where: { name_in: $names }) {
    resolvedAddress {
      id
    }
    name
  }
}
`

const resolveEnsDomainsBatch = async (
  ensNames: string[],
): Promise<{ [name: string]: string }> => {
  const variables = {
    names: ensNames,
  }
  const response = await fetch(ensSubgraphUrl, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ query, variables }),
  })
  const {
    data,
  }: {
    data?: {
      domains: { resolvedAddress: { id: string } | null; name: string }[]
    }
  } = await response.json()
  if (!data) {
    throw new Error('No data returned from subgraph')
  }

  return data.domains.reduce((acc, domain) => {
    if (domain.resolvedAddress !== null) {
      acc[domain.name] = domain.resolvedAddress.id
    }
    return acc
  }, {} as { [name: string]: string })
}

export const resolveEnsDomain = async (ensName: string): Promise<string> => {
  const domains = await resolveEnsDomainsBatch([ensName])
  if (!domains[ensName]) {
    throw new Error(`No address found for ENS name ${ensName}`)
  }
  return domains[ensName]
}

export const resolveEnsDomains = async (
  ensNames: string[],
): Promise<{ [name: string]: string }> => {
  const batches = chunk(ensNames, 100)
  const results = await Promise.all(batches.map(resolveEnsDomainsBatch))
  return results.reduce((acc, batch) => {
    return { ...acc, ...batch }
  }, {} as { [name: string]: string })
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const chunk = (arr: any[], size: number): any[][] => {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const chunks: any[][] = []
  let i = 0
  while (i < arr.length) {
    chunks.push(arr.slice(i, i + size))
    i += size
  }
  return chunks
}

export const isENSLike = (query: string) => query.endsWith('.eth')
