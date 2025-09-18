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
            d="M6 8.25564L24.0086 3L42 8.25564V22.9554C42 33.3633 34.7354 41.5943 24.0026 45.0001C13.2699 41.5943 6 33.3633 6 22.9554V8.25564Z" 
            fill="none" 
            stroke="currentColor" 
            stroke-linejoin="round"
        />
    </svg>
</template>