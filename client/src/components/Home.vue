<script setup lang="ts">
import { ref } from 'vue';
import SearchBox from './SearchBox.vue';
import LoadingPage from './LoadingPage.vue';
import { findRecommendedRoute, type RecommendedRoute } from '@/services/route.service';
import { useStore } from '../store/store'
import { useToast } from 'vue-toastification'
import { capitalizeFirstLetter } from '../utils/utils';
import Map from './Map.vue';

type State = 'search' | 'loading' | 'results'

const store = useStore()
const toast = useToast()
const state = ref<State>('search')

const startRouteSearch = async () => {
    const location = store.getLocation
    const source = location.source
    const destination = location.destination

    if (!source || !destination) {
        return
    }

    store.clearRecommendedRoute()
    state.value = 'loading'

    try {
        const response = await findRecommendedRoute(source, destination)
        if (response.error) {
            throw new Error(response.error)
        }
        store.setRecommendedRoute(response.data as RecommendedRoute)
    } catch (error: any) {
        const errMsg = capitalizeFirstLetter(error.message || 'Something went wrong');
        toast.error(errMsg)
        state.value = 'search'
    }
}

const newSearch = () => {
    state.value = 'search'
    store.clearLocation()
    store.clearRecommendedRoute()
}

const showResults = () => {
    state.value = 'results'
}

</script>

<template>
    <SearchBox v-if="state === 'search'" @findRoutes="startRouteSearch" />
    <LoadingPage v-else-if="state === 'loading'" @complete="showResults" />
    <Map v-else @new-search="newSearch"/>
</template>
