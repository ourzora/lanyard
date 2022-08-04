module.exports = {
  webpack5: true,
  async rewrites() {
    const rewrites = [
      { source: '/next-api/:path*', destination: '/api/:path*' },
    ]

    // forward /api/ requests to the real api
    // this should be done be fastly in production
    rewrites.push({
      source: '/api/:path*',
      destination: process.env.API_URL + '/api/:path*',
    })

    return rewrites
  },
  sentry: {
    disableServerWebpackPlugin: process.env.SENTRY_AUTH_TOKEN === undefined,
    disableClientWebpackPlugin: process.env.SENTRY_AUTH_TOKEN === undefined,
  },
}
