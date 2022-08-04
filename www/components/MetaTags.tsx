import Head from 'next/head'
import React from 'react'

interface Props {
  title?: string
  description?: string
  imageUrl?: string
  largeImage?: boolean
}

function MetaTags({ title, description, imageUrl, largeImage = true }: Props) {
  const renderedTitle = `allowlist${title !== undefined ? ` | ${title}` : ''}`
  return (
    <Head>
      <title>{renderedTitle}</title>
      <meta key="og:title" property="og:title" content={renderedTitle} />
      <meta name="twitter:title" content={renderedTitle} />
      <meta
        key="twitter:card"
        name="twitter:card"
        content={largeImage ? 'summary_large_image' : 'summary'}
      />
      {description !== undefined && (
        <>
          <meta name="description" content={description} />
          <meta
            key="og:description"
            property="og:description"
            content={description}
          />
          <meta name="twitter:description" content={description} />
        </>
      )}
      {imageUrl !== undefined && (
        <>
          <meta key="og:image" property="og:image" content={imageUrl} />
          <meta name="twitter:image" content={imageUrl} />
        </>
      )}
    </Head>
  )
}

export default React.memo(MetaTags)
