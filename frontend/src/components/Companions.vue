<template>
    <div class="companions-container">
        <h5>Your Companions</h5>
        <div class="companions-grid">
            <button 
                v-for="companion in userStore.companions" 
                @click="selectCompanion(companion)"
                :key="companion.id"
                class="companion-button">
                {{ companion.username }}
            </button>
        </div>
    </div>
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
                this.userStore.members = [];
                console.log('conversation id ' + this.userStore.conversationId)
                console.log('url to conversation photo ' + this.userStore.conversationPhoto)
                await this.userStore.fetchMessages();
            }catch(error){console.log(error)}
        },
    }
}
</script>

<style scoped>
.companions-container {
   padding: 20px;
   text-align: center;
}

.companions-grid {
   display: grid;
   gap: 10px;
   grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
   width: fit-content;  /* Add this line */
   margin: 0 auto;     /* Add this line */
}

.companion-button {
   padding: 12px;
   border: 1px solid #ddd;
   border-radius: 8px;
   background-color: white;
   transition: all 0.2s;
}

.companion-button:hover {
   background-color: #09BC8A;
   transform: scale(1.05);
}
</style>