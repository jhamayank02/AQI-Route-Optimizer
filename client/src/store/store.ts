import type { Location } from "@/services/location.service"
import type { RecommendedRoute } from "@/services/route.service"
import { defineStore } from "pinia"

export const useStore = defineStore("main", {
    state: () => ({
        source: null as Location | null,
        destination: null as Location | null,
        recommendedRoute: null as RecommendedRoute | null
    }),

    actions: {
        setSource(source: Location) {
            this.source = source
        },

        setDestination(destination: Location) {
            this.destination = destination
        },

        clearLocation() {
            this.source = null
            this.destination = null
        },
        
        setRecommendedRoute(recommendedRoute: RecommendedRoute){
            this.recommendedRoute = recommendedRoute
        },

        clearRecommendedRoute() {
            this.recommendedRoute = null
        }
    },

    getters: {
        getLocation(): { source: Location | null, destination: Location | null } {
            return {
                source: this.source,
                destination: this.destination,
            }
        },

        getRecommendedRoute(): RecommendedRoute | null {
            return this.recommendedRoute
        }
    }
})