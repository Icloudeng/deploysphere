<script>
  import "../app.postcss";
  import "./styles.css";

  import {
    CompressOutline,
    SortHorizontalSolid,
    HourglassOutline,
  } from "flowbite-svelte-icons";

  import { page } from "$app/stores";
  import Navbar from "./components/navbar.svelte";
  import cn from "$lib/utils/cn";

  $: activeUrl = $page.url.pathname;

  $: menus = [
    {
      name: "Resources",
      icon: CompressOutline,
      href: "/",
      active: activeUrl == "/",
    },
    {
      name: "Jobs",
      icon: SortHorizontalSolid,
      href: "/jobs",
      active: activeUrl == "/jobs",
    },
    {
      name: "History",
      icon: HourglassOutline,
      href: "/history",
      active: activeUrl == "/history",
    },
  ];
</script>

<div class="container mx-auto px-4 mt-3">
  <Navbar />

  <div class="border-b border-gray-200 dark:border-gray-700">
    <ul
      class="flex flex-wrap -mb-px text-sm font-medium text-center text-gray-500 dark:text-gray-400"
    >
      {#each menus as menu, index (index)}
        <li class="mr-2">
          <a
            href={menu.href}
            class={cn(
              "inline-flex items-center justify-center p-4",
              menu.active && "border-b-2 border-b-primary-500 rounded-t-lg"
            )}
            aria-current="page"
          >
            <svelte:component this={menu.icon} class="mr-3 w-3 h-3" />
            {menu.name}
          </a>
        </li>
      {/each}
    </ul>
  </div>

  <slot />
</div>
