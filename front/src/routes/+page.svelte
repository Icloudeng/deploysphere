<script>
  // @ts-nocheck
  import { Seo, Tabs } from "@svelteuidev/core";

  import { Stack, ColorWheel, PieChart } from "radix-icons-svelte";
  import Resources from "./components/home/Resources.svelte";
  import Jobs from "./components/home/Jobs.svelte";
  import History from "./components/home/History.svelte";
</script>

<Seo title="Resources" titleTemplate="%t% | Platform Installer" />

<Tabs>
  <Tabs.Tab label="Resources" icon={Stack}>
    <Resources />
  </Tabs.Tab>

  <Tabs.Tab label="Jobs" icon={PieChart}>
    <Jobs />
  </Tabs.Tab>

  <Tabs.Tab label="History" icon={ColorWheel}>
    <History />
  </Tabs.Tab>
</Tabs>
