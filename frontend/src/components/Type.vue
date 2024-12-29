<template>
    <input v-model="message" placeholder="Type a message..." required>
    <button @click="sendMessage">Send</button>
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
                this.userStore.fetchMessages();
            }catch(error){console.log(error)}
        }
    }
}
</script>