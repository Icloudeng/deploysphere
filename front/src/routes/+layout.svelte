<script>
  import "./styles.css";
  import {
    colorScheme,
    SvelteUIProvider,
    AppShell,
    Navbar,
  } from "@svelteuidev/core";
  import NavbarContent from "./components/NavbarContent.svelte";
</script>

<SvelteUIProvider withGlobalStyles themeObserver={$colorScheme}>
  <AppShell class="app-shell">
    <Navbar slot="navbar" pt="md" height="auto" fixed>
      <div class="container">
        <NavbarContent />
      </div>
    </Navbar>

    <!-- -->
    <div class="content container">
      <slot />
    </div>
  </AppShell>
</SvelteUIProvider>

<style lang="css">
  .content {
    margin-top: 5rem;
  }

  :root {
    --container-padding: 20px;
    --container-width: calc(100% - var(--container-padding) * 1.3);
  }
  .container {
    box-sizing: border-box;
    width: 100%;
    padding-left: calc(50% - var(--container-width) / 1.3);
    padding-right: calc(50% - var(--container-width) / 1.3);
  }

  @media only screen and (min-width: 1040px) {
    :root {
      --container-width: 1000px;
    }
  }
</style>
