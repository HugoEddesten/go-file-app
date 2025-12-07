<script lang="ts">
	import { onMount } from 'svelte';
  import { Folder } from '@lucide/svelte';

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

	function enterDirectory(name: string) {
		currentDir = currentDir === '/' ? `/${name}` : `${currentDir}/${name}`;
	}

	function goUpOneDirectory() {
		if (currentDir === '/') return;

		const parts = currentDir.split('/').filter(Boolean);
		parts.pop();
		currentDir = parts.length ? `/${parts.join('/')}` : '/';
	}
</script>

<div class="flex items-center justify-center h-full w-full">

	<div class="w-2xl min-h-4/5 bg-stone-300 p-4 border-3">
			<h1>Your Files</h1>
		
			<p>Current: <strong>{currentDir}</strong></p>
		
			<!-- go up a level -->
			<button on:click={goUpOneDirectory} disabled={currentDir === '/'}> ⬆ Up </button>
		
			{#if error}
				<p style="color:red">{error}</p>
			{/if}
		
			<div class="flex gap-2 p-2">
				{#each files as file}
					<div class="flex items-center justify-center box-content border-3 h-32 w-32 overflow-clip">
						{#if !file.Name.includes('.')}
							<button on:click={() => enterDirectory(file.Name)}><Folder class="stroke-[0.5px] w-32 h-32"/></button>
						{:else}
							<img
								class=" min-w-full min-h-full"
								src="http://127.0.0.1:3000/files/download/{file.Key}"
								alt={file.Name}
							/>
						{/if}
					</div>
				{/each}
			</div>
		</div>
	</div>
