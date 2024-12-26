<template>
    <h1>Hello my dear friend {{ userStore.username }}</h1>
</template>

<script>
import axios from 'axios'
import { useUserInformation } from '@/stores/myStore';
export default {
    name: 'Chat',
    setup(){
        const userStore = useUserInformation();
        return {userStore};
    },
}
</script>