<script>
  import { onMount } from "svelte";
  import "../app.postcss";
  import "./styles.css";
  import { Spinner } from "flowbite-svelte";
  import { api_root_query } from "$lib/authenticate";

  let loading = true;

  onMount(() => {
    api_root_query().then((res) => {
      loading = false;
    });
  });
</script>

{#if loading}
  <div class="fixed inset-0 flex justify-center items-center">
    <Spinner size="12" />
  </div>
{:else}
  <slot />
{/if}
