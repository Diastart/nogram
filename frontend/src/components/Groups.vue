<template>
    <h3>groups ↓</h3>
    <button 
        v-for="group in userStore.groups" 
        @click="selectGroup(group)"
        :key="group.id">
        {{ group.groupname }}
    </button>
</template>

<script>
import { useUserInformation } from '@/stores/myStore';
import axios from 'axios'
export default {
    name: 'Groups',
    setup(){
        const userStore = useUserInformation();
        return {userStore};
    },
    mounted() {
        this.userStore.fetchGroups();
    },
    methods: {
        async selectGroup(group) {
            try{
                const response = await axios.get('api/conversations/groups', {params:{groupId:group.id}})
                this.userStore.conversationId = response.data.conversationId
                this.userStore.conversationPhoto = response.data.conversationPhoto
                this.userStore.conversationName = group.groupname
                console.log('conversation id ' + this.userStore.conversationId)
                console.log('url to conversation photo ' + this.userStore.conversationPhoto)
                await this.userStore.fetchMessages();
            }catch(error){console.log(error)}
        }
    }
}
</script>