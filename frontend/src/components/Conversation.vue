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
 </template>

<script>
import { useUserInformation } from '@/stores/myStore';
export default {
    name: 'Conversation',
    setup(){
        const userStore = useUserInformation();
        return {userStore};
    },
}
</script>

<style scoped>
.conversation-header {
   width: 100%;
   padding: 20px 0;
   background-color: #808080;  /* Changed to grey */
   display: flex;
   justify-content: flex-start;  /* Changed from center to flex-start */
   border-bottom: 1px solid #2A3241;
}

.header-content {
   display: flex;
   align-items: center;
   margin-left: 40px;  /* Add space from left edge */
   gap: 20px;  /* Space between photo and text */
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