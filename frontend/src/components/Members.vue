<template>
    <div v-if="userStore.members.length > 0" class="members-container">
        <h5>Members of the group</h5>
        <div class="members-list">
            <div v-for="member in userStore.members" 
                style="color: black;"
                :key="member.id"
                class="member-item">
                {{ member.username }}
            </div>
        </div>
    </div>
</template>

<script>
import { useUserInformation } from '@/stores/myStore';
export default {
    name: 'Members',
    setup(){
        const userStore = useUserInformation();
        return {userStore};
    },
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
</style>