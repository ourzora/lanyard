import React from 'react'
import NextHead from 'next/head'

function Head() {
  return (
    <NextHead>
      <link rel="apple-touch-icon" sizes="180x180" href="/icon180.png" />
      <link rel="icon" type="image/png" sizes="32x32" href="/icon32.png" />
      <link rel="icon" type="image/png" sizes="16x16" href="/icon16.png" />
      <link rel="manifest" href="/site.webmanifest" />
      <link rel="shortcut icon" href="/favicon.ico" />
      <meta name="msapplication-TileColor" content="#aaaaaa" />
      <meta name="theme-color" content="#aaaaaa" />
      <meta
        name="viewport"
        content="width=device-width, initial-scale=1, maximum-scale=1"
      />
    </NextHead>
  )
}

export default React.memo(Head)
