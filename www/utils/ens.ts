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

export const resolveEnsDomains = async (
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
    data?: { domains: { resolvedAddress: { id: string }; name: string }[] }
  } = await response.json()
  if (!data) {
    throw new Error('No data returned from subgraph')
  }

  return data.domains.reduce((acc, domain) => {
    acc[domain.name] = domain.resolvedAddress.id
    return acc
  }, {} as { [name: string]: string })
}
