<script>
export default {
  data() {
    localStorage.clear();
    return {
      errormsg: null,
      name: "", 
      profile: {
        id: "",
        name: "",
      },
    };
  },
  methods: {
    async loadDefaultPhoto() {
      try {
          const response = await fetch('/nopfp.jpg');
          const blob = await response.blob();
          const reader = new FileReader();
          return new Promise((resolve, reject) => {
              reader.onload = () => {
                  const base64String = reader.result.toString().split(',')[1];
                  resolve(base64String);
              };
              reader.onerror = reject;
              reader.readAsDataURL(blob);
          });
      } catch (error) {
          console.error('Error loading default photo:', error);
          return null;
      }
  },
    blobToBase64(blob) {
      return new Promise((resolve) => {
        const reader = new FileReader();
        reader.onloadend = () => resolve(reader.result);
        reader.readAsDataURL(blob);
      });
    },
    async doLogin() {
      if (this.name.trim() === "") {
        this.errormsg = "Name cannot be empty.";
        return;
      }
      try {
        const photoData = await this.loadDefaultPhoto();
        const response = await this.$axios.post("/session", {
          name: this.name,
          photo: photoData,
        }, {
          headers: {
            'Content-Type': 'application/json'
          }
        });
        if (response.data.identifier) {
          this.profile.id = response.data.identifier;
          this.profile.name = this.name; 
        } else {
          throw new Error("Unexpected server response. Missing 'identifier'.");
        }
        localStorage.setItem("token", this.profile.id);
        localStorage.setItem("name", this.profile.name);
        this.$router.push({ path: "/home" });
      } catch (e) {
        if (e.response && e.response.status === 400) {
          this.errormsg =
            "Form error, please check all fields and try again.";
        } else if (e.response && e.response.status === 500) {
          this.errormsg =
            "An internal error occurred. Please try again later.";
        } else {
          this.errormsg = e.toString();
        }
      }
    },
  },
};
</script>

<template>
  <div class="login-page">
    <div class="login-wrapper">
      <div class="login-brand">
        <div class="brand-logo">
          <svg class="feather brand-icon"><use href="/feather-sprite-v4.29.0.svg#message-square"/></svg>
        </div>
        <h1 class="brand-name">Messenger Hub</h1>
      </div>
      
      <div class="login-container">
        <div class="login-header">
          <h2 class="login-title">Welcome Back</h2>
          <p class="login-subtitle">Enter your name to continue</p>
        </div>
        
        <div class="login-form">
          <div class="form-group">
            <label for="name" class="form-label">Display Name</label>
            <div class="input-wrapper">
              <span class="input-icon">
                <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#user"/></svg>
              </span>
              <input
                type="text"
                id="name"
                v-model="name"
                class="form-input"
                placeholder="Your name"
                autofocus
              />
            </div>
          </div>
          
          <button class="login-button" type="button" @click="doLogin">
            <span>Sign In</span>
            <svg class="feather login-icon"><use href="/feather-sprite-v4.29.0.svg#log-in"/></svg>
          </button>
        </div>
        
        <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
      </div>
      
      <div class="login-footer">
        <p class="copyright">© 2025 Messenger Hub. All rights reserved.</p>
        <div class="footer-links">
          <a href="#" class="footer-link">Privacy Policy</a>
          <a href="#" class="footer-link">Terms of Service</a>
        </div>
      </div>
    </div>
    
    <div class="login-decoration">
      <div class="decoration-content">
        <div class="decoration-icon">
          <svg class="feather decoration-feather"><use href="/feather-sprite-v4.29.0.svg#message-circle"/></svg>
        </div>
        <h2 class="decoration-title">Connect with everyone</h2>
        <p class="decoration-text">Stay in touch with friends, family, and colleagues through instant messaging and group conversations.</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.login-page {
  min-height: 100vh;
  display: grid;
  grid-template-columns: 1fr 1fr;
}

