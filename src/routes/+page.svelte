<script lang="ts">
  import { page } from "$app/stores";
  import { browser } from "$app/environment";
  import { goto } from "$app/navigation";
  import Axios from "axios";

  let dragging = false;

  let form: HTMLFormElement;
  async function uploadHandler() {
    // form.submit();
    // submit using fetch
    let data = new FormData(form);
    let resp = await Axios.post($page.url.toString(), data, {
      onUploadProgress: (e) => {
        updateProgress((e.loaded / (e.total ?? 1)) * 100);
      },
    });

    let json = resp.data;
    let url = json.location as string;
    goto(url);
  }

  let progressBar: HTMLProgressElement;
  function updateProgress(value: number) {
    progressBar.value = value;
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
    class="bg-gray-3 w-2/3 p5 rounded-lg h50 relative flex items-center justify-center {dragging
      ? 'shadow-2xl shadow-blue-2 shadow-inset'
      : ''}"
    on:dragenter={() => (dragging = true)}
    on:dragleave={() => (dragging = false)}
    on:drop={() => (dragging = false)}
  >
    <span class="z10 text-xl font-bold text-dark"> Choose a file </span>
    <progress
      max="100"
      value="50"
      class="wfull hfull absolute top-0 left-0 z0 m0 p0"
      bind:this={progressBar}
    />
    <input
      type="file"
      name="file"
      class="absolute top-0 left-0 h-full w-full opacity-0"
      on:change={uploadHandler}
    />
  </label>
  <noscript>
    <button type="submit">Upload</button>
  </noscript>
</form>

<style>
  progress {
    -webkit-appearance: none;
    appearance: none;
  }
  ::-webkit-progress-bar {
    background-color: transparent;
    border-radius: 0.5rem;
  }
  ::-webkit-progress-value {
    background-color: rgb(59, 130, 246);
    border-radius: 0.5rem;
  }
</style>
