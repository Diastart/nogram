<template>
  <div class="page-wrapper">
    <div class="page-header">
      <div class="header-content">
        <div class="welcome-container">
          <h1 class="welcome-title">
            <svg class="feather header-icon"><use href="/feather-sprite-v4.29.0.svg#message-circle"/></svg>
            Hello, {{ username }}!
          </h1>
          <p class="welcome-subtitle">Welcome to your chat dashboard</p>
        </div>
        <div class="action-buttons">
          <button type="button" class="action-btn outline-btn" @click="refresh">
            <svg class="feather btn-icon"><use href="/feather-sprite-v4.29.0.svg#refresh-cw"/></svg>
            <span>Refresh</span>
          </button>
          <button type="button" class="action-btn secondary-btn" @click="newGroup">
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
    
    <ErrorMsg v-if="errormsg" :msg="errormsg" />
    
    <div class="conversations-section">
      <div class="section-header">
        <h2 class="section-title">Recent Conversations</h2>
        <div class="badge badge-primary">{{ conversations.length }} Total</div>
      </div>
      
      <div v-if="conversations.length === 0" class="empty-state">
        <div class="empty-icon">
          <svg class="feather empty-feather"><use href="/feather-sprite-v4.29.0.svg#message-square"/></svg>
        </div>
        <h3>No conversations yet</h3>
        <p>Start chatting with friends or create a new group</p>
      </div>
      
      <div v-else class="grid grid-auto">
        <div
          v-for="conv in conversations"
          :key="conv.id"
          class="card conversation-card"
          @click="viewConversation(conv.id, conv.name)"
        >
          <div class="conversation-content">
            <div class="conversation-avatar">
              <img
                v-if="conv.conversationPhoto.String"
                :src="'data:image/png;base64,' + conv.conversationPhoto.String"
                alt="Profile Picture"
                class="avatar-img"
              />
              <div v-else class="avatar-placeholder">
                {{ getInitials(conv.name) }}
              </div>
              <span v-if="isGroup(conv)" class="conversation-type">
                <svg class="feather type-icon"><use href="/feather-sprite-v4.29.0.svg#users"/></svg>
              </span>
            </div>
            
            <div class="conversation-details">
              <div class="conversation-header">
                <h3 class="conversation-name">{{ conv.name }}</h3>
                <div class="conversation-meta">
                  <span v-if="conv.lastMessage" class="conversation-time">
                    {{ formatTimeAgo(new Date(conv.lastMessage.timestamp)) }}
                  </span>
                </div>
              </div>
              
              <div v-if="conv.lastMessage" class="conversation-message">
                <div class="message-header">
                  <span class="sender-name">{{ conv.lastMessage.senderName }}</span>
                </div>
                <div class="message-content">
                  <div v-if="conv.lastMessage.attachment" class="attachment-preview">
                    <img 
                      :src="'data:image/*;base64,' + conv.lastMessage.attachment"
                      class="attachment-thumbnail"
                      alt="Attachment">
                    <span class="attachment-label">Image</span>
                  </div>
                  <div class="message-text">
                    <span v-if="isForwarded(conv.lastMessage)" v-html="getFormattedMessage(conv.lastMessage)"></span>
                    <span v-else>{{ getFormattedMessage(conv.lastMessage) }}</span>
                  </div>
                </div>
              </div>
              <div v-else class="no-messages">
                <span>No messages yet</span>
              </div>
            </div>
          </div>
          
          <div class="card-action">
            <button class="icon-btn action-icon">
              <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#chevron-right"/></svg>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import ErrorMsg from "../components/ErrorMsg.vue";

