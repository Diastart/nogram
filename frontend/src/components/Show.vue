<template>
    <div class="messages-container">
        <h3>Messages:</h3>
        <small>tap to see options</small>
        <div class="messages">
            <div class="message" v-for="message in userStore.messages" @click="options(message)">
                <div class="message-content">
                    <div class="sender">{{ message.senderName }}</div>
                    <div class="content">
                        {{ message.content }}
                        <span class="checkmark">{{ message.checkmark }}</span>
                        <span v-if="message.reaction" class="reaction">{{ message.reaction }}</span>
                    </div>
                    <hr>
                    <div v-if="message.photo">
                        <img  
                            :src="message.photo" 
                            alt="Message photo"
                            style="max-width: 200px; margin-top: 5px;"
                        >
                        <hr>
                    </div>
                    <div class="time">{{ formatTime(message.time) }}</div>
                </div>

                <div v-if="selectedMessageId === message.id" >
                    <div class="options-popup">

                        <div class="response-section">
                            <h5 style="color:black">Respond:</h5>
                            <div class="response-input">
                                <input 
                                    v-model="responseText" 
                                    placeholder="Type your response..." 
                                    required
                                    class="response-field"
                                    @click.stop
                                >
                                <button @click="respondToMessage(message)" class="response-button">
                                    🚀
                                </button>
                            </div>
                        </div>
                        <hr style="color: black;">

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
                        <hr style="color: black;">

                        <div class="reaction-section">
                            <h5 style="color: black;">Comment</h5>
                            <div class="emoji-list">
                                <button 
                                    v-for="emoji in emojis" 
                                    @click="addReaction(message, emoji)"
                                    class="emoji-button"
                                >
                                    {{ emoji }}
                                </button>
                            </div>
                        </div>
                        <hr style="color: black;">
                    </div>
                    <button @click="deleteReaction(message)" class="uncomment-button">Uncomment</button>
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
            emojis: ['👍', '❤️', '😊', '😂', '🎉', '👏', '🤔', '👌'],
            responseText: '',
        }
    },
    methods: {
        options(message){
            this.selectedMessageId = this.selectedMessageId === message.id ? null : message.id;
        },
        async redirectMessage(message, companion) {
            try{
                const response = await axios.get('api/conversations', {params:{companionId:companion.id}})
                await axios.put('api/messages', {messageId: message.id, conversationId: response.data.conversationId, senderName: message.senderName, photo: message.photo});
                console.log('Redirected');
                await this.userStore.fetchMessages();
                await this.userStore.fetchLatestMessages();
                await this.userStore.fetchLatestMessagesOfGroups();
            }catch(error){console.log(error)}
        },
        async deleteMessage(message) {
            try{
                await axios.delete('api/messages', {params:{messageId: message.id}});
                console.log('Deleted');
                await this.userStore.fetchMessages();
                await this.userStore.fetchLatestMessages();
                await this.userStore.fetchLatestMessagesOfGroups();
            }catch(error){console.log(error)}
        },
        async addReaction(message, emoji) {
            try {
                await axios.put('api/reactions', {messageId: message.id, reaction: emoji});
                this.showEmojis = false;
                await this.userStore.fetchMessages();
                await this.userStore.fetchLatestMessages();
                await this.userStore.fetchLatestMessagesOfGroups();
            } catch(error) {console.log(error);}
        },
        async deleteReaction(message) {
            try {
                await axios.delete('api/reactions', {params: {messageId: message.id}});
                console.log('Uncommented');
                await this.userStore.fetchMessages();
                await this.userStore.fetchLatestMessages();
                await this.userStore.fetchLatestMessagesOfGroups();
            } catch (error) {console.log(error)}
        },
        formatTime(timeString) {
            return new Date(timeString).toLocaleTimeString([], { 
                year: 'numeric',
                month: 'numeric',
                day: 'numeric',
                hour: '2-digit',
                minute: '2-digit'
            })
        },
        async respondToMessage(message) {
            if (!this.responseText) return;
            try {
                await axios.post('api/messages', {conversationId: this.userStore.conversationId, content: `🪅Respond to ${message.content}🪅: ${this.responseText}`, photo: message.photo});
                this.responseText = '';
                await this.userStore.fetchMessages();
                await this.userStore.fetchLatestMessages();
                await this.userStore.fetchLatestMessagesOfGroups();
            } catch(error) {console.log(error);}
        },
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
   position: relative; 
}

.sender {
   font-weight: bold;
   margin-bottom: 5px;
}

.time {
   font-size: 0.9em;
   margin-top: 5px;
   background-color: #706F6F;
   border-radius: 10px;
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

.uncomment-button {
   width: 100%;
   padding: 8px;
   margin-top: 10px;
   background-color: gray;
   color: white;
   border: none;
   border-radius: 4px;
   cursor: pointer;
}

.reaction-section {
    margin-top: 10px;
}

.reaction-button {
    width: 100%;
    padding: 8px;
    background-color: white;
    border: 1px solid #ddd;
    border-radius: 4px;
    cursor: pointer;
}

.emoji-list {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 5px;
    margin-top: 10px;
}

.emoji-button {
    padding: 8px;
    font-size: 1.2em;
    background: white;
    border: 1px solid #ddd;
    border-radius: 4px;
    cursor: pointer;
}

.emoji-button:hover {
    background-color: #f5f5f5;
}

.reaction {
    position: absolute;
    bottom: -10px;
    left: 10px;
    background: #f0f0f0;
    padding: 2px 6px;
    border-radius: 12px;
    font-size: 0.9em;
    border: 1px solid #ddd;
}

.content {
    position: relative;
    padding-right: 20px;  
}

.checkmark {
    position: absolute;
    right: 5px;
    bottom: 0;
    font-size: 0.8em;
}

.response-section {
    margin-bottom: 15px;
}

.response-input {
    display: flex;
    gap: 10px;
    margin-top: 10px;
}

.response-field {
    flex: 1;
    padding: 8px;
    border: 1px solid #ddd;
    border-radius: 4px;
}

.response-button {
    padding: 8px 16px;
    background-color: #09BC8A;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
}

.response-button:hover {
    opacity: 0.9;
}
</style>
