<template>
    <div class="messages-container">
        <h3>Messages:</h3>
        <p>-----------------------------------------------------------------</p>
        <div v-for="message in userStore.messages">
            <p>{{ message.senderName }} : {{ message.content }}</p>
            <p>{{ message.time }}</p>
            <p>additionally sender id: {{ message.senderId }} and id of message: {{ message.id }}</p>
            <p>-----------------------------------------------------------------</p>
            </div>
    </div>
</template>

<script>
import { useUserInformation } from '@/stores/myStore';
export default {
    name: 'Show',
    setup(){
        const userStore = useUserInformation();
        return {userStore};
    },
    mounted() {
        this.userStore.fetchMessages();
    }
}
</script>

<style scoped>
.messages-container {
    border: 1px solid #808080; 
    padding: 20px;  
    margin: 10px;  
    border-radius: 4px;
}
</style>