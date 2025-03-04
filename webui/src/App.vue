<script setup>
import { RouterLink, RouterView } from 'vue-router'
</script>
<script>
export default {}
</script>

<template>
	<header class="navbar sticky-top">
		<div class="navbar-brand-container">
			<a class="navbar-brand" href="">
				<div class="brand-logo">
					<svg class="feather brand-icon"><use href="/feather-sprite-v4.29.0.svg#message-square"/></svg>
				</div>
				<span>Messenger Hub</span>
			</a>
			<button class="navbar-toggler d-lg-none" type="button" data-bs-toggle="collapse" data-bs-target="#sidebarMenu" aria-controls="sidebarMenu" aria-expanded="false" aria-label="Toggle navigation">
				<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#menu"/></svg>
			</button>
		</div>
		
		<div class="navbar-actions d-none d-md-flex">
			<a href="#" class="nav-icon" title="Search">
				<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#search"/></svg>
			</a>
			<a href="#" class="nav-icon" title="Notifications">
				<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#bell"/></svg>
			</a>
			<a href="#" class="nav-icon" title="Settings">
				<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#settings"/></svg>
			</a>
		</div>
	</header>

	<div class="app-container">
		<nav id="sidebarMenu" class="sidebar">
			<div class="sidebar-sticky">
				<div class="sidebar-section">
					<h6 class="sidebar-heading">
						<svg class="feather section-icon"><use href="/feather-sprite-v4.29.0.svg#activity"/></svg>
						<span>Communication</span>
					</h6>
					<ul class="nav flex-column">
						<li class="nav-item">
							<RouterLink to="/home" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#message-circle"/></svg>
								<span>Chat Inbox</span>
							</RouterLink>
						</li>
						<li class="nav-item">
							<RouterLink to="/me" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#user"/></svg>
								<span>My Profile</span>
							</RouterLink>
						</li>
						<li class="nav-item">
							<RouterLink to="/search" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#search"/></svg>
								<span>Discover People</span>
							</RouterLink>
						</li>
					</ul>
				</div>

				<div class="sidebar-section">
					<h6 class="sidebar-heading">
						<svg class="feather section-icon"><use href="/feather-sprite-v4.29.0.svg#users"/></svg>
						<span>Communities</span>
					</h6>
					<ul class="nav flex-column">
						<li class="nav-item">
							<RouterLink :to="'/groups'" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#users"/></svg>
								<span>My Groups</span>
							</RouterLink>
						</li>
						<li class="nav-item">
							<RouterLink :to="'/new-group'" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#plus-circle"/></svg>
								<span>Create Group</span>
							</RouterLink>
						</li>
					</ul>
				</div>
				
				<div class="sidebar-footer">
					<div class="app-info">
						<div class="app-version">v2.0</div>
						<div class="app-copyright">© 2025 Messenger Hub</div>
					</div>
				</div>
			</div>
		</nav>

		<main class="main-content fade-in">
			<RouterView />
		</main>
	</div>
</template>

<style>
/* App layout */
.app-container {
	display: flex;
	min-height: 100vh;
}

.main-content {
	flex: 1;
}

/* Navbar styling */
.navbar {
	background: var(--gradient-primary);
	height: var(--header-height);
	padding: 0 var(--gap-lg);
	display: flex;
	align-items: center;
	justify-content: space-between;
}

.navbar-brand-container {
	display: flex;
	align-items: center;
	justify-content: space-between;
	width: var(--sidebar-width);
}

.navbar-brand {
	display: flex;
	align-items: center;
	gap: var(--gap-sm);
	text-decoration: none;
	color: var(--white);
	font-weight: 700;
	font-size: 1.25rem;
	letter-spacing: 0.5px;
}

.brand-logo {
	display: flex;
	align-items: center;
	justify-content: center;
	width: 32px;
	height: 32px;
	background-color: rgba(255, 255, 255, 0.2);
	border-radius: var(--border-radius);
}

.brand-icon {
	color: var(--white);
	width: 18px;
	height: 18px;
}

.navbar-toggler {
	background: transparent;
	border: none;
	color: var(--white);
	cursor: pointer;
	display: flex;
	align-items: center;
	justify-content: center;
	padding: 0;
}

.navbar-actions {
	display: flex;
	align-items: center;
	gap: var(--gap);
}

/* Sidebar footer */
.sidebar-footer {
	margin-top: auto;
	padding: var(--gap-lg) var(--gap-lg) var(--gap);
	border-top: 1px solid var(--slate-200);
}

.app-info {
	display: flex;
	flex-direction: column;
	gap: var(--gap-xs);
}

.app-version {
	font-size: 0.75rem;
	color: var(--primary-color);
	font-weight: 600;
}

.app-copyright {
	font-size: 0.7rem;
	color: var(--text-muted);
}

/* Section icon in sidebar */
.section-icon {
	width: 16px;
	height: 16px;
	color: var(--primary-color);
}

@media (max-width: 991.98px) {
	.navbar-brand-container {
		width: 100%;
	}
}
</style>