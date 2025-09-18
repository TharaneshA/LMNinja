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
            d="M24 17V42" 
            stroke="currentColor" 
            stroke-linecap="round" 
            stroke-linejoin="round"
        />
        <path 
            :stroke-width="props.strokeWidth" 
            d="M36 6V42" 
            stroke="currentColor" 
            stroke-linecap="round" 
            stroke-linejoin="round"
        />
        <path 
            :stroke-width="props.strokeWidth" 
            d="M12 28V42" 
            stroke="currentColor" 
            stroke-linecap="round" 
            stroke-linejoin="round"
        />
    </svg>
</template>