export function parseAddressesFromText(text: string) {
  // split the addresses by comma, space, or newline
  // and trim each one
  return text
    .split(/[,\s\n]+/)
    .map((s) => s.trim())
    .filter((s) => s.length > 0)
}
