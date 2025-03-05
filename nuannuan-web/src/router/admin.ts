import AdminView from "@/views/admin/adminPicture.vue";
import VideoAdmin from "@/views/admin/videoAdmin.vue";
export default [
    {
        path: 'picture',
        name: 'adminPicture',
        component: AdminView
    },
    {
        path: 'video',
        // 这个video和主页的video重名了 就用不了了
        name: 'adminVideo',
        component: VideoAdmin
    }
]