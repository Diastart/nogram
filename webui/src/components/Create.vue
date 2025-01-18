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
                <h4 class="section-title">Select Members</h4>
                <div class="members-grid">
                    <div v-for="companion in userStore.companions" 
                            :key="companion.id" 
                            class="member-card"
                            :class="{ 'selected': groupMembers.includes(companion.id) }"
                            @click="toggleMember(companion.id)">
                        <span class="member-name">{{ companion.username }}</span>
                        <div class="checkbox-circle">
                            <div v-if="groupMembers.includes(companion.id)" class="check">✓</div>
                        </div>
                    </div>
                </div>
            </div>
 
            <button @click="createGroup" class="create-button">
                🚀
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
        toggleMember(id) {
            const index = this.groupMembers.indexOf(id)
            if (index === -1) {
                this.groupMembers.push(id)
            } else {
                this.groupMembers.splice(index, 1)
            }
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
   background-color: white;
   padding: 20px;
   border-radius: 12px;
   box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.section-title {
   text-align: center;
   margin-bottom: 20px;
   color: #333;
   font-size: 1.2em;
}

.members-grid {
   display: grid;
   grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
   gap: 12px;
}

.member-card {
   display: flex;
   justify-content: space-between;
   align-items: center;
   padding: 12px 16px;
   background-color: #f8f9fa;
   border: 2px solid #e9ecef;
   border-radius: 8px;
   cursor: pointer;
   transition: all 0.2s ease;
}

.member-card:hover {
   background-color: #e9ecef;
   transform: translateY(-2px);
}

.member-card.selected {
   border-color: #09BC8A;
   background-color: #e6f7f2;
}

.member-name {
   font-weight: 500;
   color: #333;
}

.checkbox-circle {
   width: 24px;
   height: 24px;
   border: 2px solid #ddd;
   border-radius: 50%;
   display: flex;
   align-items: center;
   justify-content: center;
}

.selected .checkbox-circle {
   border-color: #09BC8A;
   background-color: #09BC8A;
}

.check {
   color: white;
   font-size: 14px;
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