<template>
    <h3>Hello my dear friend {{ userStore.username }}</h3>
    <p>-----------------------------------------------------------------</p>
    <Find/>
    <p>-----------------------------------------------------------------</p>
    <Companions/>
    <p>-----------------------------------------------------------------</p>
    <Conversation/><br>
    <Messages/>
    
</template>

<script>
import { useUserInformation } from '@/stores/myStore';
import Find from '../components/Find.vue';
import Companions from '../components/Companions.vue'
import Conversation from '@/components/Conversation.vue';
import Messages from '@/components/Messages.vue';
export default {
    name: 'Chat',
    components: {
        Find,
        Companions,
        Conversation,
        Messages,
    },
    setup(){
        const userStore = useUserInformation();
        return {userStore};
    },
}
</script>