import {createRouter, createWebHistory} from 'vue-router'

const router = createRouter({
    history: createWebHistory(),
    routes: [
        {
            path: '/',
            name: 'home',
            component: () => import('../components/Home.vue')
        },
        {
            path: '/about',
            name: 'about',
            component: () => import('../components/AboutPage.vue')
        },
        {
            path: '/how-it-works',
            name: 'how it works',
            component: () => import('../components/HowItWorksPage.vue')
        }
    ]
})

export default router