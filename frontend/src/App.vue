<script setup lang="ts">
import { RouterView } from "vue-router";
</script>

<template>
  <main class="app-shell">
    <header class="app-header">
      <strong>TrackLink</strong>
      <span class="app-subtitle">Frontend foundation</span>
    </header>

    <section class="page-card">
      <RouterView />
    </section>
  </main>
</template>
