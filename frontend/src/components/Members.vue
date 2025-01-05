<template>
    <div v-if="userStore.members && userStore.members.length > 0" class="members-container">
        <h5>Members of the group</h5>
        <div class="members-list">
            <div v-for="member in userStore.members" 
                style="color: black;"
                :key="member.id"
                class="member-item">
                {{ member.username }}
            </div>
        </div>
        <Candidates/>
        <button @click="leave()" class="leave-button">Leave</button>
    </div>
</template>

<script>
import { useUserInformation } from '@/stores/myStore';
import Candidates from './Candidates.vue'
import axios from 'axios';
export default {
    name: 'Members',
    setup(){
        const userStore = useUserInformation();
        return {userStore};
    },
    components: {
        Candidates,
    },
    methods: {
        async leave(){
            try {
                await axios.delete('api/members', {params:{conversationId: this.userStore.conversationId}});
                console.log('you left the group')
                this.userStore.conversationId = 0;
                this.userStore.members = [];
                this.userStore.fetchGroups();
            } catch (error) {console.log(error)}
        },
    }
}
</script>

<style scoped>
.members-container {
   padding: 20px;
}

.members-list {
   display: grid;
   gap: 10px;
   grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
}

.member-item {
   padding: 12px;
   background-color: white;
   border: 1px solid #ddd;
   border-radius: 8px;
   text-align: center;
}

.leave-button {
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