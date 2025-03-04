<script>
export default {
	props: ['msg']
}
</script>

<template>
	<div class="error-message" role="alert">
		<div class="error-icon-container">
			<svg class="feather error-icon"><use href="/feather-sprite-v4.29.0.svg#alert-circle"/></svg>
		</div>
		<div class="error-content">
			<span class="error-title">Error</span>
			<span class="error-text">{{ msg }}</span>
		</div>
	</div>
</template>

<style scoped>
.error-message {
	display: flex;
	align-items: flex-start;
	background-color: rgba(239, 68, 68, 0.08);
	border-radius: var(--border-radius);
	padding: 16px;
	margin: 16px 0;
	box-shadow: var(--shadow-sm);
	border: 1px solid rgba(239, 68, 68, 0.2);
	animation: fadeIn 0.3s ease;
}

.error-icon-container {
	background-color: rgba(239, 68, 68, 0.15);
	width: 32px;
	height: 32px;
	border-radius: var(--border-radius-full);
	display: flex;
	align-items: center;
	justify-content: center;
	margin-right: 16px;
	flex-shrink: 0;
}

.error-icon {
	color: var(--error-color);
	width: 18px;
	height: 18px;
}

.error-content {
	display: flex;
	flex-direction: column;
}

.error-title {
	font-weight: 600;
	color: var(--error-color);
	margin-bottom: 4px;
}

.error-text {
	color: var(--text-secondary);
	font-size: 0.9rem;
	line-height: 1.5;
}

@keyframes fadeIn {
	from { opacity: 0; transform: translateY(-10px); }
	to { opacity: 1; transform: translateY(0); }
}
</style>
