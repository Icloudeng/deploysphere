<script>
  import "./styles.css";
  import { base } from "$app/paths";
  import {
    colorScheme,
    SvelteUIProvider,
    AppShell,
    Navbar,
    Image,
    Flex,
    Container,
  } from "@svelteuidev/core";

  function toggleTheme() {
    colorScheme.update((v) => (v === "light" ? "dark" : "light"));
  }
</script>

<SvelteUIProvider withGlobalStyles themeObserver={"dark"}>
  <AppShell>
    <Navbar fixed className="">
      <!-- <Container> -->
      <Flex justify="space-between">
        <Image width={50} height={50} src="{base}/favicon.png" />
        <span />
      </Flex>
      <!-- </Container> -->
    </Navbar>

    <!-- <Stack align="center">
      <Text>Press to change the theme</Text>
      <Switch on:change={toggleTheme} />
    </Stack> -->

    <Container>
      <slot />
    </Container>
  </AppShell>
</SvelteUIProvider>
