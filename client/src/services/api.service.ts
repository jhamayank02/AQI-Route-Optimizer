export interface ApiResponse<T> {
    message: string
    data: T | null
    error: string | null
}

export async function get<T>(url: string): Promise<T> {
    const response = await fetch(url)

    if (!response.ok) {
        let errorMessage = `Request failed: ${response.status}`

        try {
            const errorBody = await response.json()

            errorMessage =
                errorBody.message ||
                errorBody.error ||
                errorMessage
        } catch {
            // Response body wasn't JSON
        }

        throw new Error(errorMessage)
    }

    return response.json()
}

export async function post<T>(url: string, headers: Record<string, string>, body: any): Promise<T> {
    const response = await fetch(url, {
        method: 'POST',
        headers,
        body: JSON.stringify(body)
    });

    if (!response.ok) {
        let errorMessage = `Request failed: ${response.status}`

        try {
            const errorBody = await response.json()

            errorMessage =
                errorBody.message ||
                errorBody.error ||
                errorMessage
        } catch {
            // Response body wasn't JSON
        }

        throw new Error(errorMessage)
    }

    return response.json()
}