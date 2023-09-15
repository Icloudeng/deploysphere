<script>
  import {
    Button,
    Dropdown,
    DropdownItem,
    Navbar,
    NavBrand,
    DarkMode,
  } from "flowbite-svelte";
  import {
    PlusOutline,
    CompressOutline,
    CaretDownOutline,
    AdjustmentsHorizontalOutline,
    CodePullRequestOutline,
  } from "flowbite-svelte-icons";
  import Brand from "./brand.svelte";

  let open = false;

  const menus = [
    {
      name: "Resources",
      icon: CompressOutline,
    },
    {
      name: "Provisioning",
      icon: CaretDownOutline,
    },
    {
      name: "Configuration",
      icon: AdjustmentsHorizontalOutline,
    },
    {
      name: "Proxy Host",
      icon: CodePullRequestOutline,
    },
  ];
</script>

<Navbar let:hidden let:toggle rounded color="form">
  <Brand />

  <div class="flex space-x-3">
    <div>
      <Button>
        New Pipline
        <PlusOutline class="w-3 h-3 ml-2 text-white dark:text-white" />
      </Button>

      <Dropdown {open}>
        {#each menus as menu, index (index)}
          <DropdownItem
            on:click={() => (open = false)}
            class="flex items-center space-x-3"
          >
            <svelte:component this={menu.icon} class="w-3 h-3 ml-2" />
            <span>{menu.name}</span>
          </DropdownItem>
        {/each}
      </Dropdown>
    </div>

    <DarkMode class="text-2xl" />
  </div>
</Navbar>
