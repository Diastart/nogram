<template>
  <div class="profile-page">
    <div class="page-header">
      <div class="header-content">
        <div class="welcome-container">
          <h1 class="welcome-title">Your Profile</h1>
          <p class="welcome-subtitle">Manage your personal information</p>
        </div>
        <div class="action-buttons">
          <button type="button" class="btn action-btn refresh-btn" @click="refresh">
            <svg class="feather btn-icon"><use href="/feather-sprite-v4.29.0.svg#refresh-cw"/></svg>
            <span>Refresh</span>
          </button>
          <button type="button" class="btn action-btn new-group-btn" @click="newGroup">
            <svg class="feather btn-icon"><use href="/feather-sprite-v4.29.0.svg#users"/></svg>
            <span>New Group</span>
          </button>
          <button type="button" class="btn action-btn logout-btn" @click="logOut">
            <svg class="feather btn-icon"><use href="/feather-sprite-v4.29.0.svg#log-out"/></svg>
            <span>Sign Out</span>
          </button>
        </div>
      </div>
    </div>
    
    <div class="profile-content-wrapper">
      <div class="profile-card">
        <div class="profile-card-header">
          <div class="profile-image-container">
            <div class="profile-photo-wrapper">
              <img v-if="userPhoto" :src="userPhoto" alt="User Photo" class="profile-photo" />
              <div v-else class="profile-avatar">
                {{ getInitials(userName) }}
              </div>
            </div>
          </div>
          <div class="profile-info">
            <h2 class="profile-name">{{ userName }}</h2>
            <p class="profile-id">User ID: {{ userId }}</p>
          </div>
        </div>
        
        <div class="profile-card-body">
          <div class="profile-section">
            <h3 class="section-title">
              <svg class="feather section-icon"><use href="/feather-sprite-v4.29.0.svg#user"/></svg>
              Display Name
            </h3>
            <div class="profile-form-group">
              <label for="username" class="form-label">Update your display name</label>
              <div class="input-group">
                <input
                  id="username"
                  v-model="newUserName"
                  placeholder="Enter new display name"
                  maxlength="16"
                  minlength="3"
                  class="form-control"
                />
                <button
                  class="profile-btn primary-btn"
                  @click="updateUsername"
                  :disabled="!newUserName || newUserName === userName"
                >
                  <svg class="feather btn-icon"><use href="/feather-sprite-v4.29.0.svg#check"/></svg>
                  Save
                </button>
              </div>
            </div>
          </div>
          
          <div class="profile-section">
            <h3 class="section-title">
              <svg class="feather section-icon"><use href="/feather-sprite-v4.29.0.svg#image"/></svg>
              Profile Photo
            </h3>
            <div class="profile-form-group">
              <label for="photo-upload" class="form-label">Upload a new profile photo</label>
              <div class="file-upload-container">
                <div class="file-upload-wrapper">
                  <input 
                    type="file" 
                    id="photo-upload" 
                    @change="handlePhotoUpload" 
                    accept="image/*" 
                    class="file-input"
                  />
                  <div class="file-upload-button">
                    <svg class="feather upload-icon"><use href="/feather-sprite-v4.29.0.svg#upload"/></svg>
                    <span>Choose File</span>
                  </div>
                  <span class="file-name">{{ newPhotoName || 'No file chosen' }}</span>
                </div>
                <button 
                  class="profile-btn primary-btn"
                  @click="updatePhoto" 
                  :disabled="!newPhoto"
                >
                  <svg class="feather btn-icon"><use href="/feather-sprite-v4.29.0.svg#check"/></svg>
                  Upload
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <ErrorMsg v-if="errormsg" :msg="errormsg" />
    </div>
  </div>
</template>

<script>
import axios from "../services/axios";
import ErrorMsg from "../components/ErrorMsg.vue";

