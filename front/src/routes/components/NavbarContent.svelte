<script>
  //   @ts-nocheck

  import { base } from "$app/paths";
  import {
    colorScheme,
    Image,
    Flex,
    Switch,
    Text,
    Button,
    Menu,
    Space,
  } from "@svelteuidev/core";
  import {
    Gear,
    ComponentPlaceholder,
    Stack,
    Plus,
    Commit,
    // Group,
  } from "radix-icons-svelte";

  function toggleTheme() {
    colorScheme.update((v) => (v === "light" ? "dark" : "light"));
  }
</script>

<Flex justify="space-between">
  <Flex align="center">
    <Image width={35} height={35} src="{base}/favicon.png" />
    <Text>PInstall</Text>
  </Flex>

  <Flex>
    <Switch on:change={toggleTheme} />
    <Space w="10" />
    <Menu>
      <Button slot="control">
        <Plus />
        <span>New Pipline</span>
      </Button>

      <!-- <Menu.Item icon={Group}>Workspace</Menu.Item> -->
      <Menu.Item icon={Stack}>Resources</Menu.Item>
      <Menu.Item icon={ComponentPlaceholder}>Provisioning</Menu.Item>
      <Menu.Item icon={Gear}>Configuration</Menu.Item>
      <Menu.Item icon={Commit}>Proxy Host</Menu.Item>
    </Menu>
  </Flex>
</Flex>
