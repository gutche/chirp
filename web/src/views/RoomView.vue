<script setup>
	import { ref, onBeforeUnmount, onMounted } from "vue";
	import { useRoute, useRouter } from "vue-router";

	const route = useRoute();
	const router = useRouter();
	const roomID = ref(route.query.roomID);
	const endsAt = ref(0);
	const expired = ref(true);
	const message = ref("");
	const messages = ref([]);
	const connected = ref(false);
	const remainingTime = ref("00:00:00");
	const redirectCountdown = ref(3);

	let socket = null;
	let countdownInterval = null;

	function formatTime(seconds) {
		const h = Math.floor(seconds / 3600)
			.toString()
			.padStart(2, "0");
		const m = Math.floor((seconds % 3600) / 60)
			.toString()
			.padStart(2, "0");
		const s = Math.floor(seconds % 60)
			.toString()
			.padStart(2, "0");
		return `${h}:${m}:${s}`;
	}

	function startCountdown(endTimestamp) {
		countdownInterval = setInterval(() => {
			const now = Math.floor(Date.now() / 1000);
			const remaining = endTimestamp - now;

			if (remaining <= 0) {
				remainingTime.value = "00:00:00";
				clearInterval(countdownInterval);
				socket?.close();
				expired.value = true;

				const redirectTimer = setInterval(() => {
					redirectCountdown.value--;

					if (redirectCountdown.value <= 0) {
						clearInterval(redirectTimer);
						router.push("/home");
					}
				}, 1000);

				return;
			}

			remainingTime.value = formatTime(remaining);
		}, 1000);
	}

	onMounted(() => {
		const { expiration, username, create } = route.query;

		const params = new URLSearchParams({
			roomID: roomID.value,
			username: username,
		});

		if (create && expiration) {
			params.append("expiration", expiration);
			params.append("create", create);
		}

		const url = `ws://localhost:8080/ws?${params.toString()}`;
		socket = new WebSocket(url);

		socket.onopen = () => {
			connected.value = true;
		};

		socket.onmessage = (event) => {
			let data;
			try {
				data = JSON.parse(event.data);
				if (data.type === "init" && data.endsAt) {
					endsAt.value = data.endsAt;
					startCountdown(data.endsAt);
					return;
				}
			} catch (err) {
				messages.value.push(event.data);
				return;
			}

			if (data.message) {
				messages.value.push(data.message);
			}
		};

		socket.onclose = (event) => {
			connected.value = false;

			if (!expired.value) {
				expired.value = true;
				redirectCountdown.value = 3;

				const redirectTimer = setInterval(() => {
					redirectCountdown.value--;
					if (redirectCountdown.value <= 0) {
						clearInterval(redirectTimer);
						router.push("/home");
					}
				}, 1000);
			}
		};

		socket.onerror = (err) => {
			console.error("WebSocket error:", err);
		};
	});

	const sendMessage = () => {
		if (socket && message.value.trim() !== "") {
			socket.send(message.value);
			message.value = "";
		}
	};

	const leaveRoom = () => {
		if (socket) socket.close();
		if (countdownInterval) clearInterval(countdownInterval);
		router.push("/home");
	};

	onBeforeUnmount(() => {
		if (socket) socket.close();
		if (countdownInterval) clearInterval(countdownInterval);
	});
</script>

<template>
	<div
		class="flex flex-col items-center justify-center min-h-screen bg-gray-100 p-4"
	>
		<div
			class="w-full max-w-2xl bg-white rounded-xl shadow p-6 flex flex-col h-[80vh]"
		>
			<div
				v-if="connected"
				class="h-full flex flex-col"
			>
				<div class="flex justify-between items-center mb-4">
					<h2 class="text-xl font-semibold text-gray-800">
						Room: {{ roomID }}
					</h2>
					<div class="text-sm text-gray-600">
						⏳ {{ remainingTime }}
					</div>
					<button
						@click="leaveRoom"
						class="ml-2 px-2 py-1 border border-gray-300 rounded text-xs text-gray-600 hover:bg-red-400"
					>
						Leave Room
					</button>
				</div>

				<div
					class="flex-1 overflow-y-auto border border-gray-400 rounded p-4 mb-4 bg-gray-50 space-y-2"
				>
					<div
						v-for="(msg, index) in messages"
						:key="index"
						class="text-sm text-gray-800 bg-blue-100 px-3 py-2 rounded-md w-fit max-w-[80%]"
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
			<div
				v-else-if="expired && !connected"
				class="flex flex-col items-center justify-center h-full"
			>
				<p class="text-red-400 text-lg">Room has expired.</p>
				<p class="text-sm mt-1 text-gray-500">
					Redirecting in {{ redirectCountdown }}
				</p>
			</div>
			<div
				v-else
				class="flex items-center justify-center h-full text-gray-600 text-lg"
			>
				Connecting...
			</div>
		</div>
	</div>
</template>
