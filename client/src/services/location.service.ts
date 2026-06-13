import { API_URLS } from "@/constants/api";
import { get, type ApiResponse } from "./api.service";

export interface Location {
    confidence: number,
    country: string,
    lat: number,
    lng: number,
    name: string,
    region: string,
    label: string
}

export async function searchLocation(query: string) {
    try {
        const response: ApiResponse<Location[]> = await get<ApiResponse<Location[]>>(`${API_URLS.SEARCH_LOCATION}?query=${query}`)

        return response
    } catch (error: any) {
        return { error: error.message || 'Something went wrong, please try again later', data: null }
    }
}