<script>
  import { DarkMode } from "flowbite-svelte";
</script>

<div class="container mx-auto px-4 mt-3">
  <DarkMode class="text-2xl" />

  <div class="mt-4 w-full">
    <slot />
  </div>
</div>
