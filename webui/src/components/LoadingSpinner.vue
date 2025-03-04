<script>
export default {
	props: ["loading"]
}
</script>

<template>
	<div v-if="loading" class="spinner-container">
		<div class="spinner-wrapper">
			<div class="custom-spinner">
				<div class="spinner-ring"></div>
				<div class="spinner-dot"></div>
			</div>
			<p class="spinner-text">Loading data...</p>
		</div>
	</div>
	<div v-if="!loading"><slot /></div>
</template>

<style scoped>
.spinner-container {
	display: flex;
	justify-content: center;
	align-items: center;
	width: 100%;
	padding: 2rem 0;
}

.spinner-wrapper {
	display: flex;
	flex-direction: column;
	align-items: center;
}

.custom-spinner {
	position: relative;
	width: 60px;
	height: 60px;
	margin-bottom: 1rem;
}

.spinner-ring {
	position: absolute;
	top: 0;
	left: 0;
	width: 100%;
	height: 100%;
	border: 3px solid rgba(108, 92, 231, 0.1);
	border-top: 3px solid var(--primary-color);
	border-radius: 50%;
	animation: spin 1s linear infinite;
}

.spinner-dot {
	position: absolute;
	top: 50%;
	left: 50%;
	width: 12px;
	height: 12px;
	background-color: var(--primary-color);
	border-radius: 50%;
	transform: translate(-50%, -50%);
	box-shadow: 0 0 10px rgba(108, 92, 231, 0.5);
}

.spinner-text {
	color: var(--text-secondary);
	font-size: 0.95rem;
	font-weight: 500;
	margin: 0;
}

@keyframes spin {
	0% { transform: rotate(0deg); }
	100% { transform: rotate(360deg); }
}
</style>
