// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-nocheck
import Prism from 'prism-react-renderer/prism'

// Ref: https://github.com/FormidableLabs/prism-react-renderer#faq
// (open the "How do I add more language highlighting support?" question)

export default function loadSolidityLanguageHighlighting() {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const _global: { Prism?: any } | undefined =
    typeof global !== 'undefined'
      ? global
      : typeof window !== 'undefined'
      ? window
      : undefined
  if (_global !== undefined) {
    _global.Prism = Prism
  }

  require('prismjs/components/prism-solidity')
}
