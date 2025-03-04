<template>
  <div class="page-wrapper">
    <div class="page-header">
      <div class="header-content">
        <div class="welcome-container">
          <h1 class="welcome-title">
            <svg class="feather header-icon"><use href="/feather-sprite-v4.29.0.svg#plus-circle"/></svg>
            Create New Group
          </h1>
          <p class="welcome-subtitle">Start a new community with your friends</p>
        </div>
      </div>
    </div>
    
    <div class="create-group-content">
      <div class="card">
        <div class="card-header">
          <h2 class="card-title">Group Details</h2>
          <p class="card-subtitle">Fill out the group information</p>
        </div>
        
        <div class="card-body">
          <div class="form-section">
            <div class="form-group">
              <label for="group-name" class="form-label">Group Name</label>
              <div class="input-wrapper">
                <span class="input-icon">
                  <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#users"/></svg>
                </span>
                <input
                  id="group-name"
                  v-model="groupName"
                  class="form-input"
                  type="text"
                  placeholder="Enter a name for your group"
                />
              </div>
            </div>
            
            <div class="form-group">
              <label for="group-image" class="form-label">Group Photo</label>
              <div class="file-upload-container">
                <div class="file-upload-wrapper">
                  <input
                    id="group-image"
                    ref="fileInput"
                    type="file"
                    @change="handleImageUpload"
                    accept="image/*"
                    class="file-input"
                  />
                  <div class="file-upload-button">
                    <svg class="feather upload-icon"><use href="/feather-sprite-v4.29.0.svg#image"/></svg>
                    <span>Choose Image</span>
                  </div>
                  <span class="file-name">{{ fileName || 'No file chosen' }}</span>
                </div>
              </div>
              
              <div v-if="previewImage" class="image-preview-container">
                <img :src="previewImage" class="preview-image" alt="Group Image Preview"/>
              </div>
            </div>
          </div>
          
          <div class="member-section">
            <h3 class="section-title">
              <svg class="feather section-icon"><use href="/feather-sprite-v4.29.0.svg#users"/></svg>
              Add Members
            </h3>
            
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
            
            <div v-if="error" class="error-box">
              {{ error }}
            </div>
            
            <LoadingSpinner v-if="loading" :loading="loading" />
            
            <div v-if="!loading && showResults" class="search-results">
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
                    :class="isUserAdded(user) ? 'accent-btn' : 'secondary-btn'"
                    @click="addUserToGroup(user)"
                    :disabled="isUserAdded(user)"
                  >
                    <svg class="feather">
                      <use :href="isUserAdded(user) ? '/feather-sprite-v4.29.0.svg#check' : '/feather-sprite-v4.29.0.svg#plus'"/>
                    </svg>
                  </button>
                </div>
              </div>
            </div>
            
            <div v-if="selectedUsers.length > 0" class="selected-members">
              <div class="selected-header">
                <h4 class="selected-title">
                  Selected Members
                  <span class="selected-count">({{ selectedUsers.length }})</span>
                </h4>
              </div>
              
              <div class="selected-users">
                <div v-for="user in selectedUsers" :key="user.id" class="selected-user">
                  <div class="user-info">
                    <div class="user-avatar">
                      {{ getInitials(user.name) }}
                    </div>
                    <span class="user-name">{{ user.name }}</span>
                  </div>
                  <button class="action-btn icon-btn remove-btn" @click="removeUserFromGroup(user)">
                    <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#x"/></svg>
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
        
        <div class="card-footer">
          <button class="action-btn outline-btn" @click="cancelCreation">
            Cancel
          </button>
          <button 
            class="action-btn primary-btn create-btn" 
            @click="createGroup" 
            :disabled="!canCreateGroup"
          >
            <svg class="feather btn-icon"><use href="/feather-sprite-v4.29.0.svg#check"/></svg>
            Create Group
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import axios from "../services/axios";
import LoadingSpinner from "../components/LoadingSpinner.vue";

