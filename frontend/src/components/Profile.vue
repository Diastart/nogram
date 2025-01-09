<template>
    <div class="profile-container">
        <div class="profile-content">
            <div class="photo-section">
                <div class="photo-wrapper">
                    <img :src="userStore.profilePhoto" alt="Profile" class="profile-circle">
                    <button @click="$refs.photoInput.click()" class="change-button">
                        Change Photo
                    </button>
                    <input 
                        type="file" 
                        ref="photoInput" 
                        @change="setProfilePhoto" 
                        accept="image/*" 
                        style="display: none"
                    >
                </div>
            </div>

            <div class="info-section">
                <div class="name-display">
                    <h1>{{ userStore.profileName }}</h1>
                </div>
                <div class="name-update">
                    <input 
                        v-model="newName" 
                        placeholder="New username..." 
                        class="username-input"
                    >
                    <button @click="setProfileName" class="update-button">
                        Update
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import axios from 'axios'
import { useUserInformation } from '@/stores/myStore';
export default {
    name: 'Profile',
    setup(){
        const userStore = useUserInformation();
        return {userStore};
    },
    data() {
        return {
            newName: '',
        }
    },
    async mounted() {
        await this.userStore.fetchProfile();
    },
    methods: {
        async setProfileName(){
            if (!this.newName) return;
            try {
                await axios.put('api/profile/username', { username: this.newName });
                await this.userStore.fetchProfile();
                this.newName = '';
            } catch (error) {console.log(error)}
        },
        async setProfilePhoto(event){
            const file = event.target.files[0];
            if (!file) return;

            const formData = new FormData();
            formData.append('photo', file);

            try {
                await axios.put('api/profile/photo', formData);
                await this.userStore.fetchProfile();
            } catch (error) {console.log(error)}
        }
    }
}
</script>

<style scoped>
.profile-container {
    display: flex;
    justify-content: center;
    align-items: center;
    padding: 2rem;
    width: 100%;
    min-height: 300px;
}

.profile-content {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 2rem;
    width: 100%;
    max-width: 800px;
}

.photo-section {
    width: 100%;
    display: flex;
    justify-content: center;
}

.photo-wrapper {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1rem;
}

.profile-circle {
    width: 120px;
    height: 120px;
    border-radius: 50%;
    border: 3px solid #09BC8A;
    object-fit: cover;
}

.info-section {
    width: 100%;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1.5rem;
}

.name-display h1 {
    color: #B118C8;
    font-size: 2.5em;
    margin: 0;
    text-align: center;
}

.name-update {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    align-items: center;
    width: 100%;
    max-width: 250px;
}

.username-input {
    padding: 10px 15px;
    border-radius: 8px;
    border: 2px solid #ddd;
    font-size: 1em;
    width: 250px;
    height: 42px;
}

.change-button, .update-button {
    padding: 10px 20px;
    background-color: #09BC8A;
    color: white;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    font-size: 1em;
    transition: opacity 0.2s ease;
    width: 250px;
}

.change-button:hover, .update-button:hover {
    opacity: 0.9;
}

/* Responsive design */
@media (min-width: 768px) {
    .profile-content {
        flex-direction: row;
        align-items: center;
        gap: 3rem;
    }

    .photo-section {
        width: auto;
    }

    .info-section {
        flex: 1;
    }
}
</style>