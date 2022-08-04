import React from 'react'
import { AppProps } from 'next/app'

import Providers from 'components/Providers'
import Layout from 'components/Layout'

import '../styles/global.css'

function App({ Component, pageProps }: AppProps) {
  return (
    <Providers>
      <Layout>
        <Component {...pageProps} />
      </Layout>
    </Providers>
  )
}

export default App
