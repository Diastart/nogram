<template>
    <Profile/>
    <hr><h1 style="color: #09BC8A;">USERS OF THE APP</h1></hr>
    <Users/>
    <hr><h1 style="color: #09BC8A;">ADD COMPANION</h1></hr>
    <Find/>
    <Companions/>
    <hr><h1 style="color: #09BC8A;">CREATE GROUP</h1></hr>
    <Create/>
    <Groups/>
    <hr><h1 style="color: #09BC8A;">CHAT</h1></hr>
    <Conversation/>
    <Messages/>
</template>

<script>
import axios from 'axios';
import Find from '../components/Find.vue';
import Companions from '../components/Companions.vue';
import Conversation from '@/components/Conversation.vue';
import Messages from '@/components/Messages.vue';
import Create from '@/components/Create.vue';
import Groups from '@/components/Groups.vue';
import Users from '@/components/Users.vue';
import Profile from '@/components/Profile.vue';
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
        Profile,
    },
    setup(){
        const token = localStorage.getItem('token');
        axios.defaults.headers.common['Authorization'] = `Bearer ${token}`
        console.log('token is set')
    },
}
</script>