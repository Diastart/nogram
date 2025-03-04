
<template>
  <div class="page-wrapper">
    <div class="page-header">
      <div class="header-content">
        <div class="welcome-container">
          <h1 class="welcome-title">
            <svg class="feather header-icon"><use href="/feather-sprite-v4.29.0.svg#users"/></svg>
            My Communities
          </h1>
          <p class="welcome-subtitle">Manage and discover your group conversations</p>
        </div>
        <div class="action-buttons">
          <button type="button" class="action-btn outline-btn" @click="refresh">
            <svg class="feather btn-icon"><use href="/feather-sprite-v4.29.0.svg#refresh-cw"/></svg>
            <span>Refresh</span>
          </button>
          <button type="button" class="action-btn primary-btn" @click="newGroup">
            <svg class="feather btn-icon"><use href="/feather-sprite-v4.29.0.svg#plus"/></svg>
            <span>Create Group</span>
          </button>
          <button type="button" class="action-btn outline-btn" @click="logOut">
            <svg class="feather btn-icon"><use href="/feather-sprite-v4.29.0.svg#log-out"/></svg>
            <span>Sign Out</span>
          </button>
        </div>
      </div>
    </div>
    
    <ErrorMsg v-if="errormsg" :msg="errormsg" />
    
    <div class="groups-section">
      <LoadingSpinner v-if="loading" :loading="loading" />
      
      <div v-else-if="groups.length === 0" class="empty-state">
        <div class="empty-icon">
          <svg class="feather empty-feather"><use href="/feather-sprite-v4.29.0.svg#users"/></svg>
        </div>
        <h3>No Groups Yet</h3>
        <p>Create a new group or join existing ones to start collaborating.</p>
        <button class="action-btn primary-btn empty-action" @click="newGroup">
          <svg class="feather btn-icon"><use href="/feather-sprite-v4.29.0.svg#plus"/></svg>
          <span>Create Your First Group</span>
        </button>
      </div>
      
      <div v-else>
        <div class="section-header">
          <h2 class="section-title">Your Groups</h2>
          <div class="badge badge-secondary">{{ groups.length }} Groups</div>
        </div>
        
        <div class="grid grid-auto">
          <div
            v-for="group in groups"
            :key="group.id"
            class="card group-card"
            @click="viewGroup(group.id, group.name)"
          >
            <div class="group-cover-image" :style="groupCoverStyle(group)">
              <div class="group-badge">
                <svg class="feather badge-icon"><use href="/feather-sprite-v4.29.0.svg#users"/></svg>
              </div>
            </div>
            
            <div class="card-body group-info">
              <div class="group-avatar">
                <img
                  v-if="group.conversationPhoto && group.conversationPhoto.String"
                  :src="'data:image/png;base64,' + group.conversationPhoto.String"
                  alt="Group"
                  class="avatar-img"
                />
                <div v-else class="avatar-placeholder">
                  {{ getInitials(group.name) }}
                </div>
              </div>
              
              <div class="group-details">
                <h3 class="group-name">{{ group.name }}</h3>
                <div class="group-meta">
                  <span class="group-members">{{ group.membersCount || 'Unknown' }} members</span>
                  <span class="group-date">Created {{ formatDate(group.createdAt) }}</span>
                </div>
              </div>
            </div>
            
            <div class="card-footer group-action">
              <button class="icon-btn view-group-btn">
                <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#chevron-right"/></svg>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import ErrorMsg from "../components/ErrorMsg.vue";
import LoadingSpinner from "../components/LoadingSpinner.vue";

