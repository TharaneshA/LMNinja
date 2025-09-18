<script setup>
import { ref, watch, nextTick } from 'vue';
import { useThemeVars } from 'naive-ui';
import useChatStore from 'stores/chat.js';

const chatStore = useChatStore();
const prompt = ref('');
const scrollArea = ref(null);
const themeVars = useThemeVars();

const scrollToBottom = () => {
  nextTick(() => {
    if (scrollArea.value) {
      scrollArea.value.scrollTop = scrollArea.value.scrollHeight;
    }
  });
};

watch(() => chatStore.messages, () => {
    scrollToBottom();
}, { deep: true });

const handleSendMessage = () => {
    if (!prompt.value.trim()) return;
    chatStore.sendMessage(prompt.value);
    prompt.value = '';
};
</script>

<template>
    <div class="chat-container">
        <div class="message-history" ref="scrollArea">
            <div v-for="(msg, index) in chatStore.messages" :key="index" :class="['message', `message-${msg.sender}`]">
                <div class="message-bubble" :class="{ 'error-bubble': msg.isError }">
                    <p>{{ msg.text }}</p>
                </div>
            </div>
            <div v-if="chatStore.isStreamingResponse" class="message message-model">
                <div class="message-bubble">
                    <n-spin size="small" />
                </div>
            </div>
        </div>
        <div class="input-area">
            <n-input
                v-model:value="prompt"
                type="textarea"
                :autosize="{ minRows: 1, maxRows: 5 }"
                placeholder="Send a message... (Enter to send, Shift+Enter for newline)"
                @keypress.enter.prevent.exact="handleSendMessage"
                :disabled="chatStore.isStreamingResponse"
            />
            <n-button @click="handleSendMessage" type="primary" :loading="chatStore.isStreamingResponse">Send</n-button>
        </div>
    </div>
</template>

<style lang="scss" scoped>
.chat-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  background-color: v-bind('themeVars.bodyColor');
}
.message-history {
  flex-grow: 1;
  overflow-y: auto;
  padding: 16px;
}
.message {
  display: flex;
  margin-bottom: 12px;
}
.message-user { justify-content: flex-end; }
.message-model { justify-content: flex-start; }
.message-bubble {
  max-width: 85%;
  padding: 10px 15px;
  border-radius: 18px;
  line-height: 1.5;
  p { margin: 0; white-space: pre-wrap; word-wrap: break-word; }
}
.message-user .message-bubble {
  background-color: v-bind('themeVars.primaryColor');
  color: white;
}
.message-model .message-bubble {
  background-color: v-bind('themeVars.cardColor');
}
.error-bubble {
  background-color: v-bind('themeVars.errorColor');
  color: white;
}
.input-area {
  display: flex;
  gap: 10px;
  align-items: flex-end;
  padding: 8px;
  border-top: 1px solid v-bind('themeVars.borderColor');
}
</style>