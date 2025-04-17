import { createMemoryHistory, createRouter } from "vue-router";

import HomeView from "./views/HomeView.vue";
import RoomView from "./views/RoomView.vue";
import LandingView from "./views/LandingView.vue";

const routes = [
	{ path: "/", component: LandingView },
	{ path: "/home", component: HomeView },
	{ path: "/room", component: RoomView },
];

const router = createRouter({
	history: createMemoryHistory(),
	routes,
});

export default router;
