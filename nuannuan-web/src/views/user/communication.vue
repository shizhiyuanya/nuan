<!-- <template>
  <div>
    <h1>Chat</h1>
    <div>
      <input v-model="message" placeholder="Type a message" />
      <input v-model="toUserID" placeholder="Recipient User ID" />
      <button @click="sendMessage">Send</button>
    </div>
    <ul>
      <li v-for="(msg, index) in messages" :key="index">
        <strong>{{ msg.fromUserID }}:</strong> {{ msg.message }}
      </li>
    </ul>
  </div>
</template>


<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';
import { useUserStore } from '@/stores/userStore'; // 导入 Pinia store
import { ElMessage } from 'element-plus';

const message = ref<string>('');
const toUserID = ref<string>('');
const messages = ref<Array<{ fromUserID: string; message: string }>>([]);
const ws = ref<WebSocket | null>(null);

const userStore = useUserStore(); // 使用 Pinia store

onMounted(() => {
  const currentUser = userStore.currentUser;
  if (!currentUser) {
    console.error('用户未登录');
    ElMessage.error('用户未登录!')
    return;
  }
  console.log("我看看用户：", currentUser.id)
  // 建立 WebSocket 连接
  ws.value = new WebSocket(`ws://localhost:8080/ws?userID=${currentUser.id}`);

  ws.value.onmessage = (event) => {
    const data = JSON.parse(event.data);
    messages.value.push({ fromUserID: data.fromUserID, message: data.message });
  };
});

onUnmounted(() => {
  if (ws.value) {
    ws.value.close();
  }
});

const sendMessage = () => {
  if (ws.value && ws.value.readyState === WebSocket.OPEN) {
    ws.value.send(JSON.stringify({ toUserID: toUserID.value, message: message.value }));
    message.value = '';
  } else {
    console.error('WebSocket is not connected');
  }
};
</script>

<style scoped>
/* 样式代码 */
</style> -->




<template>
  <div class="chat-container">
    <!-- 左侧好友列表 -->
    <div class="friend-list">
      <div class="header">
        <el-input
          v-model="searchQuery"
          placeholder="搜索好友"
          clearable
          class="search-input"
        />
        <el-button type="primary" class="add-button" @click="showAddDialog">
          添加好友
        </el-button>
      </div>
      
      <el-scrollbar class="scroll-area">
        <div
          v-for="friend in filteredFriends"
          :key="friend.id"
          class="friend-item"
          :class="{ active: selectedFriend?.id === friend.id }"
          @click="selectFriend(friend)"
        >
          <el-avatar :size="40" :src="friend.avatar" class="avatar" />
          <div class="friend-info">
            <div class="name-line">
              <span class="name">{{ friend.name }}</span>
              <span class="time">{{ formatTime(friend.lastMessageTime) }}</span>
            </div>
            <div class="message-preview">
              <span class="last-message">{{ friend.lastMessage }}</span>
              <el-badge v-if="friend.unread > 0" :value="friend.unread" class="badge" />
            </div>
          </div>
        </div>
      </el-scrollbar>
    </div>

    <!-- 右侧聊天窗口 -->
    <div class="chat-window">
      <div v-if="selectedFriend" class="chat-content">
        <div class="chat-header">
          <el-avatar :size="40" :src="selectedFriend.avatar" />
          <div class="friend-name">{{ selectedFriend.name }}</div>
          <div class="online-status" :class="{ online: selectedFriend.online }">
            {{ selectedFriend.online ? '在线' : '离线' }}
          </div>
        </div>

        <el-scrollbar ref="messageScroll" class="message-area">
          <div
            v-for="(msg, index) in currentMessages"
            :key="index"
            class="message-bubble"
            :class="{ own: msg.fromUserID === userStore.currentUser?.id }"
          >
            <div class="message-content">{{ msg.message }}</div>
            <div class="message-time">{{ formatMessageTime(msg.timestamp) }}</div>
          </div>
        </el-scrollbar>

        <div class="message-input">
          <el-input
            v-model="message"
            type="textarea"
            :rows="3"
            placeholder="输入消息"
            @keyup.enter="sendMessage"
          />
          <div class="send-button">
            <el-button type="primary" @click="sendMessage">发送</el-button>
          </div>
        </div>
      </div>

      <div v-else class="empty-chat">
        请选择好友开始聊天
      </div>
    </div>

    <!-- 添加好友对话框 -->
    <el-dialog v-model="addDialogVisible" title="添加好友" width="30%">
      <el-input v-model="searchId" placeholder="输入用户ID" />
      <template #footer>
        <el-button @click="addDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="addFriend">添加</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { useUserStore } from '@/stores/userStore'
import { ElMessage } from 'element-plus'
import type { Friend } from '@/types/chat'

// 用户状态管理
const userStore = useUserStore()

// WebSocket相关
const ws = ref<WebSocket | null>(null)
const message = ref('')
const selectedFriend = ref<Friend | null>(null)
const messageScroll = ref()

// 好友列表相关
const searchQuery = ref('')
const searchId = ref('')
const addDialogVisible = ref(false)

// 模拟好友数据（实际应从接口获取）
const friends = ref<Friend[]>([
  {
    id: '1002',
    name: '张三',
    avatar: 'https://example.com/avatar1.jpg',
    online: true,
    lastMessage: '你好呀！',
    lastMessageTime: Date.now(),
    unread: 1
  },
  // 可以添加更多测试数据
])

// 所有聊天消息（按好友分组）
const allMessages = ref<{ [key: string]: Array<{
  fromUserID: string
  message: string
  timestamp: number
}> }>({})

