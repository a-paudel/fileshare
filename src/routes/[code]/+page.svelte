<script lang="ts">
  import { page } from "$app/stores";
  import type { PageData } from "./$types";

  export let data: PageData;

  function copyHandler(
    event: MouseEvent & { currentTarget: EventTarget & HTMLElement }
  ) {
    navigator.clipboard.writeText($page.url.toString());
  }
</script>

<a
  href="/"
  class="decoration-none bg-blue-5 hover:bg-blue-6 text-white p4 rounded-md fixed top-0 left-0 m2"
>
  Upload another file
</a>
<div class="flex flex-col h-screen items-center justify-center m0">
  <h1>File Share</h1>
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
  <code
    class="bg-gray-2 p2 rounded shadow shadow-inset select-all"
    on:click={copyHandler}
  >
    {$page.url}
  </code>
  <h2>{data.file.name}</h2>
  <a
    href={data.file.url}
    class="decoration-none bg-blue-5 hover:bg-blue-6 text-white p4 rounded-md"
  >
    Download File
  </a>
</div>
