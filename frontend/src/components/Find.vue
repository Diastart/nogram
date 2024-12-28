<template>
    <input v-model="companion" placeholder="Find & add companion" required><br>
    <button @click="findCompanion">press</button>
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