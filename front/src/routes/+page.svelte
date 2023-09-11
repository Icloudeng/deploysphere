<script>
  import Navbar from "./components/navbar.svelte";
  import { Tabs, TabItem } from "flowbite-svelte";
  import { CompressOutline, SortHorizontalSolid } from "flowbite-svelte-icons";
</script>

<svelte:head>
  <title>Platfomr Installer</title>
  <meta name="description" content="Svelte demo app" />
</svelte:head>

<Navbar />

<Tabs style="underline" class="mt-5">
  <TabItem>
    <div slot="title" class="flex items-center gap-2">
      <CompressOutline size="sm" />
      Resources
    </div>
  </TabItem>

  <TabItem>
    <div slot="title" class="flex items-center gap-2">
      <SortHorizontalSolid size="sm" />
      Jobs
    </div>
  </TabItem>
</Tabs>
