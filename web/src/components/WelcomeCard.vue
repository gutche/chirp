<template>
	<div class="flex flex-col items-center justify-center min-h-screen p-6">
		<div
			class="w-full max-w-sm bg-white rounded-2xl shadow-md p-6 space-y-4"
		>
			<div class="flex flex-col">
				<img
					src="@/assets/chirp.png"
					alt="Chirp logo"
					class="w-20 mb-6 self-center"
				/>

				<label
					for="username"
					class="text-gray-700 font-medium mb-1"
					>Username</label
				>
				<input
					id="username"
					type="text"
					v-model="username"
					class="border border-gray-300 rounded-md px-4 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
				/>

				<label
					for="room"
					class="text-gray-700 font-medium mt-4 mb-1"
					>Room</label
				>
				<div class="relative">
					<input
						id="room"
						type="text"
						v-model="room"
						placeholder="room id"
						class="w-full border border-gray-300 rounded-md px-4 py-2 pr-20 focus:outline-none focus:ring-2 focus:ring-blue-500"
					/>
					<button
						:disabled="room.length < 6"
						class="absolute right-1 top-1 bottom-1 px-3 bg-gray-200 text-gray-700 text-sm font-medium rounded-md hover:bg-gray-300 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
					>
						Join
					</button>
				</div>

				<hr class="my-4 border-gray-300" />

				<button
					class="w-full bg-blue-500 text-white font-medium py-2 rounded-md hover:bg-blue-600 transition-colors"
				>
					Create a room
				</button>
			</div>
		</div>
	</div>
</template>

<script setup>
	import { ref } from "vue";

	const username = ref("");
	const room = ref("");
</script>
