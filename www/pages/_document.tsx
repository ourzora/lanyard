import { Html, Head, Main, NextScript } from 'next/document'
import { bodyStyles } from 'utils/theme'

export default function Document() {
  return (
    <Html>
      <Head>
        <link rel="preconnect" href="https://fonts.googleapis.com" />
        <link
          rel="preconnect"
          href="https://fonts.gstatic.com"
          crossOrigin="anonymous"
        />
        <link
          href="https://fonts.googleapis.com/css2?family=Chivo:wght@400;700&display=swap"
          rel="stylesheet"
        />
      </Head>
      <body className={bodyStyles}>
        <Main />
        <NextScript />
      </body>
    </Html>
  )
}
