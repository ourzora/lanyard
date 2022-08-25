export function parseAddressesFromText(
  text: string,
  removeDuplicates: boolean,
): {
  addresses: string[]
  hasDuplicates: boolean
} {
  let hasDuplicates = false

  let addrs: { [key: string]: boolean } = {}

  // split the addresses by comma, space, or newline
  // and trim each one
  return {
    addresses: text
      .split(/[,\s\n]+/)
      .map((s) => {
        let a = s.trim().toLowerCase()

        if (addrs[a]) {
          if (removeDuplicates) {
            return ''
          }

          hasDuplicates = true
        }

        addrs[a] = true

        return a
      })
      .filter((s) => s.length > 0),
    hasDuplicates,
  }
}
