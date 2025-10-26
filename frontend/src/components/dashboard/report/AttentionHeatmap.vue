<script setup>
import { computed } from 'vue';
import { useThemeVars } from 'naive-ui';

const props = defineProps({
  prompt: String,
  explainability: Object,
});

const themeVars = useThemeVars();

const getHighlightColor = (score) => {

    const alpha = Math.min(1, Math.max(0.1, props.explainability.distraction_score || 0));
    return `rgba(232, 128, 128, ${alpha})`;
};

const highlightedHtml = computed(() => {
  if (!props.prompt || !props.explainability?.hotspot_tokens) {
    return props.prompt.replace(/</g, "&lt;").replace(/>/g, "&gt;");
  }

  let html = props.prompt.replace(/</g, "&lt;").replace(/>/g, "&gt;");
  const color = getHighlightColor();
  const sortedHotspots = [...props.explainability.hotspot_tokens].sort((a, b) => b.length - a.length);

  sortedHotspots.forEach(token => {
    const sanitizedToken = token.replace(/</g, "&lt;").replace(/>/g, "&gt;");
    const regex = new RegExp(`(${escapeRegExp(sanitizedToken)})`, 'gi');
    html = html.replace(regex, `<span class="hotspot" style="background-color: ${color};">$1</span>`);
  });
  
  return html;
});

function escapeRegExp(string) {
  return string.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
}
</script>

<template>
  <div class="heatmap-container">
    <p v-html="highlightedHtml"></p>
  </div>
</template>

<style scoped>
.heatmap-container {
  background-color: v-bind('themeVars.codeColor');
  padding: 16px;
  border-radius: 6px;
  white-space: pre-wrap;
  word-break: break-all;
  line-height: 1.6;
  color: v-bind('themeVars.textColor2');
  user-select: text;
}
:deep(.hotspot) {
  padding: 1px 3px;
  border-radius: 3px;
  color: #FFFFFF;
}
</style>