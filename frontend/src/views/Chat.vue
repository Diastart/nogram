<template>
    <h1 style="color: #B118C8;">Welcome {{username}} 🥳🥳🥳</h1>
    <hr><h1 style="color: #09BC8A;">ADD   COMPANION</h1></hr>
    <Find/>
    <Users/>
    <hr><h1 style="color: #09BC8A;">CREATE   GROUP</h1></hr>
    <Create/>
    <hr><h1 style="color: #09BC8A;">SELECT</h1></hr>
    <Companions/>
    <Groups/>
    <hr><h1 style="color: #09BC8A;">CHAT</h1></hr>
    <Conversation/>
    <Messages/>
</template>

<script>
import axios from 'axios'
import Find from '../components/Find.vue';
import Companions from '../components/Companions.vue'
import Conversation from '@/components/Conversation.vue';
import Messages from '@/components/Messages.vue';
import Create from '@/components/Create.vue';
import Groups from '@/components/Groups.vue';
import Users from '@/components/Users.vue';
export default {
    name: 'Chat',
    components: {
        Find,
        Companions,
        Conversation,
        Messages,
        Create,
        Groups,
        Users,
    },
    data(){
        return {
            username: localStorage.getItem('username')
        }
    },
    setup(){
        const token = localStorage.getItem('token');
        axios.defaults.headers.common['Authorization'] = `Bearer ${token}`
        console.log('token is set')
    },
}
</script>