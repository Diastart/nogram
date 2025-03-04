<template>
  <div class="page-wrapper">
    <div class="page-header">
      <div class="header-content">
        <div class="welcome-container">
          <h1 class="welcome-title">
            <svg class="feather header-icon"><use href="/feather-sprite-v4.29.0.svg#settings"/></svg>
            Manage Group
          </h1>
          <p class="welcome-subtitle">Edit details and members of "{{ groupName }}"</p>
        </div>
        <div class="action-buttons">
          <button type="button" class="action-btn danger-btn" @click="leaveGroup">
            <svg class="feather btn-icon"><use href="/feather-sprite-v4.29.0.svg#log-out"/></svg>
            <span>Leave Group</span>
          </button>
        </div>
      </div>
    </div>
    
    <div class="group-edit-content">
      <div class="grid grid-2">
        <!-- Group Info Card -->
        <div class="card group-info-card">
          <div class="card-header">
            <h2 class="card-title">Group Information</h2>
            <p class="card-subtitle">Update group details</p>
          </div>
          
          <div class="card-body">
            <div class="group-photo-container">
              <div class="group-photo-wrapper">
                <img v-if="groupPhoto" :src="groupPhoto" alt="Group Photo" class="group-photo" />
                <div v-else class="group-avatar">{{ getInitials(groupName) }}</div>
              </div>
              
              <div class="group-photo-update">
                <div class="file-upload-container">
                  <div class="file-upload-wrapper">
                    <input 
                      type="file" 
                      @change="handleGroupPhotoUpload" 
                      accept="image/*" 
                      class="file-input"
                    />
                    <div class="file-upload-button">
                      <svg class="feather upload-icon"><use href="/feather-sprite-v4.29.0.svg#image"/></svg>
                      <span>Change Photo</span>
                    </div>
                  </div>
                  <button 
                    class="action-btn primary-btn photo-btn"
                    @click="updateGroupPhoto" 
                    :disabled="!newGroupPhoto"
                  >
                    <svg class="feather btn-icon"><use href="/feather-sprite-v4.29.0.svg#check"/></svg>
                    Update Photo
                  </button>
                </div>
              </div>
            </div>
            
            <div class="group-name-container">
              <label for="group-name" class="form-label">Group Name</label>
              <div class="input-group">
                <div class="input-wrapper">
                  <span class="input-icon">
                    <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#users"/></svg>
                  </span>
                  <input
                    id="group-name"
                    v-model="newGroupName"
                    class="form-input"
                    placeholder="Enter new group name"
                    maxlength="16"
                    minlength="3"
                  />
                </div>
                <button
                  class="action-btn primary-btn"
                  @click="updateGroupName"
                  :disabled="!newGroupName || newGroupName === groupName"
                >
                  <svg class="feather btn-icon"><use href="/feather-sprite-v4.29.0.svg#check"/></svg>
                  Update Name
                </button>
              </div>
            </div>
          </div>
        </div>
        
        <!-- Members Management Card -->
        <div class="card members-card">
          <div class="card-header">
            <h2 class="card-title">Members</h2>
            <p class="card-subtitle">Add new members to the group</p>
          </div>
          
          <div class="card-body">
            <div class="search-container">
              <form @submit.prevent="searchUsers" class="search-form">
                <div class="input-wrapper">
                  <span class="input-icon">
                    <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#search"/></svg>
                  </span>
                  <input
                    id="username"
                    v-model="query"
                    class="form-input"
                    type="text"
                    placeholder="Search users by name"
                  />
                </div>
                <button class="search-button action-btn secondary-btn" type="submit">
                  <svg class="feather btn-icon"><use href="/feather-sprite-v4.29.0.svg#search"/></svg>
                  Search
                </button>
              </form>
            </div>
            
            <ErrorMsg v-if="errormsg" :msg="errormsg" />
            
            <div v-if="showResults" class="search-results">
              <div class="results-header">
                <h4 class="results-title">
                  Search Results 
                  <span class="result-count">({{ users.length }})</span>
                </h4>
              </div>
              
              <div v-if="users.length === 0" class="empty-state small-empty">
                <svg class="feather empty-icon-sm"><use href="/feather-sprite-v4.29.0.svg#users"/></svg>
                <p>No users found matching "{{ lastQuery }}"</p>
              </div>
              
              <div v-else class="user-list">
                <div v-for="user in users" :key="user.id" class="user-card">
                  <div class="user-info">
                    <div class="user-avatar">
                      {{ getInitials(user.name) }}
                    </div>
                    <span class="user-name">{{ user.name }}</span>
                  </div>
                  <button
                    class="action-btn icon-btn"
                    :class="isMember(user.id) ? 'accent-btn' : 'secondary-btn'"
                    @click="handleAddToGroup(user.id)"
                    :disabled="isMember(user.id)"
                  >
                    <svg class="feather">
                      <use :href="isMember(user.id) ? '/feather-sprite-v4.29.0.svg#check' : '/feather-sprite-v4.29.0.svg#plus'"/>
                    </svg>
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import axios from "../services/axios";
import ErrorMsg from "../components/ErrorMsg.vue";