// 计算属性
const filteredFriends = computed(() => {
  return friends.value.filter(friend => 
    friend.name.includes(searchQuery.value) || 
    friend.id.includes(searchQuery.value)
  )
})

const currentMessages = computed(() => {
  if (!selectedFriend.value) return []
  return allMessages.value[selectedFriend.value.id] || []
})

// 方法
const selectFriend = (friend: Friend) => {
  selectedFriend.value = friend
  friend.unread = 0
  nextTick(() => {
    scrollToBottom()
  })
}

const sendMessage = () => {
  if (!message.value.trim() || !selectedFriend.value) return

  const msg = {
    toUserID: selectedFriend.value.id,
    message: message.value,
    timestamp: Date.now()
  }

  // 添加到本地消息记录
  addMessageToStore(msg, true)
  
  // 通过WebSocket发送
  if (ws.value && ws.value.readyState === WebSocket.OPEN) {
    ws.value.send(JSON.stringify(msg))
    message.value = ''
    scrollToBottom()
  }
}

const addMessageToStore = (
  msg: {
    fromUserID: string
    message: string
    timestamp: number
  }, 
  isOwn: boolean
) => {
  const friendId = isOwn ? msg.toUserID : msg.fromUserID
  if (!allMessages.value[friendId]) {
    allMessages.value[friendId] = []
  }
  allMessages.value[friendId].push(msg)
}

const scrollToBottom = () => {
  if (messageScroll.value) {
    messageScroll.value.scrollTo({ bottom: 0 })
  }
}

const formatTime = (timestamp: number) => {
  return new Date(timestamp).toLocaleTimeString()
}

const formatMessageTime = (timestamp: number) => {
  return new Date(timestamp).toLocaleString()
}

const showAddDialog = () => {
  addDialogVisible.value = true
}

const addFriend = async () => {
  // 这里应调用添加好友接口
  const newFriend = {
    id: searchId.value,
    name: '新用户',
    avatar: 'https://example.com/avatar-default.jpg',
    online: false,
    lastMessage: '',
    lastMessageTime: Date.now(),
    unread: 0
  }
  friends.value.push(newFriend)
  addDialogVisible.value = false
  searchId.value = ''
}

// 生命周期
onMounted(() => {
  if (!userStore.currentUser) {
    ElMessage.error('请先登录')
    return
  }

  // 初始化WebSocket
  ws.value = new WebSocket(`ws://localhost:8080/ws?userID=${userStore.currentUser.id}`)

  ws.value.onmessage = (event) => {
    const data = JSON.parse(event.data)
    const isCurrent = data.fromUserID === selectedFriend.value?.id
    
    addMessageToStore({
      fromUserID: data.fromUserID,
      message: data.message,
      timestamp: Date.now()
    }, false)

    // 更新好友最后消息状态
    const friend = friends.value.find(f => f.id === data.fromUserID)
    if (friend) {
      friend.lastMessage = data.message
      friend.lastMessageTime = Date.now()
      if (!isCurrent) friend.unread++
    }

    if (isCurrent) scrollToBottom()
  }
})

onUnmounted(() => {
  if (ws.value) ws.value.close()
})
</script>

<style scoped>
.chat-container {
  display: flex;
  height: 100vh;
  background: #f0f2f5;
}

.friend-list {
  width: 300px;
  background: white;
  border-right: 1px solid #e6e6e6;
}

.header {
  padding: 15px;
  border-bottom: 1px solid #e6e6e6;
}

.search-input {
  margin-bottom: 10px;
}

.add-button {
  width: 100%;
}

.scroll-area {
  height: calc(100vh - 120px);
}

.friend-item {
  display: flex;
  padding: 12px;
  cursor: pointer;
  transition: background 0.3s;
}

.friend-item:hover {
  background: #f5f5f5;
}

.friend-item.active {
  background: #e6f4ff;
}

.avatar {
  margin-right: 12px;
}

.friend-info {
  flex: 1;
  min-width: 0;
}

.name-line {
  display: flex;
  justify-content: space-between;
  margin-bottom: 4px;
}

.name {
  font-weight: 500;
  overflow: hidden;
  text-overflow: ellipsis;
}

.time {
  font-size: 12px;
  color: #999;
}

.message-preview {
  display: flex;
  justify-content: space-between;
  font-size: 13px;
  color: #666;
}

.last-message {
  flex: 1;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}

.badge {
  margin-left: 8px;
}

.chat-window {
  flex: 1;
  background: white;
}

.empty-chat {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
  color: #999;
}

.chat-content {
  height: 100vh;
  display: flex;
  flex-direction: column;
}

.chat-header {
  display: flex;
  align-items: center;
  padding: 15px;
  border-bottom: 1px solid #e6e6e6;
}

.friend-name {
  margin-left: 12px;
  font-weight: 500;
}

.online-status {
  margin-left: 8px;
  font-size: 12px;
  color: #666;
}

.online-status.online {
  color: #67c23a;
}

.message-area {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
}

.message-bubble {
  margin: 10px 0;
  max-width: 70%;
  padding: 10px 15px;
  border-radius: 8px;
  background: #f0f0f0;
}

.message-bubble.own {
  margin-left: auto;
  background: #95ec69;
}

.message-content {
  word-break: break-word;
}

.message-time {
  font-size: 12px;
  color: #666;
  margin-top: 4px;
  text-align: right;
}

.message-input {
  border-top: 1px solid #e6e6e6;
  padding: 15px;
}

.send-button {
  margin-top: 10px;
  text-align: right;
}
</style>
