<template>
    <div v-if="userStore.conversationId">
        <Show/>
        <br>
        <Type/>
    </div>
</template>

<script>
import Show from './Show.vue';
import Type from './Type.vue';
import { useUserInformation } from '@/stores/myStore';
export default {
    name: 'Messages',
    components: {
        Show,
        Type,
    },
    setup(){
        const userStore = useUserInformation();
        return {userStore};
    },
}
</script>