export default {
  name: "GroupEditView",
  components: {
    ErrorMsg,
  },
  data() {
    return {
      token: localStorage.getItem("token"),
      groupId: this.$route.params.uuid,
      groupName: localStorage.getItem("groupName"),
      groupPhoto: null,
      newGroupName: "", 
      newGroupPhoto: null, 
      errormsg: null, 
      query: "",
      lastQuery: "",
      users: [],
      loading: false,
      showResults: false,
      members: [],
    };
  },
  methods: {
    async fetchGroupProfile() {
      try {
        const token = localStorage.getItem("token");
        if (!token) {
            this.$router.push({ path: "/" });
            return;
        }
        const response = await axios.get(`/groups/${this.groupId}`, {
            headers: {
                Authorization: `Bearer ${token}`,
            },
        });
        const groupPhoto  = response.data.groupPhoto;
        this.groupPhoto = groupPhoto ? `data:image/*;base64,${groupPhoto}` : null;
        this.members = response.data.members;
      } catch (error) {
        console.error("Failed to fetch user profile:", error);
        this.errormsg = "Failed to load user profile. Please try again later.";
      }
    },
    getInitials(name) {
      if (!name) return '?';
      return name.split(' ').map(part => part.charAt(0)).join('').toUpperCase().substring(0, 2);
    },
    handleGroupPhotoUpload(event) {
      const file = event.target.files[0];
      if (file) {
        this.newGroupPhoto = file;
      }
    },
    async updateGroupPhoto() {
      if (!this.newGroupPhoto) return;
      try {
        const formData = new FormData();
        formData.append("photo", this.newGroupPhoto);
        await axios.put(`/groups/${this.groupId}/photo`, formData, {
          headers: {
            Authorization: `Bearer ${this.token}`,
          },
        });
        alert("Group photo updated successfully!");
        this.newGroupPhoto = null;
        this.fetchGroupProfile(); 
      } catch (error) {
        console.error("Failed to update group photo:", error);
        this.errormsg = "Failed to update group photo. Please try again.";
      }
    },
    async updateGroupName() {
      if (!this.newGroupName || this.newGroupName === this.groupName) return;
      try {
        await axios.put(
          `/groups/${this.groupId}/name`,
          { groupName: this.newGroupName },
          {
            headers: {
              Authorization: `Bearer ${this.token}`,
            },
          }
        );
        alert("Group name updated successfully!");
        localStorage.setItem("groupName", this.newGroupName) 
        this.groupName = this.newGroupName;
        this.newGroupName = "";
      } catch (error) {
        console.error("Failed to update group name:", error);
        this.errormsg = "Failed to update group name. Please try again.";
      }
    },
    async leaveGroup() {
      if (!confirm('Are you sure you want to leave this group?')) {
        return;
      }
      try {
        await axios.delete(`/groups/${this.groupId}`, {
          headers: {
            Authorization: `Bearer ${this.token}`,
          },
        });
        this.$router.push({ path: "/groups" });
      } catch (error) {
        console.error("Failed to leave group:", error);
        this.errormsg = "Failed to leave group. Please try again.";
      }
    },  
    async searchUsers() {
      if (!this.query.trim()) {
        this.errormsg = "Please enter a search query";
        this.showResults = false;
        return;
      }

      this.loading = true;
      this.errormsg = null;
      try {
        const response = await axios.get(`/search`, {
          params: { username: this.query }
        });
        this.users = response.data;
        this.lastQuery = this.query;
        this.showResults = true;
      } catch (error) {
        console.error("Search failed:", error);
        this.errormsg = "Failed to search users. Please try again.";
      } finally {
        this.loading = false;
      }
    },
    isMember(userId) {
      return this.members.includes(userId);
    },
    async handleAddToGroup(userId) {
      if (this.isMember(userId)) return;
      try {
        await axios.post(`/groups/${this.groupId}`, 
        {  
          userId: userId,
        },
        {
          headers: {
            Authorization: `Bearer ${this.token}`,
          },}
        );
        this.members.push(userId);
        this.errormsg = null;
      } catch (error) {
        console.error("Failed to add user:", error);
        this.errormsg = "Failed to add user to group. Please try again.";
      }
    },
  },
  mounted() {
    this.fetchGroupProfile();
  },
};
</script>