export default {
  name: "ProfileView",
  components: {
    ErrorMsg,
  },
  data() {
    return {
      userName: "", 
      userPhoto: null, 
      newUserName: "", 
      newPhoto: null,
      newPhotoName: "",
      userId: localStorage.getItem("token") || "Unknown",
      errormsg: null, 
    };
  },
  methods: {
    async fetchUserProfile() {
      try {
        const token = localStorage.getItem("token");
        if (!token) {
          this.$router.push({ path: "/" });
          return;
        }
        const response = await axios.get("/users/photo", {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });
        const { photo } = response.data;
        this.userName = localStorage.getItem("name");
        this.userPhoto = photo ? `data:image/jpeg;base64,${photo}` : null;
      } catch (error) {
        console.error("Failed to fetch user profile:", error);
        this.errormsg = "Failed to load user profile. Please try again later.";
      }
    },
    handlePhotoUpload(event) {
      const file = event.target.files[0];
      if (file) {
        this.newPhoto = file;
        this.newPhotoName = file.name;
      }
    },
    getInitials(name) {
      if (!name) return '?';
      return name.split(' ').map(part => part.charAt(0)).join('').toUpperCase().substring(0, 2);
    },
    async updatePhoto() {
      if (!this.newPhoto) return;
      try {
        const token = localStorage.getItem("token");
        const formData = new FormData();
        formData.append("photo", this.newPhoto);
        await axios.put("/users/photo", formData, {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });
        alert("Photo updated successfully!");
        this.newPhoto = null;
        this.fetchUserProfile(); 
      } catch (error) {
        console.error("Failed to update photo:", error);
        this.errormsg = "Failed to update photo. Please try again.";
      }
    },
    async updateUsername() {
      if (!this.newUserName || this.newUserName === this.userName) return;
      try {
        const token = localStorage.getItem("token");
        const response = await axios.put(
          "/users/name",
          { name: this.newUserName },
          {
            headers: {
              Authorization: `Bearer ${token}`,
            },
          }
        );
        alert("Username updated successfully!");
        localStorage.setItem("name", this.newUserName);
        this.userName = response.data.name;
        this.newUserName = response.data.name;
      } catch (error) {
        console.error("Failed to update username:", error);
        this.errormsg = "Failed to update username. Please try again.";
      }
    },
    refresh() {
      this.fetchUserProfile();
    },
    logOut() {
      localStorage.clear();
      this.$router.push({ path: "/" });
    },
    newGroup() {
      this.$router.push({ path: "/new-group" });
    }
  },
  mounted() {
    this.fetchUserProfile();
  },
};
</script>

<style scoped>
.profile-page {
  padding: 1.5rem 0;
}

.page-header {
  margin-bottom: 2rem;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 1rem;
}

.welcome-container {
  margin-bottom: 1rem;
}

.welcome-title {
  font-size: 1.75rem;
  font-weight: 700;
  margin: 0;
  color: var(--text-primary);
}

.welcome-subtitle {
  color: var(--text-secondary);
  margin: 0.25rem 0 0;
  font-size: 1rem;
}

.action-buttons {
  display: flex;
  gap: 0.75rem;
  flex-wrap: wrap;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 0.5rem 1rem;
  border-radius: var(--border-radius);
  font-weight: 500;
  transition: all 0.3s ease;
}

.btn-icon {
  width: 16px;
  height: 16px;
}

.refresh-btn {
  background-color: var(--white);
  color: var(--gray-700);
  border: 1px solid var(--gray-300);
}

.refresh-btn:hover {
  background-color: var(--gray-100);
}

.new-group-btn {
  background-color: var(--primary-color);
  color: white;
  border: none;
}

.new-group-btn:hover {
  background-color: var(--primary-dark);
}

.logout-btn {
  background-color: var(--white);
  color: var(--gray-700);
  border: 1px solid var(--gray-300);
}

.logout-btn:hover {
  background-color: var(--gray-100);
}

.profile-content-wrapper {
  max-width: 800px;
  margin: 0 auto;
}

.profile-card {
  background-color: var(--white);
  border-radius: var(--border-radius);
  box-shadow: var(--shadow-md);
  overflow: hidden;
  margin-bottom: 2rem;
}

