<template>
	<div
		class="flex flex-col items-center justify-center min-h-screen p-6 bg-gray-50"
	>
		<div
			class="w-full max-w-sm bg-white rounded-2xl shadow-md p-6 space-y-6 text-center"
		>
			<h1 class="text-xl font-semibold text-gray-800">
				Welcome to Chirp
			</h1>
			<p class="text-gray-600">
				Join or create chat rooms that vanish after a countdown. Say
				what you want, then let it disappear.
			</p>
			<RouterLink
				to="/home"
				class="inline-block w-full bg-blue-500 text-white font-medium py-2 rounded-md hover:bg-blue-600 transition-colors"
			>
				Try it
			</RouterLink>
		</div>
	</div>
</template>
