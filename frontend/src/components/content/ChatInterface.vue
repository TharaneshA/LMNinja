<script setup>
import { ref, watch, nextTick } from 'vue';
import { useThemeVars } from 'naive-ui';
import useChatStore from 'stores/chat.js';

const chatStore = useChatStore();
const prompt = ref('');
const scrollArea = ref(null); // Ref for the message history div
const themeVars = useThemeVars();

// This function ensures the view scrolls to the latest message
const scrollToBottom = () => {
  nextTick(() => {
    if (scrollArea.value) {
      scrollArea.value.scrollTop = scrollArea.value.scrollHeight;
    }
  });
};

// Watch for new messages and auto-scroll
watch(() => chatStore.messages, () => {
    scrollToBottom();
}, { deep: true });

const handleSendMessage = () => {
    const trimmedPrompt = prompt.value.trim();
    if (!trimmedPrompt) return; // Don't send empty messages
    
    // This calls the action in our Pinia store, which handles the backend call
    chatStore.sendMessage(trimmedPrompt);

    prompt.value = ''; // Clear the input field immediately
};
</script>

<template>
    <div class="chat-container">
        <!-- Message History Area -->
        <div class="message-history" ref="scrollArea">
            <!-- Welcome Message -->
            <div v-if="chatStore.messages.length === 0 && !chatStore.isStreamingResponse" class="welcome-message">
                <n-empty :description="`Ready to chat with ${chatStore.activeModel?.name}.`" />
            </div>

            <!-- Conversation -->
            <div v-for="(msg, index) in chatStore.messages" :key="index" :class="['message', `message-${msg.sender}`]">
                <div class="message-bubble" :class="{ 'error-bubble': msg.isError }">
                    <p>{{ msg.text }}</p>
                </div>
            </div>
            
            <!-- Loading Indicator for Model Response -->
            <div v-if="chatStore.isStreamingResponse" class="message message-model">
                <div class="message-bubble">
                    <n-spin size="small" />
                </div>
            </div>
        </div>

        <!-- Input Area -->
        <div class="input-area">
            <n-input
                v-model:value="prompt"
                type="textarea"
                :autosize="{ minRows: 1, maxRows: 5 }"
                placeholder="Send a message... (Enter to send, Shift+Enter for newline)"
                @keypress.enter.prevent.exact="handleSendMessage"
                :disabled="chatStore.isStreamingResponse || !chatStore.isModelLoaded"
            />
            <n-button 
                @click="handleSendMessage" 
                type="primary" 
                :loading="chatStore.isStreamingResponse"
                :disabled="!prompt.trim() || !chatStore.isModelLoaded"
            >
                Send
            </n-button>
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
  max-width: 90%;
}
.message-user { 
  margin-left: auto; /* Aligns user messages to the right */
  justify-content: flex-end;
}
.message-model { 
  margin-right: auto; /* Aligns model messages to the left */
  justify-content: flex-start;
}
.message-bubble {
  padding: 10px 15px;
  border-radius: 18px;
  line-height: 1.5;
  p { 
    margin: 0; 
    white-space: pre-wrap; 
    word-wrap: break-word;
  }
}
.message-user .message-bubble {
  background-color: v-bind('themeVars.primaryColor');
  color: white;
  border-bottom-right-radius: 4px;
}
.message-model .message-bubble {
  background-color: v-bind('themeVars.cardColor');
  border-bottom-left-radius: 4px;
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
  flex-shrink: 0;
}
.welcome-message {
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
}
</style>