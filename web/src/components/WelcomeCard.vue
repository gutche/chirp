<template>
	<div class="flex flex-col items-center justify-center min-h-screen p-6">
		<div
			class="w-full max-w-sm bg-white rounded-2xl shadow-md p-6 space-y-4"
		>
			<div
				v-if="!showTimerPicker"
				class="flex flex-col"
			>
				<h1 class="self-center mb-6 text-lg text-gray-700 font-medium">
					Join or create a room and enjoy!
				</h1>
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
					:class="[
						'border rounded-md px-4 py-2 focus:outline-none focus:ring-2',
						invalidUsername
							? 'border-red-500 focus:ring-red-500'
							: 'border-gray-300 focus:ring-blue-500',
					]"
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
						:disabled="room.length < 6 || !username"
						@click="connectToRoom"
						class="absolute right-1 top-1 bottom-1 px-3 bg-gray-200 text-gray-700 text-sm font-medium rounded-md hover:bg-gray-300 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
					>
						Join
					</button>
				</div>

				<hr class="my-4 border-gray-300" />

				<button
					class="w-full bg-blue-500 text-white font-medium py-2 rounded-md hover:bg-blue-600 transition-colors"
					@click="createRoom"
				>
					Create a room
				</button>
			</div>
			<div
				v-else
				class="flex flex-col items-center"
			>
				<p class="text-gray-700 mb-2">Set Room Timer</p>
				<div class="flex gap-6 w-[200px]">
					<VueScrollPicker
						v-model="hours"
						:options="hourOptions"
					/>
					<span class="self-center">h</span>
					<VueScrollPicker
						v-model="minutes"
						:options="minuteOptions"
					/>
					<span class="self-center">m</span>
				</div>

				<button
					@click="startRoomWithTimer"
					class="mt-4 bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600"
				>
					Start Room
				</button>
			</div>
		</div>
	</div>
</template>
<script setup>
	import { ref } from "vue";
	import { useRouter } from "vue-router";
	import { VueScrollPicker } from "vue-scroll-picker";
	import "vue-scroll-picker/style.css";

	const username = ref("");
	const invalidUsername = ref(false);
	const room = ref("");

	const showTimerPicker = ref(false);
	const hours = ref(0);
	const minutes = ref(5);

	const hourOptions = Array.from({ length: 24 }, (_, i) => i);

	const minuteOptions = Array.from({ length: 60 }, (_, i) => i);

	const router = useRouter();

	const createRoom = () => {
		if (!username.value.trim()) {
			flashUsernameInvalid();
			return;
		}

		showTimerPicker.value = true;
	};

	const connectToRoom = () => {
		if (!username.value.trim()) {
			flashUsernameInvalid();
			return;
		}

		if (!room.value.trim() || room.value.length < 6) {
			return;
		}

		// Optional: use a default short expiration (e.g., 5 minutes from now)
		const expiration = Math.floor(Date.now() / 1000) + 300;

		router.push({
			path: "/room",
			query: {
				id: room.value,
				username: username.value,
				expiration: expiration.toString(),
			},
		});
	};

	const startRoomWithTimer = () => {
		const totalSeconds =
			parseInt(hours.value) * 3600 + parseInt(minutes.value) * 60;

		if (totalSeconds <= 0) return;

		const now = Math.floor(Date.now() / 1000); // current time in seconds
		const expiration = now + totalSeconds;

		const id = generateRoomId();

		router.push({
			path: "/room",
			query: {
				id,
				username: username.value,
				expiration: expiration.toString(), // send as string
			},
		});
	};

	const flashUsernameInvalid = () => {
		invalidUsername.value = true;
		setTimeout(() => (invalidUsername.value = false), 1000);
	};

	const generateRoomId = () => {
		return Math.random().toString(36).substring(2, 10); // 8-char random string
	};
</script>
