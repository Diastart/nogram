<template>
    <h3>Hello my dear friend {{ userStore.username }}</h3>
    <p>-----------------------------------------------------------------</p>
    <Find/>
    <p>-----------------------------------------------------------------</p>
    <Companions/>
</template>

<script>
import { useUserInformation } from '@/stores/myStore';
import Find from '../components/Find.vue';
import Companions from '../components/Companions.vue'
export default {
    name: 'Chat',
    components: {
        Find,
        Companions,
    },
    setup(){
        const userStore = useUserInformation();
        return {userStore};
    },
}
</script>