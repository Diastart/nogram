<template>
    <div class="groups-container">
        <h5 class="groups-title">Your Groups</h5>
        <div v-if="userStore.groups && userStore.groups.length > 0">
            <small class="groups-hint">tap to start conversation</small>
            <div class="groups-list">
                <button 
                    v-for="group in sortedGroups" 
                    @click="selectGroup(group)"
                    :key="group.id"
                    class="group-card">
                    <div class="group-header">
                        <div class="group-name">{{ group.groupname }}</div>
                        <div class="message-time" v-if="getLatestMessage(group.id)">
                            {{ formatTime(getLatestMessage(group.id).time) }}
                        </div>
                    </div>
                    <div class="message-preview" v-if="getLatestMessage(group.id)">
                        {{ getLatestMessage(group.id).content }}
                    </div>
                </button>
            </div>
        </div>
        <div v-else class="no-groups">
            <h3>NO GROUPS YET</h3>
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
    async mounted() {
        await this.userStore.fetchGroups();
        await this.userStore.fetchLatestMessagesOfGroups();
    },
    computed: {
        sortedGroups() {
            return [...this.userStore.groups].sort((a, b) => {
                const msgA = this.getLatestMessage(a.id)
                const msgB = this.getLatestMessage(b.id)
                if (!msgA) return 1
                if (!msgB) return -1
                return new Date(msgB.time) - new Date(msgA.time)
            })
        }
    },
    methods: {
        async selectGroup(group) {
            try{
                const response = await axios.get('api/conversations/groups', {params:{groupId:group.id}})
                this.userStore.conversationId = response.data.conversationId
                this.userStore.conversationPhoto = response.data.conversationPhoto
                this.userStore.conversationName = group.groupname
                await this.userStore.fetchMembers(group.id);
                await this.userStore.fetchMessages();
            }catch(error){console.log(error)}
        },
        getLatestMessage(groupId) {
            return this.userStore.latestMessagesOfGroups.find(msg => msg.groupId === groupId)
        },
        formatTime(timeString) {
            return new Date(timeString).toLocaleTimeString([], { 
                hour: '2-digit', 
                minute: '2-digit' 
            })
        }
    }
}
</script>

<style scoped>
.groups-container {
    padding: 20px;
    max-width: 600px;
    margin: 0 auto;
}

.groups-title {
    text-align: center;
    margin-bottom: 20px;
    font-size: 1.5em;
}

.groups-hint {
    display: block;
    text-align: center;
    color: #666;
    margin-bottom: 15px;
}

.groups-list {
    display: flex;
    flex-direction: column;
    gap: 10px;
}

.group-card {
    position: relative;
    padding: 15px;
    background: white;
    border: 1px solid #ddd;
    border-radius: 8px;
    text-align: left;
    transition: all 0.2s;
    cursor: pointer;
}

.group-card:hover {
    background-color: #09BC8A;
    transform: translateY(-2px);
}

.group-header {
    margin-bottom: 8px;
}

.group-name {
    font-weight: bold;
    font-size: 1.1em;
}

.message-time {
    position: absolute;
    top: 10px;
    right: 10px;
    color: #666;
    font-size: 0.8em;
}

.message-preview {
    color: #666;
    font-size: 0.9em;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.no-groups {
    text-align: center;
    color: gray;
    margin-top: 40px;
}
</style>