export function copyToClipboard(text: string) {
  if (typeof navigator !== 'undefined' && navigator.clipboard !== undefined) {
    navigator.clipboard.writeText(text)
  }
}