.login-wrapper {
  display: flex;
  flex-direction: column;
  padding: var(--gap-xl);
  background-color: var(--white);
}

.login-brand {
  display: flex;
  align-items: center;
  gap: var(--gap-sm);
  margin-bottom: var(--gap-xl);
}

.brand-logo {
  width: 40px;
  height: 40px;
  background-color: var(--primary-color);
  border-radius: var(--border-radius);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: var(--shadow-md);
}

.brand-icon {
  color: var(--white);
}

.brand-name {
  font-size: 1.5rem;
  font-weight: 700;
  margin: 0;
  background: var(--gradient-primary);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.login-container {
  max-width: 400px;
  width: 100%;
  margin: auto 0;
}

.login-header {
  margin-bottom: var(--gap-lg);
}

.login-title {
  font-size: 2rem;
  font-weight: 700;
  margin: 0 0 var(--gap-xs) 0;
  color: var(--text-primary);
}

.login-subtitle {
  color: var(--text-secondary);
  margin: 0;
  font-size: 1.1rem;
}

.login-form {
  margin-bottom: var(--gap-lg);
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
  height: 54px;
  padding: 0 var(--gap) 0 var(--gap-xl);
  font-size: 1rem;
  border: 1px solid var(--slate-300);
  border-radius: var(--border-radius);
  transition: all var(--transition-speed);
}

.form-input:focus {
  outline: none;
  border-color: var(--primary-color);
  box-shadow: 0 0 0 4px rgba(79, 70, 229, 0.1);
}

.login-button {
  width: 100%;
  height: 54px;
  background-color: var(--primary-color);
  color: var(--white);
  border: none;
  border-radius: var(--border-radius);
  font-size: 1rem;
  font-weight: 600;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--gap-sm);
  cursor: pointer;
  transition: all var(--transition-speed) var(--transition-bounce);
}

.login-button:hover {
  background-color: var(--primary-dark);
  transform: translateY(-3px);
  box-shadow: var(--shadow-md);
}

.login-button:active {
  transform: translateY(-1px);
}

.login-icon {
  width: 18px;
  height: 18px;
}

.login-footer {
  margin-top: auto;
  text-align: center;
}

.copyright {
  color: var(--text-muted);
  font-size: 0.875rem;
  margin-bottom: var(--gap-sm);
}

.footer-links {
  display: flex;
  justify-content: center;
  gap: var(--gap-md);
}

.footer-link {
  color: var(--primary-color);
  font-size: 0.875rem;
  text-decoration: none;
  transition: color var(--transition-speed);
}

.footer-link:hover {
  color: var(--primary-dark);
  text-decoration: underline;
}

.login-decoration {
  background: var(--gradient-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
}

.decoration-content {
  color: var(--white);
  text-align: center;
  max-width: 400px;
  padding: var(--gap-xl);
  position: relative;
  z-index: 1;
}

.decoration-icon {
  width: 80px;
  height: 80px;
  background-color: rgba(255, 255, 255, 0.2);
  border-radius: var(--border-radius-full);
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto var(--gap-lg);
}

.decoration-feather {
  width: 40px;
  height: 40px;
}

.decoration-title {
  font-size: 2rem;
  font-weight: 700;
  margin: 0 0 var(--gap) 0;
}

.decoration-text {
  font-size: 1.1rem;
  line-height: 1.6;
  opacity: 0.9;
}

/* Responsive design */
@media (max-width: 991px) {
  .login-page {
    grid-template-columns: 1fr;
    grid-template-rows: auto 1fr;
  }
  
  .login-decoration {
    display: none;
  }
  
  .login-wrapper {
    padding: var(--gap-lg);
  }
}

@media (max-width: 576px) {
  .login-wrapper {
    padding: var(--gap);
  }
  
  .login-container {
    max-width: 100%;
  }
  
  .login-title {
    font-size: 1.75rem;
  }
}
</style>
