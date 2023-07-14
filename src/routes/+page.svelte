<script lang="ts">
  import { page } from "$app/stores";
  import { browser } from "$app/environment";
  import { goto } from "$app/navigation";

  let dragging = false;

  let fileInput: HTMLInputElement;
  let form: HTMLFormElement;
  async function uploadHandler() {
    // form.submit();
    // submit using fetch
    let data = new FormData(form);
    let resp = await fetch($page.url, {
      method: "POST",
      body: data,
    });
    // get redirect url
    let json = await resp.json();
    let url = json.location as string;
    goto(url);
  }
  $: formData = $page.form;
  $: {
    if (browser && formData && formData.error) {
      console.log(formData.error);
    }
  }
</script>

<form
  class="flex flex-col h-screen items-center justify-center m0"
  method="POST"
  enctype="multipart/form-data"
  bind:this={form}
>
  <h1>File Share</h1>

  <label
    class="bg-gray-2 w-2/3 p5 rounded-lg h50 relative flex items-center justify-center {dragging
      ? 'shadow-2xl shadow-blue-2 shadow-inset'
      : ''}"
    on:dragenter={() => (dragging = true)}
    on:dragleave={() => (dragging = false)}
    on:drop={() => (dragging = false)}
  >
    Choose a file
    <input
      type="file"
      name="file"
      class="absolute top-0 left-0 h-full w-full opacity-0"
      bind:this={fileInput}
      on:change={uploadHandler}
    />
  </label>
  <noscript>
    <button type="submit">Upload</button>
  </noscript>
</form>
