<template>
    <div class="search-container">
        <div class="search-box">
            <input 
                v-model="username" 
                placeholder="Enter your username..." 
                class="search-input"
                required
            >
            <button v-if="username" @click="doLogin" class="add-button">
                Login
            </button>
        </div>
    </div>
 </template>

<script>
import axios from 'axios'
import { useUserInformation } from '@/stores/myStore';
    export default {
        name: 'Signin',
        data(){
            return {
                username: '',
            }
        },
        setup(){
            const userStore = useUserInformation();
            return {userStore};
        },
        methods: {
            async doLogin(){
                try{
                    const response = await axios.post('api/session', { username: this.username });
                    const token = response.data.token
                    console.log(token)
                    localStorage.setItem('token', token)
                    localStorage.setItem('username', this.username)
                    this.$router.push('/chat');
                }catch (error){console.log(error)}
            },
        }
    }
</script>

<style scoped>
.search-container {
   padding: 20px;
}

.search-box {
   display: flex;
   gap: 10px;
}

.search-input {
   padding: 12px;
   border: 1px solid #ddd;
   border-radius: 8px;
   flex: 1;
}

.add-button {
   padding: 12px 24px;
   background-color: #09BC8A;
   color: white;
   border: none;
   border-radius: 8px;
   transition: transform 0.2s;
}

.add-button:hover {
   transform: scale(1.05);
}
</style>

