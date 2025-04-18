<template>
	<div class="flex flex-col items-center justify-center min-h-screen p-6">
		<div
			class="w-full max-w-lg bg-white rounded-2xl shadow-md p-6 flex flex-col h-[80vh]"
		>
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

	function joinRoom() {
		const url = `ws://localhost:8080/ws?room=${room.value}&username=${username.value}`;
		socket = new WebSocket(url);

		socket.onopen = () => {
			console.log("Connected to server");
			connected.value = true;
		};

		socket.onmessage = (event) => {
			messages.value.push(event.data);
		};

		socket.onclose = () => {
			console.log("Disconnected");
			connected.value = false;
		};

		socket.onerror = (error) => {
			console.error("WebSocket error:", error);
		};
	}

	function sendMessage() {
		if (socket && message.value.trim() !== "") {
			socket.send(message.value);
			message.value = "";
		}
	}

	onBeforeUnmount(() => {
		if (socket) socket.close();
	});
</script>
