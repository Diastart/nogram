<template>
    <h3>Hello my dear friend {{ userStore.username }}</h3>
    <p>-----------------------------------------------------------------</p>
    <Find/>
    <p>-----------------------------------------------------------------</p>
</template>

<script>
import { useUserInformation } from '@/stores/myStore';
import Find from '../components/Find.vue';
export default {
    name: 'Chat',
    components: {
        Find,
    },
    setup(){
        const userStore = useUserInformation();
        return {userStore};
    },
}
</script>