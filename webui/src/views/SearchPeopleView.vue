<template>
  <div class="page-wrapper">
    <div class="page-header">
      <div class="header-content">
        <div class="welcome-container">
          <h1 class="welcome-title">
            <svg class="feather header-icon"><use href="/feather-sprite-v4.29.0.svg#search"/></svg>
            Discover People
          </h1>
          <p class="welcome-subtitle">Find and chat with other users</p>
        </div>
        <div class="action-buttons">
          <button type="button" class="action-btn outline-btn" @click="refresh">
            <svg class="feather btn-icon"><use href="/feather-sprite-v4.29.0.svg#refresh-cw"/></svg>
            <span>Refresh</span>
          </button>
          <button type="button" class="action-btn primary-btn" @click="newGroup">
            <svg class="feather btn-icon"><use href="/feather-sprite-v4.29.0.svg#users"/></svg>
            <span>New Group</span>
          </button>
          <button type="button" class="action-btn outline-btn" @click="logOut">
            <svg class="feather btn-icon"><use href="/feather-sprite-v4.29.0.svg#log-out"/></svg>
            <span>Sign Out</span>
          </button>
        </div>
      </div>
    </div>
    
    <div class="search-content">
      <div class="card">
        <div class="card-header">
          <h2 class="card-title">Find Users</h2>
          <p class="card-subtitle">Search for people to chat with</p>
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
          
          <ErrorMsg v-if="error" :msg="error" />
          
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
                  class="action-btn secondary-btn chat-btn"
                  @click="navigateToConversation(user.id, user.name)"
                >
                  <svg class="feather btn-icon"><use href="/feather-sprite-v4.29.0.svg#message-circle"/></svg>
                  Message
                </button>
              </div>
            </div>
          </div>
          
          <div v-if="!query && !showResults" class="no-search">
            <div class="empty-state">
              <div class="empty-icon">
                <svg class="feather empty-feather"><use href="/feather-sprite-v4.29.0.svg#search"/></svg>
              </div>
              <h3>Search for Users</h3>
              <p>Enter a name in the search box above to find people to chat with</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import axios from "../services/axios";
import LoadingSpinner from "../components/LoadingSpinner.vue";
import ErrorMsg from "../components/ErrorMsg.vue";

export default {
  name: "SearchPeopleView",
  components: {
    LoadingSpinner,
    ErrorMsg
  },
  data() {
    return {
      userName: localStorage.getItem("name"),
      query: "",
      lastQuery: "",
      users: [],
      loading: false,
      showResults: false,
      error: "",
    };
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
        this.users = response.data;
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
    navigateToConversation(recipientId, recipientName) {
      localStorage.setItem("conversationName", recipientName);
      const senderId = localStorage.getItem("token");
      axios
        .post(`/conversations`, { senderId, recipientId })
        .then((response) => {
          const conversationId = response.data.conversationId;
          this.$router.push({
            path: `/conversations/${conversationId}`
          });
        })
        .catch((error) => {
          console.error("Error starting conversation:", error);
          this.error = "Failed to start conversation. Please try again.";
        });
    },
    getInitials(name) {
      if (!name) return '?';
      return name.split(' ').map(part => part.charAt(0)).join('').toUpperCase().substring(0, 2);
    },
    refresh() {
      this.query = "";
      this.users = [];
      this.showResults = false;
      this.error = "";
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
    const token = localStorage.getItem("token");
    if (!token) {
      this.$router.push({ path: "/" });
    }
  }
};
</script>

<style scoped>
.page-wrapper {
  padding: var(--gap-lg);
}

.header-icon {
  color: var(--secondary-color);
}

.search-content {
  margin-top: var(--gap-lg);
  max-width: 800px;
  margin-left: auto;
  margin-right: auto;
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
  max-height: 400px;
  overflow-y: auto;
}

.user-card {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--gap);
  border-bottom: 1px solid var(--slate-200);
  transition: background-color var(--transition-speed);
}

.user-card:hover {
  background-color: var(--slate-50);
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
  width: 40px;
  height: 40px;
  border-radius: var(--border-radius-full);
  background-color: var(--secondary-color);
  color: var(--white);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 1rem;
}

.user-name {
  font-weight: 500;
  color: var(--text-primary);
}

.chat-btn {
  background-color: var(--secondary-color);
  color: white;
  padding: 0.5rem 1rem;
  border-radius: var(--border-radius);
  font-weight: 500;
  border: none;
  transition: all var(--transition-speed);
}

.chat-btn:hover {
  background-color: var(--secondary-dark);
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.chat-btn .btn-icon {
  width: 18px;
  height: 18px;
}

.no-search {
  padding: var(--gap-md) 0;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--gap-xl) var(--gap-lg);
  text-align: center;
}

.empty-icon {
  width: 80px;
  height: 80px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: rgba(79, 70, 229, 0.1);
  border-radius: var(--border-radius-full);
  margin-bottom: var(--gap-md);
}

.empty-feather {
  width: 40px;
  height: 40px;
  color: var(--primary-color);
}

.empty-state h3 {
  font-size: 1.25rem;
  font-weight: 600;
  margin-bottom: var(--gap-sm);
  color: var(--text-primary);
}

.empty-state p {
  color: var(--text-secondary);
  max-width: 350px;
  margin: 0;
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