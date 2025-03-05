import Register from "@/views/user/register.vue";
import Login from "@/views/user/login.vue";
import communication from "@/views/user/communication.vue";
export default [
    {
        path: 'register',
        name: 'register',
        component: Register
    },
    {
        path: 'login',
        name: 'login',
        component: Login
    },
    {
        path: 'communication',
        name: 'communication',
        component: communication
    }
]