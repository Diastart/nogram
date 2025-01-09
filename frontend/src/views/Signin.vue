<template>
    <div class="login-container">
        <img src="../assets/iconca.png" class="logo">
        <div class="login-box">
            <input 
                v-model="username" 
                placeholder="Enter your username..." 
                class="login-input"
                required
            >
            <button @click="doLogin" class="login-button">
                🚀
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
                if (this.username) {
                    try{
                        const response = await axios.post('api/session', { username: this.username });
                        const token = response.data.token
                        console.log(token)
                        localStorage.setItem('token', token)
                        localStorage.setItem('username', this.username)
                        this.$router.push('/chat');
                    }catch (error){console.log(error)}
                }else {console.log('Please type a name')}
            },
        }
    }
</script>

<style scoped>
.login-container {
   height: 100vh;
   display: flex;
   flex-direction: column;
   align-items: center;
   justify-content: center;
   gap: 30px;
}

.logo {
   width: 150px;  /* Adjust size as needed */
   height: auto;
   margin-bottom: 20px;
}

.login-box {
   display: flex;
   gap: 10px;
   width: 100%;
   max-width: 400px;  /* Limit width on larger screens */
   padding: 0 20px;
}

.login-input {
   padding: 12px;
   border: 1px solid #ddd;
   border-radius: 8px;
   flex: 1;
}

.login-button {
   padding: 12px 24px;
   background-color: #09BC8A;
   color: white;
   border: none;
   border-radius: 8px;
   transition: transform 0.2s;
}

.login-button:hover {
   transform: scale(1.05);
}
</style>

