<template>
    <ul>
        <li v-for="companion in userStore.companions" :key="companion.id"> {{ companion.username }} </li>
    </ul>
</template>

<script>
import axios from 'axios';
import { useUserInformation } from '@/stores/myStore';
export default {
    name: 'Companions',
    setup(){
        const userStore = useUserInformation();
        return {userStore};
    },
    mounted() {
        this.userStore.fetchCompanions()
    },
}
</script>