.profile-card-header {
  display: flex;
  align-items: center;
  padding: 2rem;
  background: linear-gradient(135deg, rgba(108, 92, 231, 0.1), rgba(0, 206, 201, 0.1));
  border-bottom: 1px solid var(--gray-200);
}

.profile-image-container {
  margin-right: 2rem;
}

.profile-photo-wrapper {
  width: 120px;
  height: 120px;
  border-radius: 50%;
  overflow: hidden;
  box-shadow: var(--shadow-md);
  border: 4px solid white;
}

.profile-photo {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.profile-avatar {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: var(--primary-color);
  color: white;
  font-size: 2.5rem;
  font-weight: 600;
}

.profile-info {
  flex: 1;
}

.profile-name {
  font-size: 1.75rem;
  font-weight: 700;
  margin: 0 0 0.5rem;
  color: var(--text-primary);
}

.profile-id {
  color: var(--text-secondary);
  margin: 0;
  font-size: 0.9rem;
}

.profile-card-body {
  padding: 2rem;
}

.profile-section {
  margin-bottom: 2rem;
  padding-bottom: 2rem;
  border-bottom: 1px solid var(--gray-200);
}

.profile-section:last-child {
  margin-bottom: 0;
  padding-bottom: 0;
  border-bottom: none;
}

.section-title {
  font-size: 1.25rem;
  font-weight: 600;
  margin-bottom: 1.25rem;
  color: var(--text-primary);
  display: flex;
  align-items: center;
  gap: 10px;
}

.section-icon {
  width: 20px;
  height: 20px;
  color: var(--primary-color);
}

.profile-form-group {
  margin-bottom: 1rem;
}

.form-label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 500;
  color: var(--text-secondary);
}

.input-group {
  display: flex;
  gap: 0.75rem;
}

.form-control {
  flex: 1;
  padding: 0.75rem 1rem;
  border: 1px solid var(--gray-300);
  border-radius: var(--border-radius);
  font-size: 1rem;
  transition: all 0.3s ease;
}

.form-control:focus {
  border-color: var(--primary-color);
  box-shadow: 0 0 0 0.25rem rgba(108, 92, 231, 0.25);
}

.profile-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 0.75rem 1.5rem;
  border-radius: var(--border-radius);
  font-weight: 600;
  font-size: 0.95rem;
  cursor: pointer;
  transition: all 0.3s ease;
}

.primary-btn {
  background-color: var(--primary-color);
  color: white;
  border: none;
}

.primary-btn:hover:not(:disabled) {
  background-color: var(--primary-dark);
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.primary-btn:disabled {
  background-color: var(--gray-300);
  color: var(--gray-500);
  cursor: not-allowed;
}

.file-upload-container {
  display: flex;
  gap: 0.75rem;
}

.file-upload-wrapper {
  flex: 1;
  position: relative;
  display: flex;
  align-items: center;
  border: 1px solid var(--gray-300);
  border-radius: var(--border-radius);
  overflow: hidden;
}

.file-input {
  position: absolute;
  top: 0;
  left: 0;
  opacity: 0;
  width: 100%;
  height: 100%;
  cursor: pointer;
  z-index: 2;
}

.file-upload-button {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 0.75rem 1rem;
  background-color: var(--gray-100);
  border-right: 1px solid var(--gray-300);
  font-weight: 500;
  color: var(--gray-700);
}

.upload-icon {
  width: 16px;
  height: 16px;
}

.file-name {
  padding: 0 1rem;
  color: var(--text-secondary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

@media (max-width: 768px) {
  .profile-card-header {
    flex-direction: column;
    text-align: center;
    padding: 1.5rem;
  }
  
  .profile-image-container {
    margin-right: 0;
    margin-bottom: 1.5rem;
  }
  
  .input-group,
  .file-upload-container {
    flex-direction: column;
  }
  
  .action-buttons {
    margin-top: 1rem;
  }
}

/* Utility for getting initials */
.getInitials {
  /* This is just to ensure the method is styled */
  display: none;
}
</style>