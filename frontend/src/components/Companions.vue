<template>
    <h3>Companions:</h3>
    <button 
        v-for="companion in userStore.companions" 
        @click="selectCompanion(companion)"
        :key="companion.id">
        {{ companion.username }}
    </button>
</template>

<script>
import { useUserInformation } from '@/stores/myStore';
import axios from 'axios'
export default {
    name: 'Companions',
    setup(){
        const userStore = useUserInformation();
        return {userStore};
    },
    mounted() {
        this.userStore.fetchCompanions()
    },
    methods: {
        async selectCompanion(companion){
            try{
                const response = await axios.get('api/conversations', {params:{companionId:companion.id}})
                this.userStore.conversationId = response.data.conversationId
                this.userStore.conversationPhoto = response.data.conversationPhoto
                this.userStore.conversationName = companion.username
                console.log('conversation id ' + this.userStore.conversationId)
                console.log('url to conversation photo ' + this.userStore.conversationPhoto)
                await this.userStore.fetchMessages();
            }catch(error){console.log(error)}
        },
    }
}
</script>