export default {
  name: "GroupCreateView",
  components: {
    LoadingSpinner,
  },
  data() {
    const token = localStorage.getItem("token");
    if (!token) {
      this.$router.push({ path: "/" });
      return;
    }
    return {
      groupName: "",
      query: "",
      lastQuery: "",
      users: [],
      loading: false,
      showResults: false,
      error: "",
      selectedUsers: [],
      previewImage: null,
      file: null,
      fileName: "",
    };
  },
  computed: {
    canCreateGroup() {
      return (
        this.groupName.trim() !== "" &&
        this.selectedUsers.length > 0 &&
        this.file !== null
      );
    },
  },
  methods: {
    async searchUsers() {
      if (!this.query.trim()) {
        this.error = "Please enter a valid search query.";
        this.showResults = false;
        return;
      }
      this.loading = true;
      this.error = "";
      this.users = [];
      this.showResults = false;
      try {
        const response = await axios.get(`/search`, {
          params: { username: this.query },
        });
        this.users = response.data.filter(user => user.id !== localStorage.getItem("token"));
        this.lastQuery = this.query;
        this.showResults = true;
      } catch (err) {
        const status = err.response?.status;
        const reason = err.response?.data?.message || "Failed to fetch users.";
        this.error = `Status ${status}: ${reason}`;
      } finally {
        this.loading = false;
      }
    },
    addUserToGroup(user) {
      if (!this.isUserAdded(user)) {
        this.selectedUsers.push(user);
      }
    },
    isUserAdded(user) {
      return this.selectedUsers.some((u) => u.id === user.id);
    },
    removeUserFromGroup(user) {
      this.selectedUsers = this.selectedUsers.filter((u) => u.id !== user.id);
    },
    handleImageUpload(event) {
      const file = event.target.files[0];
      if (file) {
        this.file = file;
        this.fileName = file.name;
        const reader = new FileReader();
        reader.onload = (e) => {
          this.previewImage = e.target.result;
        };
        reader.readAsDataURL(file);
      } else {
        this.previewImage = null;
        this.file = null;
        this.fileName = "";
      }
    },
    getInitials(name) {
      if (!name) return '?';
      return name.split(' ').map(part => part.charAt(0)).join('').toUpperCase().substring(0, 2);
    },
    async createGroup() {
      if (!this.canCreateGroup) {
        alert("Please fill in all required fields and select a group image.");
        return;
      }
      this.loading = true;
      this.error = "";
      const formData = new FormData();
      formData.append("name", this.groupName);
      formData.append("image", this.file);
      formData.append("members", JSON.stringify([...this.selectedUsers.map(u => u.id), localStorage.getItem("token")]));
      try {
        await axios.post(`/groups`, formData, {
          headers: {
            'Content-Type': 'multipart/form-data'
          }
        });
        alert("Group created successfully!");
        this.$router.push(`/groups`);
      } catch (err) {
        const status = err.response?.status;
        const reason = err.response?.data?.message || "Failed to create group.";
        this.error = `Status ${status}: ${reason}`;
      } finally {
        this.loading = false;
      }
    },
    cancelCreation() {
      this.$router.push('/groups');
    }
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

.create-group-content {
  margin-top: var(--gap-lg);
}

.card {
  max-width: 800px;
  margin: 0 auto;
}

.form-section {
  margin-bottom: var(--gap-lg);
  padding-bottom: var(--gap-lg);
  border-bottom: 1px solid var(--slate-200);
}

.form-group {
  margin-bottom: var(--gap-md);
}

.form-label {
  display: block;
  margin-bottom: var(--gap-xs);
  font-weight: 500;
  color: var(--text-secondary);
}

.input-wrapper {
  position: relative;
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

.file-upload-container {
  margin-top: var(--gap-xs);
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
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 0.6rem 1rem;
  background-color: var(--slate-100);
  border-right: 1px solid var(--slate-300);
  font-weight: 500;
  color: var(--slate-700);
  white-space: nowrap;
}

.upload-icon {
  width: 16px;
  height: 16px;
}

.file-name {
  padding: 0 var(--gap);
  color: var(--text-secondary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.image-preview-container {
  margin-top: var(--gap);
  border-radius: var(--border-radius);
  overflow: hidden;
  border: 1px solid var(--slate-200);
  display: inline-block;
}

.preview-image {
  max-width: 200px;
  max-height: 200px;
  display: block;
}

.member-section {
  margin-bottom: var(--gap);
}

.section-title {
  font-size: 1.25rem;
  font-weight: 600;
  margin-bottom: var(--gap);
  color: var(--text-primary);
  display: flex;
  align-items: center;
  gap: var(--gap-xs);
}

.section-icon {
  color: var(--secondary-color);
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

.error-box {
  background-color: rgba(239, 68, 68, 0.1);
  color: var(--error-color);
  border-radius: var(--border-radius);
  padding: var(--gap);
  margin: var(--gap) 0;
  border: 1px solid rgba(239, 68, 68, 0.2);
}

.search-results {
  margin-top: var(--gap);
  margin-bottom: var(--gap-lg);
  border: 1px solid var(--slate-200);
  border-radius: var(--border-radius);
  overflow: hidden;
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
  padding: var(--gap);
  text-align: center;
  min-height: auto;
}

.empty-icon-sm {
  width: 24px;
  height: 24px;
  color: var(--text-muted);
  margin-bottom: var(--gap-xs);
}

.user-list {
  max-height: 250px;
  overflow-y: auto;
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

.selected-members {
  margin-top: var(--gap-md);
  border: 1px solid var(--slate-200);
  border-radius: var(--border-radius);
  overflow: hidden;
}

.selected-header {
  padding: var(--gap-sm) var(--gap);
  background-color: var(--slate-50);
  border-bottom: 1px solid var(--slate-200);
}

.selected-title {
  font-size: 1rem;
  font-weight: 600;
  margin: 0;
  color: var(--text-primary);
  display: flex;
  align-items: center;
  gap: var(--gap-xs);
}

.selected-count {
  color: var(--text-secondary);
  font-weight: normal;
}

.selected-users {
  display: flex;
  flex-wrap: wrap;
  gap: var(--gap-xs);
  padding: var(--gap);
}

.selected-user {
  display: flex;
  align-items: center;
  gap: var(--gap-xs);
  background-color: var(--slate-100);
  border-radius: var(--border-radius);
  padding: 0.25rem 0.25rem 0.25rem 0.5rem;
  border: 1px solid var(--slate-200);
}

.remove-btn {
  padding: 2px;
  min-width: 0;
  width: 24px;
  height: 24px;
  color: var(--text-secondary);
  background-color: transparent;
}

.remove-btn:hover {
  color: var(--error-color);
  background-color: rgba(239, 68, 68, 0.1);
}

.card-footer {
  display: flex;
  justify-content: flex-end;
  gap: var(--gap);
}

.create-btn {
  min-width: 140px;
}

@media (max-width: 767px) {
  .page-wrapper {
    padding: var(--gap);
  }
  
  .search-form {
    flex-direction: column;
    gap: var(--gap-sm);
  }
  
  .card-footer {
    flex-direction: column-reverse;
    gap: var(--gap-sm);
  }
  
  .card-footer .action-btn {
    width: 100%;
  }
}
</style>