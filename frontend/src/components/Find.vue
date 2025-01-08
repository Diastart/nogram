<template>
    <div class="search-container">
        <div class="search-box">
            <input 
                v-model="companion" 
                placeholder="Companion name..." 
                class="search-input"
                required
            >
            <button @click="findCompanion" class="add-button">
                Add
            </button>
        </div>
    </div>
</template>

<script>
import axios from 'axios'
import { useUserInformation } from '@/stores/myStore';
export default {
    name : 'Find',
    setup(){
        const userStore = useUserInformation();
        return {userStore};
    },
    data(){
        return {
            companion: '',
        }
    },
    methods: {
        async findCompanion(){
            try {
                await axios.get('api/dialogs', {params: {companion: this.companion}});
                console.log(this.companion + ' is added successfully');
                this.companion = '';
                await this.userStore.fetchCompanions();
                this.userStore.candidates = this.userStore.companions.filter(c => !this.userStore.members.some(m => m.id === c.id))
            }catch(error){console.log(error)}
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