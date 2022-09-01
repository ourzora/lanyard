export function parseAddressesFromText(text: string) {
  // split the addresses by comma, space, or newline
  // and trim each one
  return text
    .split(/[,\s\n]+/)
    .map((s) => s.trim())
    .filter((s) => s.length > 0)
}

export const prepareAddresses = (
  addressesOrENSNames: string[],
  ensMap: Record<string, string>,
): {
  addresses: string[]
  dedupedAddresses: string[]
  unresolvedEnsNames: string[]
} => {
  const seenAddresses = new Set<string>()
  const unresolvedEnsNames: string[] = []
  const addresses: string[] = []
  const dedupedAddresses: string[] = []

  for (const addressOrENSName of addressesOrENSNames) {
    const lowercasedAddressOrENSName = addressOrENSName.toLowerCase()
    if (addressOrENSName.includes('.')) {
      const addressFromEns: string | undefined =
        ensMap[lowercasedAddressOrENSName]?.toLowerCase()

      if (addressFromEns !== undefined) {
        addresses.push(addressFromEns)
        if (seenAddresses.has(addressFromEns)) {
          continue
        }
        seenAddresses.add(addressFromEns)
        dedupedAddresses.push(addressFromEns)
      } else {
        unresolvedEnsNames.push(lowercasedAddressOrENSName)
      }
    } else {
      addresses.push(addressOrENSName)
      if (seenAddresses.has(lowercasedAddressOrENSName)) {
        continue
      }
      seenAddresses.add(lowercasedAddressOrENSName)
      dedupedAddresses.push(addressOrENSName)
    }
  }

  return {
    addresses,
    dedupedAddresses,
    unresolvedEnsNames,
  }
}
