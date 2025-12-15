<script lang="ts">
	import { onMount } from 'svelte';
	import { ChevronUp, File, Folder } from '@lucide/svelte';

  type File = {
    Name: string;
    Key: string;
  }

	let currentDir = '/'; // starting directory
	let files: File[] = []; // files in the directory
	let error: string | null = null;

	async function loadFiles() {
		error = null;

		let cleanDir = currentDir.replace(/^\//, ''); // remove leading slash

		try {
			const res = await fetch(`http://127.0.0.1:3000/files/list/${cleanDir}`, {
				credentials: 'include'
			});

			if (!res.ok) {
				error = 'Could not fetch files';
				return;
			}

			files = await res.json();
		} catch (e) {
			error = 'Network error';
		}
	}

	// load on page mount
	onMount(loadFiles);

	// reload whenever directory changes
	$: if (currentDir !== undefined) {
		loadFiles();
	}

	function enterDirectory (name: string) {
		currentDir = currentDir === '/' ? `/${name}` : `${currentDir}/${name}`;
	}

	async function handleUpload(e: Event) {
    e.preventDefault();

    const input = e.target as HTMLInputElement;
    if (!input || !input.files || input.files.length === 0) {
        return;
    }

    const file = input.files[0];

    const formData = new FormData();
    formData.append("file", file);

    try {
			const res = await fetch("http://127.0.0.1:3000/files/upload", {
				method: "POST",
				credentials: "include",
				body: formData,
			});

			if (!res.ok) {
				error = "Could not upload file";
				return;
			}

			await loadFiles();
			
    } catch (err) {
			console.error(err);
			error = "Network error";
    }
	}

	function goUpOneDirectory() {
		if (currentDir === '/') return;

		const parts = currentDir.split('/').filter(Boolean);
		parts.pop();
		currentDir = parts.length ? `/${parts.join('/')}` : '/';
	}
</script>

<div class="flex items-center justify-center h-full w-full">
	<div class="w-2xl h-full bg-stone-300 p-4 border-3">
		<div class="w-full">
			<input value={currentDir}/>
		</div>
		<div class="flex justify-between">
			<form>
				<input type="file" on:change={handleUpload} class="hover:underline cursor-pointer" />
			</form>
		</div>

		<p>Current: <strong>{currentDir}</strong></p>
		<!-- go up a level -->
		<button class="flex" on:click={goUpOneDirectory} disabled={currentDir === '/'}> <ChevronUp /> Up </button>

		{#if error}
			<p style="color:red">{error}</p>
		{/if}

		<div class="flex gap-4 p-2 items-start h-fit">
			{#each files as file}
				<div class="w-16 hover:bg-stone-500 ">
					<div class="flex w-full justify-center">
						{#if !file.Name.includes('.')}
							<button on:click={() => enterDirectory(file.Name)}>
								<Folder class="stroke-[0.5px] w-8 h-8" />
							</button>
						{:else}
							<File class="stroke-[0.5px] w-8 h-8"/>
						{/if}
					</div>
					<p class="text-[9px] wrap-break-word text-center">{file.Name}</p>
				</div>
			{/each}
		</div>
	</div>
</div>
