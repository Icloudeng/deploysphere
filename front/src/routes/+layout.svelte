<script>
  import "./styles.css";
  import {
    fns,
    AppShell,
    Navbar,
    Header,
    Title,
    Divider,
  } from "@svelteuidev/core";

  import {
    colorScheme,
    SvelteUIProvider,
    Stack,
    Switch,
    Text,
  } from "@svelteuidev/core";

  function toggleTheme() {
    colorScheme.update((v) => (v === "light" ? "dark" : "light"));
  }
</script>

<SvelteUIProvider withGlobalStyles themeObserver={$colorScheme}>
  <AppShell>
    <!-- <Navbar slot="navbar" hidden={!opened} />
  <Header slot="header" height="" /> -->
    <Stack align="center">
      <Text>Press to change the theme</Text>
      <Switch on:change={toggleTheme} />
    </Stack>

    <slot />
  </AppShell>
</SvelteUIProvider>
