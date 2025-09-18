<script setup>
const props = defineProps({
    strokeWidth: {
        type: [Number, String],
        default: 3,
    },
});
</script>
<template>
    <svg fill="none" viewBox="0 0 48 48" xmlns="http://www.w3.org/2000/svg">
        <path 
            :stroke-width="props.strokeWidth" 
            d="M19 4H37L26 20H41L11 44L16 26H7L19 4Z" 
            stroke="currentColor" 
            stroke-linejoin="round"
        />
    </svg>
</template>