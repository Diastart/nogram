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
    <div v-if="userStore.members && userStore.members.length > 0">
        <input 
            v-model="newGroupName" 
            placeholder="New groupname..." 
            class="groupname-input"
        >
        <button @click="setGroupName" class="update-button">
            Update
        </button>

        <hr>

        <button @click="$refs.GroupPhotoInput.click()" class="change-button">
            Change Photo
        </button>
        <input 
            type="file" 
            ref="GroupPhotoInput" 
            @change="setGroupPhoto" 
            accept="image/*" 
            style="display: none"
        >
        <hr>
        <Members/>
    </div>
</template>

<script>
import axios from 'axios'
import { useUserInformation } from '@/stores/myStore';
import Members from './Members.vue';
export default {
    name: 'Conversation',
    setup(){
        const userStore = useUserInformation();
        return {userStore};
    },
    data(){
        return {
            newGroupName: '',
        }
    },
    components: {
        Members,
    },
    methods: {
        async setGroupName(){
            if (!this.newGroupName) return;
            try {
                const response = await axios.put('api/groups/groupname', { groupname: this.newGroupName, conversationId: this.userStore.conversationId });
                this.userStore.conversationName = response.data.groupname;
                this.newGroupName = '';
            } catch (error) {console.log(error)}
        },
        async setGroupPhoto(event){
            const file = event.target.files[0];
            if (!file) return;

            const formData = new FormData();
            formData.append('photo', file);
            formData.append('conversationId', this.userStore.conversationId);

            try {
                const response = await axios.put('api/groups/photo', formData);
                this.userStore.conversationPhoto = response.data.photo;
            } catch (error) {console.log(error)}
        }
    }
}
</script>

<style scoped>
.conversation-header {
    width: 100%;
    padding: 20px 0;
    background-color: #808080;  
    display: flex;
    justify-content: flex-start;  
    border-bottom: 1px solid #2A3241;
    border-radius: 10px;
}

.header-content {
   display: flex;
   align-items: center;
   margin-left: 40px; 
   gap: 20px;  
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

.change-button, .update-button {
    padding: 10px 20px;
    background-color: #09BC8A;
    color: white;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    font-size: 1em;
    transition: opacity 0.2s ease;
    width: 250px;
}

.change-button:hover, .update-button:hover {
    opacity: 0.9;
}

.groupname-input {
    padding: 10px 15px;
    border-radius: 8px;
    border: 2px solid #ddd;
    font-size: 1em;
    width: 250px;
    height: 42px;
}
</style>