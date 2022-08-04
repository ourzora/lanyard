import { Html, Head, Main, NextScript } from 'next/document'
import { bodyStyles } from 'utils/theme'

export default function Document() {
  return (
    <Html>
      <Head />
      <body className={bodyStyles}>
        <Main />
        <NextScript />
      </body>
    </Html>
  )
}
