<script setup lang="ts">
import { ref } from "vue";
import type { Link } from "@/entities/link/link.types";
import CreatedLinkResult from "@/features/create-link/CreatedLinkResult.vue";
import CreateLinkForm from "@/features/create-link/CreateLinkForm.vue";
import { ROUTES } from "@/shared/lib/routes/paths";
import { UiPageHeader } from "@/shared/ui";

const createdLink = ref<Link | null>(null);
const formKey = ref(0);

const onLinkCreated = (link: Link) => {
  createdLink.value = link;
};

const onCreateMore = () => {
  createdLink.value = null;
  formKey.value += 1;
};
</script>

<template>
  <section class="create-link-page">
    <UiPageHeader
      title="Создание короткой ссылки"
      subtitle="Укажите target URL и optional alias, чтобы получить short URL."
      :back-to="ROUTES.dashboard"
      back-label="Dashboard"
    />

    <div class="create-link-page__content">
      <CreateLinkForm v-if="!createdLink" :key="formKey" @created="onLinkCreated" />
      <CreatedLinkResult v-else :link="createdLink" @create-more="onCreateMore" />
    </div>
  </section>
</template>

<style scoped>
.create-link-page {
  width: 100%;
}

.create-link-page__content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}
</style>
