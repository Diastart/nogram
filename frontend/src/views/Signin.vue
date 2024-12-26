<template>
    <input v-model="userStore.username" placeholder="Enter you username..." required><br>
    <button @click="doLogin">Login</button>
</template>

<script>
import axios from 'axios'
import { useUserInformation } from '@/stores/myStore';
    export default {
        name: 'Signin',
        setup(){
            const userStore = useUserInformation();
            return {userStore};
        },
        methods: {
            async doLogin(){
                try{
                    const response = await axios.post('api/session', { username: this.userStore.username });
                    const token = response.data.token
                    console.log(token)
                    localStorage.setItem('token', token)
                    axios.defaults.headers.common['Authorization'] = `Bearer ${token}`
                    this.$router.push('/chat');
                }catch (error){console.log(error)}
            },
        }
    }
</script>

