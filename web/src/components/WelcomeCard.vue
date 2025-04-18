<template>
	<div class="flex flex-col items-center justify-center min-h-screen p-6">
		<div
			v-if="!connected"
			class="w-full max-w-sm bg-white rounded-2xl shadow-md p-6 space-y-4"
		>
			<!-- Join/Create View -->
			<div class="flex flex-col">
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
					>Username*</label
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
					:disabled="!username"
					@click="createRoom"
				>
					Create a room
				</button>
			</div>
		</div>

		<div
			v-else
			class="w-full max-w-lg bg-white rounded-2xl shadow-md p-6 flex flex-col h-[80vh]"
		>
			<!-- Chat View -->
			<h2 class="text-xl font-semibold text-gray-800 mb-4">
				Room: {{ room }}
			</h2>

			<div
				class="flex-1 overflow-y-auto border rounded p-4 mb-4 bg-gray-50 space-y-2"
			>
				<div
					v-for="(msg, index) in messages"
					:key="index"
					class="text-sm"
				>
					{{ msg }}
				</div>
			</div>

			<div class="flex gap-2">
				<input
					v-model="message"
					@keyup.enter="sendMessage"
					class="flex-1 border border-gray-300 rounded-md px-4 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
					placeholder="Type a message..."
				/>
				<button
					@click="sendMessage"
					class="bg-blue-500 text-white font-medium px-4 py-2 rounded-md hover:bg-blue-600"
				>
					Send
				</button>
			</div>
		</div>
	</div>
</template>

<script setup>
	import { ref, onBeforeUnmount } from "vue";

	const username = ref("");
	const room = ref("");
	const message = ref("");
	const messages = ref([]);
	const connected = ref(false);

	let socket = null;

	const connectToRoom = () => {
		const url = `ws://localhost:8080/ws?room=${room.value}&username=${username.value}`;
		socket = new WebSocket(url);

		socket.onopen = () => {
			console.log("✅ WebSocket connected");
			connected.value = true;
		};

		socket.onmessage = (event) => {
			messages.value.push(event.data);
		};

		socket.onclose = () => {
			console.log("🔌 WebSocket disconnected");
			connected.value = false;
		};

		socket.onerror = (err) => {
			console.error("❌ WebSocket error:", err);
		};
	};

	const createRoom = () => {
		const id = generateRoomId();
		room.value = id;
		connectToRoom();
	};

	const generateRoomId = () => {
		return Math.random().toString(36).substring(2, 10); // 8-char random string
	};

	const sendMessage = () => {
		if (socket && message.value.trim() !== "") {
			socket.send(message.value);
			message.value = "";
		}
	};

	onBeforeUnmount(() => {
		if (socket) socket.close();
	});
</script>
