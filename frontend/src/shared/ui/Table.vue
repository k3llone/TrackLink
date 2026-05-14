<script setup lang="ts">
import UiPageState from "./PageState.vue";

export type UiTableColumn = {
  key: string;
  label: string;
  width?: string;
  align?: "left" | "center" | "right";
};

const props = withDefaults(
  defineProps<{
    columns?: UiTableColumn[];
    rows?: unknown[];
    loading?: boolean;
    emptyText?: string;
    rowClickable?: boolean;
  }>(),
  {
    columns: () => [],
    rows: () => [],
    loading: false,
    emptyText: "No data yet",
    rowClickable: false,
  },
);

const emit = defineEmits<{
  "row-click": [row: unknown, rowIndex: number];
}>();

const onRowClick = (row: unknown, rowIndex: number) => {
  if (props.rowClickable) {
    emit("row-click", row, rowIndex);
  }
};
</script>

<template>
  <div class="ui-table">
    <slot v-if="loading" name="loading">
      <UiPageState type="loading" title="Loading data" />
    </slot>

    <slot v-else-if="!rows.length" name="empty">
      <UiPageState type="empty" :description="emptyText" />
    </slot>

    <template v-else>
      <slot name="header" :columns="columns">
        <table class="ui-table__native">
          <thead>
            <tr>
              <th
                v-for="column in columns"
                :key="column.key"
                :style="{ width: column.width }"
                :class="`align-${column.align ?? 'left'}`"
              >
                {{ column.label }}
              </th>
              <th v-if="$slots.actions" class="align-right">Actions</th>
            </tr>
          </thead>

          <tbody>
            <slot name="body" :rows="rows" :columns="columns">
              <tr
                v-for="(row, rowIndex) in rows"
                :key="rowIndex"
                :class="{ 'is-clickable': rowClickable }"
                :tabindex="rowClickable ? 0 : undefined"
                :role="rowClickable ? 'button' : undefined"
                @click="onRowClick(row, rowIndex)"
                @keydown.enter.prevent="onRowClick(row, rowIndex)"
                @keydown.space.prevent="onRowClick(row, rowIndex)"
              >
                <slot name="row" :row="row" :row-index="rowIndex">
                  <td
                    v-for="column in columns"
                    :key="column.key"
                    :class="`align-${column.align ?? 'left'}`"
                  >
                    <slot name="cell" :row="row" :column="column">
                      {{ (row as Record<string, unknown>)[column.key] ?? "—" }}
                    </slot>
                  </td>
                </slot>
                <td v-if="$slots.actions" class="align-right">
                  <slot name="actions" :row="row" :row-index="rowIndex" />
                </td>
              </tr>
            </slot>
          </tbody>
        </table>
      </slot>
    </template>
  </div>
</template>

<style scoped>
.ui-table {
  width: 100%;
  overflow-x: auto;
}

.ui-table__native {
  width: 100%;
  min-width: 680px;
  border-collapse: separate;
  border-spacing: 0 8px;
}

.ui-table__native thead th {
  padding: 8px 12px;
  color: var(--tl-color-text-muted);
  font-size: 12px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.02em;
}

.ui-table__native tbody td {
  background: var(--tl-color-surface-muted);
  padding: 14px 12px;
  font-size: 14px;
  color: var(--tl-color-text);
  transition: background-color 0.2s ease;
}

.ui-table__native tbody tr.is-clickable {
  cursor: pointer;
}

.ui-table__native tbody tr.is-clickable:hover td,
.ui-table__native tbody tr.is-clickable:focus-visible td {
  background: #e8e2f4;
}

.ui-table__native tbody tr.is-clickable:focus-visible {
  outline: 2px solid rgb(109 74 255 / 35%);
  outline-offset: 2px;
}

.ui-table__native tbody td:first-child {
  border-top-left-radius: var(--tl-radius-md);
  border-bottom-left-radius: var(--tl-radius-md);
}

.ui-table__native tbody td:last-child {
  border-top-right-radius: var(--tl-radius-md);
  border-bottom-right-radius: var(--tl-radius-md);
}

.align-left {
  text-align: left;
}

.align-center {
  text-align: center;
}

.align-right {
  text-align: right;
}

@media (max-width: 767px) {
  .ui-table__native {
    min-width: 560px;
  }
}
</style>
