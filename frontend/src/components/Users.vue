<template>
    <div class="users-container">
        <h5>Users of the app</h5>
        <div class="users-list">
            <div v-for="user in userStore.users" 
                style="color: black;"
                 :key="user.id" 
                 class="user-item">
                {{ user.username }}
            </div>
        </div>
    </div>
 </template>

<script>
import { useUserInformation } from '@/stores/myStore';
export default {
    name: 'Users',
    setup(){
        const userStore = useUserInformation();
        return {userStore};
    },
    mounted() {
        this.userStore.fetchUsers();
    },
}
</script>

<style scoped>
.users-container {
   padding: 20px;
}

.users-list {
   display: grid;
   gap: 10px;
   grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
}

.user-item {
   padding: 12px;
   background-color: white;
   border: 1px solid #ddd;
   border-radius: 8px;
   text-align: center;
}
</style>