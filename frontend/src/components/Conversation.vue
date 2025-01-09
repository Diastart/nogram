<template>
    <div v-if="userStore.conversationId" class="conversation-header">
        <div class="header-content">
            <div class="profile-container">
                <img 
                    :src="userStore.conversationPhoto" 
                    class="profile-photo"
                    alt="Profile photo"
                >
            </div>
            <div class="chat-title">
                {{userStore.conversationName}}
            </div>
        </div>
    </div>
    <Members/>
</template>

<script>
import { useUserInformation } from '@/stores/myStore';
import Members from './Members.vue';
export default {
    name: 'Conversation',
    setup(){
        const userStore = useUserInformation();
        return {userStore};
    },
    components: {
        Members,
    }
}
</script>

<style scoped>
.conversation-header {
    width: 100%;
    padding: 20px 0;
    background-color: #808080;  
    display: flex;
    justify-content: flex-start;  
    border-bottom: 1px solid #2A3241;
    border-radius: 10px;
}

.header-content {
   display: flex;
   align-items: center;
   margin-left: 40px; 
   gap: 20px;  
}

.profile-container {
   position: relative;
}

.profile-photo {
   width: 60px;
   height: 60px;
   border-radius: 50%;
   border: 2px solid #2A3241;
   object-fit: cover;
}

.chat-title {
   color: white;
   font-size: 18px;
}
</style>