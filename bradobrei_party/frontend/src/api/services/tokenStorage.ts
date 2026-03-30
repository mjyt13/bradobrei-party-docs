const TOKEN_STORAGE_KEY = 'bradobrei-party.auth.token'

export const tokenStorage = {
  get() {
    return window.localStorage.getItem(TOKEN_STORAGE_KEY)
  },
  set(token: string) {
    window.localStorage.setItem(TOKEN_STORAGE_KEY, token)
  },
  clear() {
    window.localStorage.removeItem(TOKEN_STORAGE_KEY)
  },
  has() {
    return Boolean(window.localStorage.getItem(TOKEN_STORAGE_KEY))
  },
}
