<template>
    <div class="container">
        <input v-model="groupName" placeholder="Type group name..." required>
        <p>choose people to add</p>
        <div v-for="companion in userStore.companions" class="member-row">
            <label>{{ companion.username }}</label>
            <input v-model="groupMembers" type="checkbox" :value="companion.id">
        </div>
        <button @click="createGroup">create</button>
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
            try{
                await axios.post('api/groups', {groupName: this.groupName, groupMembers: this.groupMembers});
                console.log(this.groupName + 'is created succesfully');
                this.groupName = '';
                this.groupMembers = [];
                this.userStore.fetchGroups();
            }catch(error){console.log(error)}
        },
    }
}
</script>

<style scoped>
.container {
   display: flex;
   flex-direction: column;
   align-items: center;
}

.member-row {
   display: flex;
   justify-content: space-between;
   align-items: center;
   margin: 5px 0;
   width: 100px;
}
</style>