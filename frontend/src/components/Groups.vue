<template>
    <div class="groups-container">
        <h3>Your Groups</h3>
        <div class="groups-grid">
            <button 
                v-for="group in userStore.groups" 
                @click="selectGroup(group)"
                :key="group.id"
                class="group-button">
                {{ group.groupname }}
            </button>
        </div>
    </div>
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

<style scoped>
.groups-container {
   padding: 20px;
}

.groups-grid {
   display: grid;
   gap: 10px;
   grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
}

.group-button {
   padding: 12px;
   border: 1px solid #ddd;
   border-radius: 8px;
   background-color: white;
   transition: all 0.2s;
}

.group-button:hover {
   background-color: #09BC8A;
   transform: scale(1.05);
}
</style>