export default {
  components: {
    ErrorMsg,
    LoadingSpinner,
  },
  data() {
    return {
      username: "",
      errormsg: null,
      loading: false,
      groups: []
    };
  },
  methods: {
    async loadGroups() {
      this.errormsg = null;
      this.loading = true;
      try {
        const token = localStorage.getItem("token");
        if (!token) {
          this.$router.push({ path: "/" });
          return;
        }
        const response = await this.$axios.get("/groups", {
          headers: {
            Authorization: `Bearer ${token}`
          }
        });
        this.groups = response.data || [];
      } catch (error) {
        console.error("Error loading groups:", error);
        this.errormsg = "Failed to load groups. Please try again.";
      } finally {
        this.loading = false;
      }
    },
    viewGroup(groupId, groupName) {
      localStorage.setItem("groupName", groupName);
      this.$router.push({
        path: `/groups/${groupId}`
      });
    },
    refresh() {
      this.loadGroups();
    },
    logOut() {
      this.$router.push({ path: "/" });
    },
    newGroup() {
      this.$router.push({ path: "/new-group" });
    },
    getInitials(name) {
      if (!name) return '?';
      return name.split(' ').map(part => part.charAt(0)).join('').toUpperCase().substring(0, 2);
    },
    formatDate(dateString) {
      if (!dateString) return 'recently';
      
      const date = new Date(dateString);
      const now = new Date();
      const diffTime = Math.abs(now - date);
      const diffDays = Math.floor(diffTime / (1000 * 60 * 60 * 24));
      
      if (diffDays < 1) return 'today';
      if (diffDays < 2) return 'yesterday';
      if (diffDays < 7) return `${diffDays} days ago`;
      if (diffDays < 30) return `${Math.floor(diffDays / 7)} weeks ago`;
      
      return date.toLocaleDateString();
    },
    groupCoverStyle(group) {
      // Generate a deterministic gradient based on the group name
      const hash = this.hashString(group.name);
      const hue1 = hash % 360;
      const hue2 = (hue1 + 40) % 360;
      return {
        backgroundImage: `linear-gradient(135deg, hsl(${hue1}, 80%, 55%), hsl(${hue2}, 80%, 45%))`
      };
    },
    hashString(str) {
      let hash = 0;
      for (let i = 0; i < str.length; i++) {
        hash = ((hash << 5) - hash) + str.charCodeAt(i);
        hash = hash & hash; // Convert to 32bit integer
      }
      return Math.abs(hash);
    }
  },
  mounted() {
    this.username = localStorage.getItem("name") || "Guest";
    this.loadGroups();
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

.groups-section {
  margin-top: var(--gap-lg);
}

.section-header {
  display: flex;
  align-items: center;
  gap: var(--gap);
  margin-bottom: var(--gap-md);
}

.section-title {
  font-size: 1.25rem;
  font-weight: 600;
  margin: 0;
  color: var(--text-primary);
}

.empty-action {
  margin-top: var(--gap-md);
}

/* Group cards */
.group-card {
  overflow: hidden;
  transition: all var(--transition-speed);
  cursor: pointer;
  display: flex;
  flex-direction: column;
  height: 100%;
}

.group-card:hover {
  transform: translateY(-5px);
  box-shadow: var(--shadow-lg);
}

.group-cover-image {
  height: 80px;
  position: relative;
}

.group-badge {
  position: absolute;
  top: var(--gap-sm);
  right: var(--gap-sm);
  width: 24px;
  height: 24px;
  border-radius: var(--border-radius-full);
  background-color: rgba(255, 255, 255, 0.3);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
}

.badge-icon {
  width: 14px;
  height: 14px;
  color: var(--white);
}

.group-info {
  padding: var(--gap);
  display: flex;
  align-items: center;
  gap: var(--gap);
  flex: 1;
}

.group-avatar {
  width: 50px;
  height: 50px;
  margin-top: -40px;
  border-radius: var(--border-radius-full);
  overflow: hidden;
  border: 3px solid var(--white);
  box-shadow: var(--shadow-md);
  flex-shrink: 0;
}

.avatar-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: var(--secondary-color);
  color: var(--white);
  font-weight: 600;
  font-size: 1.1rem;
}

.group-details {
  flex: 1;
}

.group-name {
  font-size: 1.1rem;
  font-weight: 600;
  margin: 0 0 var(--gap-xs) 0;
  color: var(--text-primary);
}

.group-meta {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.group-members, .group-date {
  font-size: 0.8rem;
  color: var(--text-secondary);
}

.group-action {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  padding: var(--gap-sm) var(--gap);
  background-color: var(--slate-50);
  border-top: 1px solid var(--slate-200);
}

.view-group-btn {
  background-color: var(--white);
  border: 1px solid var(--slate-200);
  color: var(--secondary-color);
}

.view-group-btn:hover {
  background-color: var(--secondary-color);
  color: var(--white);
}

@media (max-width: 767px) {
  .page-wrapper {
    padding: var(--gap);
  }
}
</style>