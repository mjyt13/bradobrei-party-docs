const apiScheme = import.meta.env.VITE_API_SCHEME || 'http'
const apiHost = import.meta.env.VITE_API_HOST || 'localhost'
const apiPort = import.meta.env.VITE_API_PORT || '9000'
const apiBasePath = import.meta.env.VITE_API_BASE_PATH || '/api/v1'

const apiOrigin = `${apiScheme}://${apiHost}:${apiPort}`

export const apiConfig = {
  apiOrigin,
  apiBasePath,
  apiBaseUrl: import.meta.env.DEV ? apiBasePath : `${apiOrigin}${apiBasePath}`,
  docsUrl: import.meta.env.DEV ? '/swagger/index.html' : `${apiOrigin}/swagger/index.html`,
}
