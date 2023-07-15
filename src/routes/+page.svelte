<script lang="ts">
  import { page } from "$app/stores";
  import { goto } from "$app/navigation";
  import Axios from "axios";

  let dragging = false;
  let uploading = false;

  let input: HTMLInputElement;
  async function uploadHandler() {
    uploading = false;
    let files = input.files;
    if (!files) return;
    let file = files[0];
    let filename = file.name;
    let filesize = file.size;
    let data = { filename, filesize };

    let resp = await Axios.post($page.url.toString(), data);
    let respData = resp.data as { uploadurl: string; code: string };

    // put to uploadurl
    uploading = true;
    await Axios.put(respData.uploadurl, file, {
      onUploadProgress(e) {
        progressValue = (e.loaded / (e.total ?? 1)) * 100;
      },
    });
    uploading = false;

    // redirect url
    let redirectUrl = `/${respData.code}`;
    goto(redirectUrl);
  }

  let progressBar: HTMLDivElement;
  let progressValue = 0;

  $: {
    if (progressBar) {
      progressBar.style.width = `${progressValue}%`;
    }
  }
</script>

<form
  class="flex flex-col h-screen items-center justify-center m0"
  method="POST"
  enctype="multipart/form-data"
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
    <span class="z10 text-xl font-bold text-dark pointer-events-none">
      {uploading ? "Uploading" : "Choose a file"}
    </span>
    <div
      class="absolute top-0 left-0 h-full bg-green rounded"
      role="progressbar"
      aria-valuenow={progressValue}
      bind:this={progressBar}
    />
    <input
      type="file"
      name="file"
      class="absolute top-0 left-0 h-full w-full opacity-0"
      on:change={uploadHandler}
      bind:this={input}
    />
  </label>
  <noscript>
    <button type="submit">Upload</button>
  </noscript>
</form>
