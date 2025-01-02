<template>
    <div class="create-group">
        <div class="group-form">
            <input 
                v-model="groupName" 
                placeholder="Group name..." 
                class="group-input"
                required
            >
            
            <div class="members-section">
                <h4 style="color: black;">Select Members</h4>
                <div class="members-grid">
                    <div v-for="companion in userStore.companions" 
                         :key="companion.id" 
                         class="member-item">
                        <label style="color: black;" class="checkbox-label">
                            <input 
                                type="checkbox"
                                v-model="groupMembers"
                                :value="companion.id"
                            >
                            {{ companion.username }}
                        </label>
                    </div>
                </div>
            </div>
 
            <button @click="createGroup" class="create-button">
                Create Group
            </button>
        </div>
    </div>
 </template>

<script>
import axios from 'axios'
import { useUserInformation } from '@/stores/myStore';
export default {
    name: 'Create',
    data(){
        return {
            groupName: '',
            groupMembers: [],
        }
    },
    setup(){
        const userStore = useUserInformation();
        return {userStore};
    },
    mounted() {
        this.userStore.fetchCompanions();
    },
    methods: {
        async createGroup(){
            if (this.groupName){
                try{
                    await axios.post('api/groups', {groupName: this.groupName, groupMembers: this.groupMembers});
                    console.log(this.groupName + 'is created succesfully');
                    this.groupName = '';
                    this.groupMembers = [];
                    this.userStore.fetchGroups();
                }catch(error){console.log(error)}
            }else{console.log('Please type a group name')}
        },
    }
}
</script>

<style scoped>
.create-group {
   padding: 20px;
   max-width: 500px;
   margin: 0 auto;
}

.group-form {
   display: flex;
   flex-direction: column;
   gap: 20px;
}

.group-input {
   padding: 12px;
   border: 1px solid #ddd;
   border-radius: 8px;
   font-size: 16px;
}

.members-section {
   background-color: #f5f5f5;
   padding: 20px;
   border-radius: 8px;
}

.members-grid {
   display: grid;
   gap: 12px;
   grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
   text-align: center;
}

.member-item {
   background-color: white;
   padding: 10px;
   border-radius: 6px;
   text-align: center;
}

.checkbox-label {
   display: flex;
   align-items: center;
   gap: 8px;
   cursor: pointer;
   justify-content: center;
}

.create-button {
   padding: 12px 24px;
   background-color: #09BC8A;
   color: white;
   border: none;
   border-radius: 8px;
   cursor: pointer;
   transition: transform 0.2s;
}

.create-button:hover {
   transform: scale(1.02);
}
</style>