<script lang="ts">
  import { goto } from "$app/navigation";
  import { PUBLIC_API_URL } from "$env/static/public";
  import Axios from "axios";

  let fileInput: HTMLInputElement;
  let progress = 0;
  let progressBar: HTMLDivElement;
  let label = "Upload File";
  $: {
    if (progressBar) {
      progressBar.style.width = `${progress}%`;
    }
  }
  $: {
    if (progress === 0) {
      label = "Upload File";
    } else if (progress === 100) {
      label = "Uploaded";
    } else {
      label = `Uploading file: ${progress}%`;
    }
  }
  async function uploadHandler() {
    if (!fileInput?.files?.[0]) return;
    let file = fileInput.files[0];
    let formData = new FormData();
    formData.append("file", file);

    let url = `${PUBLIC_API_URL}/files`;
    let resp = await Axios.postForm(url, formData, {
      onUploadProgress: (e) => {
        progress = Math.round((e.loaded / (e.total ?? 1)) * 100);
      },
    });
    let data = resp.data;
    let code = data.Code;

    goto(`/${code}`);
  }
</script>

<!-- upload page -->
<div class="flex flex-col h-screen items-center justify-center">
  <h1 class="m0">Fileshare</h1>

  <div class="min-w-1/2">
    <label
      class="bg-dark w-full h40 flex items-center justify-center p4 rounded relative"
    >
      <div
        class="absolute bg-green-7 top-0 left-0 h-full z1 rounded"
        bind:this={progressBar}
      >
        <!-- progress bar -->
      </div>
      <span class="z-10"> {label} </span>
      <input
        type="file"
        name="file"
        id=""
        class="!hidden"
        disabled={progress > 0}
        bind:this={fileInput}
        on:input={uploadHandler}
      />
    </label>
  </div>
</div>
