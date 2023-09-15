<script>
  import { DarkMode, Navbar } from "flowbite-svelte";
  import Brand from "../components/brand.svelte";
</script>

<Navbar let:hidden let:toggle rounded color="form">
  <Brand />
  <DarkMode class="text-2xl" />
</Navbar>

<div class="container mx-auto px-4 mt-5">
  <slot />
</div>
