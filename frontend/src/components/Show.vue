<template>
    <div class="messages-container">
        <h3>Messages:</h3>
        <small>you can tap to see options</small>
        <div class="messages">
            <div class="message" v-for="message in userStore.messages" @click="options(message)">
                <div class="message-content">
                    <div class="sender">{{ message.senderName }}</div>
                    <div class="content">{{ message.content }}</div>
                    <div class="time">{{ message.time }}</div>
                </div>
 
                <div v-if="selectedMessageId === message.id" class="options-popup">
                    <div class="redirect-section">
                        <h5 style="color:black">Redirect to:</h5>
                        <div class="companions-list">
                            <button 
                                v-for="companion in userStore.companions" 
                                @click="redirectMessage(message, companion)"
                                class="companion-button"
                            >
                                {{ companion.username }}
                            </button>
                        </div>
                    </div>
                    <button @click="deleteMessage(message)" class="delete-button">Delete Message</button>
                </div>
            </div>
        </div>
    </div>
 </template>

<script>
import axios from 'axios'
import { useUserInformation } from '@/stores/myStore';
export default {
    name: 'Show',
    setup(){
        const userStore = useUserInformation();
        return {userStore};
    },
    mounted() {
        this.userStore.fetchMessages();
    },
    data(){
        return {
            selectedMessageId: null,
        }
    },
    methods: {
        options(message){
            this.selectedMessageId = this.selectedMessageId === message.id ? null : message.id;
        },
        async redirectMessage(message, companion) {
            try{
                const response = await axios.get('api/conversations', {params:{companionId:companion.id}})
                await axios.put('api/messages', {messageId: message.id, conversationId: response.data.conversationId, senderName: message.senderName});
                console.log('Redirected');
                this.userStore.fetchMessages();
            }catch(error){console.log(error)}
        },
        async deleteMessage(message) {
            try{
                await axios.delete('api/messages', {params:{messageId: message.id}});
                console.log('Deleted');
                this.userStore.fetchMessages();
            }catch(error){console.log(error)}
        }
    }
}
</script>

<style scoped>
.messages-container {
   padding: 20px;
}

.message { 
    padding: 15px;
    margin: 10px 0;
    border-radius: 8px;
    border: 1px solid #ddd;
    transition: all 0.2s;
}

.message:hover {
   background-color: #09BC8A;
   transform: scale(1.01);
}

.message-content {
   margin-bottom: 10px;
}

.sender {
   font-weight: bold;
   margin-bottom: 5px;
}

.time {
   color: #666;
   font-size: 0.9em;
   margin-top: 5px;
}

.options-popup {
   margin-top: 15px;
   padding: 15px;
   background-color: #f5f5f5;
   border-radius: 4px;
}

.companions-list {
   display: grid;
   gap: 10px;
   margin: 10px 0;
}

.companion-button {
   padding: 8px;
   border: none;
   border-radius: 4px;
   background-color: #fff;
   cursor: pointer;
}

.companion-button:hover {
   background-color: #09BC8A;
}

.delete-button {
   width: 100%;
   padding: 8px;
   margin-top: 10px;
   background-color: #ff4444;
   color: white;
   border: none;
   border-radius: 4px;
   cursor: pointer;
}
</style>
