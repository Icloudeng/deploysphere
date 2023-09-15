<script>
  import { Label, Input, Button } from "flowbite-svelte";
</script>

<svelte:head>
  <title>Login</title>
</svelte:head>

<form
  class="w-full md:w-1/2 flex justify-center flex-col mx-auto"
  autocomplete="off"
>
  <div class="mb-6">
    <Label for="email" class="block mb-2">Email</Label>
    <Input id="email" type="email" placeholder="Email input" />
  </div>

  <div class="mb-6">
    <Label for="password" class="block mb-2">Password</Label>
    <Input id="password" type="password" placeholder="Password input" />
  </div>

  <div>
    <Button>Login</Button>
  </div>
</form>