export default {
  name: "HomeView",
  components: {
    ErrorMsg,
  },
  data() {
    localStorage.removeItem("recipientId");
    return {
      username: "",
      errormsg: null,
      loading: false,
      conversations: [],
      pollIntervalId: null,
    };
  },
  methods: {
    async loadConversations() {
      this.errormsg = null;
      this.loading = true;
      try {
        const token = localStorage.getItem("token");
        if (!token) {
          this.$router.push({ path: "/" });
          return;
        }
        const response = await this.$axios.get("/conversations", {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });
        this.conversations = response.data || [];
      } catch (error) {
        console.error("Error loading conversations:", error);
        this.errormsg = "Failed to load conversations. Please try again.";
      } finally {
        this.loading = false;
      }
    },
    viewConversation(conversationId, conversationName) {
      localStorage.setItem("conversationName", conversationName);
      this.$router.push({
        path: `/conversations/${conversationId}`,
      });
    },
    truncateText(text, length = 50, clamp = '...') {
      if (!text || text.length <= length) {
        return text;
      }
      const lastSpaceIndex = text.substring(0, length).lastIndexOf(' ');
      if (lastSpaceIndex === -1) {
        return text.substring(0, length) + clamp;
      }
      return text.substring(0, lastSpaceIndex) + clamp;
    },
    isForwarded(message) {
      return message.content.includes("<strong>Forwarded from");
    },
    getFormattedMessage(message) {
      if (this.isForwarded(message)) {
        return message.content;
      }
      return this.truncateText(message.content);
    },
    getInitials(name) {
      if (!name) return '?';
      return name.split(' ').map(part => part.charAt(0)).join('').toUpperCase().substring(0, 2);
    },
    isGroup(conversation) {
      // This is a simple check - group conversations usually have multiple participants
      // Alternatively, you could check for a specific flag in your conversation data
      return conversation.name && conversation.name.includes('Group');
    },
    formatTimeAgo(date) {
      const now = new Date();
      const diffInSeconds = Math.floor((now - date) / 1000);
      
      if (diffInSeconds < 60) {
        return 'just now';
      }
      
      const diffInMinutes = Math.floor(diffInSeconds / 60);
      if (diffInMinutes < 60) {
        return `${diffInMinutes}m ago`;
      }
      
      const diffInHours = Math.floor(diffInMinutes / 60);
      if (diffInHours < 24) {
        return `${diffInHours}h ago`;
      }
      
      const diffInDays = Math.floor(diffInHours / 24);
      if (diffInDays < 7) {
        return `${diffInDays}d ago`;
      }
      
      return date.toLocaleDateString();
    },
    refresh() {
      this.loadConversations();
    },
    logOut() {
      this.$router.push({ path: "/" });
    },
    newGroup() {
      this.$router.push({ path: "/new-group" });
    },
  },
  mounted() {
    this.username = localStorage.getItem("name") || "Guest";
    this.loadConversations();
    this.pollIntervalId = setInterval(() => {
      this.loadConversations();
    }, 1000);
  },
  unmounted() {
    clearInterval(this.pollIntervalId);
  },
};
</script>

<style scoped>
.page-wrapper {
  padding: var(--gap-lg);
}

/* Section styling */
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

.header-icon {
  color: var(--primary-color);
  margin-right: var(--gap-xs);
}

/* Conversation cards */
.conversation-card {
  display: flex;
  align-items: stretch;
  cursor: pointer;
  padding: 0;
  overflow: hidden;
}

.conversation-content {
  flex: 1;
  display: flex;
  align-items: flex-start;
  gap: var(--gap);
  padding: var(--gap-md);
}

.conversation-avatar {
  position: relative;
  flex-shrink: 0;
  width: 56px;
  height: 56px;
}

.avatar-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: var(--border-radius-full);
  box-shadow: var(--shadow-sm);
}

.avatar-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background-image: linear-gradient(135deg, var(--primary-color), var(--secondary-color));
  color: var(--white);
  font-weight: 600;
  font-size: 1.25rem;
  border-radius: var(--border-radius-full);
}

.conversation-type {
  position: absolute;
  bottom: -2px;
  right: -2px;
  width: 20px;
  height: 20px;
  background-color: var(--secondary-color);
  color: var(--white);
  border-radius: var(--border-radius-full);
  display: flex;
  align-items: center;
  justify-content: center;
  border: 2px solid var(--white);
}

.type-icon {
  width: 12px;
  height: 12px;
}

.conversation-details {
  flex: 1;
  min-width: 0;
}

.conversation-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--gap-xs);
}

.conversation-name {
  font-size: 1rem;
  font-weight: 600;
  margin: 0;
  color: var(--text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.conversation-meta {
  display: flex;
  align-items: center;
  gap: var(--gap-xs);
}

.conversation-time {
  font-size: 0.75rem;
  color: var(--text-muted);
  white-space: nowrap;
}

.conversation-message {
  font-size: 0.875rem;
  color: var(--text-secondary);
}

.message-header {
  margin-bottom: var(--gap-xs);
}

.sender-name {
  font-weight: 500;
  color: var(--primary-color);
}

.message-content {
  display: flex;
  align-items: center;
  gap: var(--gap-xs);
  color: var(--text-muted);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.no-messages {
  font-size: 0.875rem;
  color: var(--text-light);
  font-style: italic;
}

.attachment-preview {
  display: flex;
  align-items: center;
  gap: 4px;
  background-color: rgba(6, 182, 212, 0.1);
  padding: 2px 8px;
  border-radius: var(--border-radius-full);
  color: var(--secondary-color);
  flex-shrink: 0;
}

.attachment-thumbnail {
  width: 16px;
  height: 16px;
  object-fit: cover;
  border-radius: var(--border-radius-sm);
}

.attachment-label {
  font-size: 0.75rem;
  font-weight: 500;
}

.message-text {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.card-action {
  display: flex;
  align-items: center;
  padding: 0 var(--gap-sm);
  background-color: var(--slate-50);
  border-left: 1px solid var(--slate-200);
}

.action-icon {
  width: 36px;
  height: 36px;
  color: var(--text-secondary);
  background: transparent;
  border: none;
  border-radius: var(--border-radius-full);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all var(--transition-speed);
}

.action-icon:hover {
  color: var(--primary-color);
  background-color: rgba(79, 70, 229, 0.1);
}

/* Responsive adjustments */
@media (max-width: 767px) {
  .page-wrapper {
    padding: var(--gap);
  }
  
  .conversation-content {
    padding: var(--gap);
  }
  
  .action-buttons {
    margin-top: var(--gap);
  }
}
</style>