const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? '/api'

export const API_URLS = {
    SEARCH_LOCATION: `${API_BASE_URL}/locations/search`,
    FIND_RECOMMENDED_ROUTE: `${API_BASE_URL}/routes/recommendation`
}
