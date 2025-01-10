<template>
    <div class="companions-container">
        <h5 class="companions-title">Your Companions</h5>
        <div v-if="userStore.companions && userStore.companions.length > 0">
            <small class="companions-hint">tap to start conversation</small>
            <div class="companions-list">
                <button 
                    v-for="companion in sortedCompanions" 
                    @click="selectCompanion(companion)"
                    :key="companion.id"
                    class="companion-card">
                    <div class="companion-header">
                        <span class="companion-name">{{ companion.username }}</span>
                        <span class="message-time" v-if="getLatestMessage(companion.id)">
                            {{ formatTime(getLatestMessage(companion.id).time) }}
                        </span>
                    </div>
                    <div class="message-preview" v-if="getLatestMessage(companion.id)">
                        {{ getLatestMessage(companion.id).content }}
                    </div>
                </button>
            </div>
        </div>
        <div v-else class="no-companions">
            <h3>NO COMPANIONS YET</h3>
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
    async mounted() {
        await this.userStore.fetchCompanions();
        await this.userStore.fetchLatestMessages();
    },
    methods: {
        async selectCompanion(companion){
            try{
                const response = await axios.get('api/conversations', {params:{companionId:companion.id}})
                this.userStore.conversationId = response.data.conversationId
                this.userStore.conversationPhoto = response.data.conversationPhoto
                this.userStore.conversationName = companion.username
                this.userStore.members = [];
                console.log('conversation id: ' + this.userStore.conversationId)
                await this.userStore.fetchMessages();
            }catch(error){console.log(error)}
        },
        getLatestMessage(companionId) {
            return this.userStore.latestMessages.find(msg => msg.companionId === companionId)
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
    },
    computed: {
        sortedCompanions() {
            return [...this.userStore.companions].sort((a, b) => {
                const msgA = this.getLatestMessage(a.id)
                const msgB = this.getLatestMessage(b.id)
                if (!msgA) return 1  // Put companions with no messages at the end
                if (!msgB) return -1
                return new Date(msgB.time) - new Date(msgA.time)
            })
    }
},
}
</script>

<style scoped>
.companions-container {
    padding: 20px;
    max-width: 600px;
    margin: 0 auto;
}

.companions-title {
    text-align: center;
    margin-bottom: 20px;
    font-size: 1.5em;
}

.companions-hint {
    display: block;
    text-align: center;
    margin-bottom: 15px;
}

.companions-list {
    display: flex;
    flex-direction: column;
    gap: 10px;
}

.companion-card {
    position: relative;
    display: flex;
    flex-direction: column;
    padding: 15px;
    background: white;
    border: 1px solid #ddd;
    border-radius: 8px;
    text-align: left;
    transition: all 0.2s;
    cursor: pointer;
}

.companion-card:hover {
    background-color: #09BC8A;
    transform: translateY(-2px);
}

.companion-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 8px;
}

.companion-name {
    font-weight: bold;
    font-size: 1.1em;
}

.message-time {
    position: absolute;
    top: 10px;      
    right: 10px;
    font-size: 0.8em;
}

.message-preview {
    font-size: 0.9em;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.no-companions {
    text-align: center;
    color: gray;
    margin-top: 40px;
}
</style>