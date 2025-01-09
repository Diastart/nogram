<template>
    <div class="message-input">
        <input 
            v-model="message" 
            placeholder="Type a message..." 
            class="input"
            required
        >
        <button @click="sendMessage" class="send-button">
            🚀
        </button>
    </div>
</template>

<script>
import axios from 'axios';
import { useUserInformation } from '@/stores/myStore';
export default {
    name: 'Type',
    setup(){
        const userStore = useUserInformation();
        return {userStore};
    },
    data(){
        return {
            message: '',
        }
    },
    methods: {
        async sendMessage(){
            try{
                await axios.post('api/messages', {conversationId: this.userStore.conversationId ,content: this.message});
                console.log('message is sent successfully');
                this.message = '';
                await this.userStore.fetchMessages();
                await this.userStore.fetchLatestMessages();
                await this.userStore.fetchLatestMessagesOfGroups();
            }catch(error){console.log(error)}
        }
    }
}
</script>

<style scoped>
.message-input {
   display: flex;
   gap: 10px;
   padding: 15px;
   border-top: 1px solid #ddd;
}

.input {
   flex: 1;
   padding: 12px;
   border: 1px solid #ddd;
   border-radius: 20px;
}

.send-button {
   padding: 10px 20px;
   background-color: #09BC8A;
   color: white;
   border: none;
   border-radius: 20px;
   transition: transform 0.2s;
}

.send-button:hover {
   transform: scale(1.05);
}
</style>