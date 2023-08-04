<script lang="ts">
  import { goto } from "$app/navigation";
  import { page } from "$app/stores";
  import { PUBLIC_API_URL } from "$env/static/public";
  import { onMount } from "svelte";
  import { sizeFormatter } from "human-readable";

  type FileType = {
    CreatedAt: string;
    UpdatedAt: string;
    Code: string;
    Filename: string;
    Filesize: number;
  };

  let file: FileType | undefined;

  onMount(async () => {
    // get file from api if not exists redirect to home
    let code = $page.params.code;
    let resp = await fetch(`${PUBLIC_API_URL}/files/${code}`);
    if (!resp.ok) {
      goto("/");
    }
    file = await resp.json();
  });

  function linkClickHandler() {
    // goto home page
    setInterval(() => {
      goto("/");
    }, 500);
  }

  function copyUrl() {
    navigator.clipboard.writeText($page.url.toString());
  }

  onMount(() => {
    // copy url to clipboard
    copyUrl();
  });
</script>

<a href="/" class="button bg-gray-4 fixed top-0 left-0 m4"
  >Upload another file</a
>
<!-- download page -->
<div class="flex flex-col h-screen items-center justify-center">
  {#if file}
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
    <pre
      class="bg-gray-4 text-dark rounded select-all"
      on:click={copyUrl}>{$page.url}</pre>
    <h2 class="m0">{file.Filename}</h2>
    <a
      href="{PUBLIC_API_URL}/files/{file?.Code}/download"
      class="button bg-gray-7 text-light"
      on:click={linkClickHandler}
    >
      <span class="text-3xl"> Download </span>
      <p class="m0 text-lg">{sizeFormatter()(file.Filesize)}</p>
    </a>
  {:else}
    <i class="fa fa-circle-notch animate-spin" />
  {/if}
</div>
