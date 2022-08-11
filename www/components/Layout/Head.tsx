import React from 'react'
import NextHead from 'next/head'

function Head() {
  return (
    <NextHead>
      <meta
        name="viewport"
        content="width=device-width, initial-scale=1, maximum-scale=1"
      />
    </NextHead>
  )
}

export default React.memo(Head)
