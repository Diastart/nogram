<template>
    <h3>Welcome {{username}} !!!</h3>
    <p>-x-x-x-x-x-x-x-x-x-x-x-<span style="color: #6BD425;">SEARCH   COMPANION</span>-x-x-x-x-x-x-x-x-x-x-x-x-x-</p>
    <Find/>
    <p>-x-x-x-x-x-x-x-x-x-x-x-x-x-<span style="color: #6BD425;">CREATE   GROUP</span>-x-x-x-x-x-x-x-x-x-x-x-x-x-x-</p>
    <Create/>
    <p>-x-x-x-x-x-x-x-x-x-x-x-x-x-x-x-x-<span style="color: #6BD425;">SELECT</span>-x-x-x-x-x-x-x-x-x-x-x-x-x-x-x-x-x-</p>
    <Companions/>
    <Groups/>
    <p>-x-x-x-x-x-x-x-x-x-x-x-x-x-x-x-x-x-<span style="color: #6BD425;">CHAT</span>-x-x-x-x-x-x-x-x-x-x-x-x-x-x-x-x-x-x-</p>
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
export default {
    name: 'Chat',
    components: {
        Find,
        Companions,
        Conversation,
        Messages,
        Create,
        Groups,
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