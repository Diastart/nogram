<template>
    <div v-if="this.userStore.candidates && this.userStore.candidates.length > 0" class="candidates-container">
        <h5>Candidates</h5>
        <small>press to add a candidate</small>
        <div class="candidates-grid">
            <button 
                v-for="candidate in userStore.candidates"
                @click="addtoGroup(candidate)"
                :key="candidate.id"
                class="candidate-button">
                {{ candidate.username }}
            </button>
        </div>
    </div>
 </template>

<script>
import { useUserInformation } from '@/stores/myStore';
import axios from 'axios'
export default {
    name: 'Candidates',
    setup(){
        const userStore = useUserInformation();
        return {userStore};
    },
    methods: {
        async addtoGroup(candidate){
            try {
                const response = await axios.post('api/candidates', {conversationId: this.userStore.conversationId, userId: candidate.id}) 
                console.log('candidate ' + candidate.username + ' is added')
                this.userStore.fetchMembers(response.data.groupId)
            }catch (error) {console.log(error)} 
        },
    }
}
</script>

<style scoped>
.candidates-container {
   padding: 20px;
   text-align: center;
}

.candidates-grid {
   display: grid;
   gap: 10px;
   grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
   margin-top: 15px;
}

.candidate-button {
   padding: 12px;
   border: 1px solid #ddd;
   border-radius: 8px;
   background-color: white;
   transition: all 0.2s;
}

.candidate-button:hover {
   background-color: #09BC8A;
   transform: scale(1.05);
}
</style>