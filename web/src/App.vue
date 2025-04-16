<script setup></script>

<template>
	<h1>Hello App!</h1>
	<p><strong>Current route path:</strong> {{ $route.fullPath }}</p>
	<nav>
		<RouterLink to="/">Go to Home</RouterLink>
		<RouterLink to="/room">Go to Room</RouterLink>
	</nav>
	<main>
		<RouterView />
	</main>
</template>