<style scoped>
.page-wrapper {
  padding: var(--gap-lg);
}

.header-icon {
  color: var(--secondary-color);
}

.group-edit-content {
  margin-top: var(--gap-lg);
}

.danger-btn {
  background-color: var(--error-color);
  color: var(--white);
  border: none;
}

.danger-btn:hover {
  background-color: #d92222;
}

/* Group Info Card */
.group-info-card {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.group-photo-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: var(--gap-lg);
}

.group-photo-wrapper {
  width: 120px;
  height: 120px;
  border-radius: var(--border-radius-full);
  overflow: hidden;
  box-shadow: var(--shadow-md);
  margin-bottom: var(--gap);
}

.group-photo {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.group-avatar {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background-image: var(--gradient-secondary);
  color: var(--white);
  font-weight: 600;
  font-size: 2.5rem;
}

.group-photo-update {
  width: 100%;
  max-width: 300px;
}

.file-upload-container {
  display: flex;
  flex-direction: column;
  gap: var(--gap-sm);
}

.file-upload-wrapper {
  position: relative;
  display: flex;
  align-items: center;
  border: 1px solid var(--slate-300);
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
  width: 100%;
  padding: 0.6rem 1rem;
  background-color: var(--slate-100);
  font-weight: 500;
  color: var(--slate-700);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--gap-sm);
}

.upload-icon {
  width: 16px;
  height: 16px;
}

.photo-btn {
  width: 100%;
}

.group-name-container {
  margin-top: var(--gap-md);
}

.input-group {
  display: flex;
  gap: var(--gap);
}

.input-wrapper {
  position: relative;
  flex: 1;
}

.input-icon {
  position: absolute;
  left: var(--gap);
  top: 50%;
  transform: translateY(-50%);
  color: var(--text-muted);
}

.form-input {
  width: 100%;
  height: 44px;
  padding: 0 var(--gap) 0 var(--gap-xl);
  font-size: 1rem;
  border: 1px solid var(--slate-300);
  border-radius: var(--border-radius);
  transition: all var(--transition-speed);
}

.form-input:focus {
  outline: none;
  border-color: var(--primary-color);
  box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1);
}

.form-label {
  display: block;
  margin-bottom: var(--gap-xs);
  font-weight: 500;
  color: var(--text-secondary);
}

/* Members Card */
.members-card {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.search-container {
  margin-bottom: var(--gap);
}

.search-form {
  display: flex;
  gap: var(--gap);
}

.search-form .input-wrapper {
  flex: 1;
}

.search-button {
  white-space: nowrap;
}

.search-results {
  margin-top: var(--gap);
  border: 1px solid var(--slate-200);
  border-radius: var(--border-radius);
  overflow: hidden;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.results-header {
  padding: var(--gap-sm) var(--gap);
  background-color: var(--slate-50);
  border-bottom: 1px solid var(--slate-200);
}

.results-title {
  font-size: 1rem;
  font-weight: 600;
  margin: 0;
  color: var(--text-primary);
  display: flex;
  align-items: center;
  gap: var(--gap-xs);
}

.result-count {
  color: var(--text-secondary);
  font-weight: normal;
}

.small-empty {
  padding: var(--gap-lg);
  text-align: center;
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.empty-icon-sm {
  width: 32px;
  height: 32px;
  color: var(--text-muted);
  margin-bottom: var(--gap-sm);
}

.user-list {
  max-height: 300px;
  overflow-y: auto;
  flex: 1;
}

.user-card {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--gap);
  border-bottom: 1px solid var(--slate-200);
}

.user-card:last-child {
  border-bottom: none;
}

.user-info {
  display: flex;
  align-items: center;
  gap: var(--gap);
}

.user-avatar {
  width: 36px;
  height: 36px;
  border-radius: var(--border-radius-full);
  background-color: var(--secondary-color);
  color: var(--white);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 0.9rem;
}

.user-name {
  font-weight: 500;
  color: var(--text-primary);
}

.accent-btn {
  background-color: var(--accent-color);
  color: white;
}

.accent-btn:hover {
  background-color: var(--accent-dark);
}

@media (max-width: 991px) {
  .grid-2 {
    grid-template-columns: 1fr;
    gap: var(--gap-lg);
  }
  
  .group-info-card, .members-card {
    height: auto;
  }
}

@media (max-width: 767px) {
  .page-wrapper {
    padding: var(--gap);
  }
  
  .search-form, .input-group {
    flex-direction: column;
    gap: var(--gap-sm);
  }
  
  .group-photo-update {
    max-width: 100%;
  }
  
  .user-list {
    max-height: 200px;
  }
}
</style>