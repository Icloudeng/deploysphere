<script>
  // @ts-nocheck

  import { Input, Space } from "@svelteuidev/core";
  import { MagnifyingGlass } from "radix-icons-svelte";
</script>

<Space h="lg" />

<Input
  icon={MagnifyingGlass}
  placeholder="Search"
  rightSectionWidth={70}
  styles={{ rightSection: { pointerEvents: "none" } }}
/>

<Space h="lg" />

<div>Resources</div>
