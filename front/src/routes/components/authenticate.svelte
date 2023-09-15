<script lang="ts">
  let loading = false;

  export const redirect_type: "unauthenticated" | "authenticated" =
    "authenticated";
</script>

<slot />
