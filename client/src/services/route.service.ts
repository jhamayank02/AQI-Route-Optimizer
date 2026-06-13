import { API_URLS } from "@/constants/api"
import { get, type ApiResponse } from "./api.service"
import type { Location } from "./location.service"

export interface RecommendedRoute {
    selected_route_index: number
    selected_route: RouteAnalysis
    all_routes: RouteAnalysis[]
}

export interface Route {
    distance_km: number,
    duration_minutes: number,
    coordinates: Coordinates[]
}

export interface RouteAnalysis {
    route: Route
    aqi_samples: AqiSample[]
    average_aqi: number
    max_aqi: number
    route_score: number
    recommendation: string
    sampling_strategy: string
    is_selected: boolean
    selection_reason: string
}

export interface Coordinates {
    lat: number,
    lng: number
}

export interface AqiSample {
    lat: number
    lng: number
    aqi: number
}

export async function findRecommendedRoute(source: Location, destination: Location) {
    try {
        const response: ApiResponse<RecommendedRoute> = await get<ApiResponse<RecommendedRoute>>(`${API_URLS.FIND_RECOMMENDED_ROUTE}?src_lat=${source.lat}&src_lng=${source.lng}&dst_lat=${destination.lat}&dst_lng=${destination.lng}`)
        return response
    } catch (error: any) {
        return { error: error.message || 'Something went wrong, please try again later', data: null }
    }
}