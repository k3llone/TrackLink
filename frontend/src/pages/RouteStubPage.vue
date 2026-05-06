<script setup lang="ts">
defineProps<{
  title: string;
}>();
</script>

<template>
  <section>
    <h1>{{ title }}</h1>
    <p class="text-muted">Страница еще не реализована в рамках текущей карточки.</p>
  </section>
</template>
