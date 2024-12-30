<template>
    <input v-model="companion" placeholder="Type companion name..." required>
    <button @click="findCompanion">add</button>
</template>

<script>
import axios from 'axios'
import { useUserInformation } from '@/stores/myStore';
export default {
    name : 'Find',
    setup(){
        const userStore = useUserInformation();
        return {userStore};
    },
    data(){
        return {
            companion: '',
        }
    },
    methods: {
        async findCompanion(){
            try {
                await axios.get('api/dialogs', {params: {companion: this.companion}});
                console.log(this.companion + ' is added successfully');
                this.companion = '';
                await this.userStore.fetchCompanions();
            }catch(error){console.log(error)}
        },
    }

}
</script>