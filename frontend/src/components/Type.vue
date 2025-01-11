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
        <hr>
        <input 
            type="file"
            ref="photoInput"
            @change="sendPhoto"
            accept="image/*"
            style="display: none"
        >
        <button @click="$refs.photoInput.click()" class="send-button">
            📷
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
        },
        async sendPhoto(event) {
            const file = event.target.files[0];
            if (!file) return;

            const formData = new FormData();
            formData.append('photo', file);
            formData.append('conversationId', this.userStore.conversationId);
            if (this.message) {
                formData.append('content', this.message);
                this.message = '';
            }

            try {
                await axios.post('api/messages/photo', formData);
                this.userStore.fetchMessages();
                await this.userStore.fetchLatestMessages();
                await this.userStore.fetchLatestMessagesOfGroups();
            } catch (error) {console.error('Error uploading photo:', error);}
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