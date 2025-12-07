<script lang="ts">
	import './layout.css';
  import { isLoggedIn, user } from '$lib/stores/user';
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
	import { derived } from 'svelte/store';

  onMount(() => {
    if (!isLoggedIn && window.location.pathname !== '/login') {
      goto('/login');
    }
  });
</script>




<main class="flex flex-col p-2 h-screen font-semibold bg-neutral-700 font-mono">
	<header class="px-4 bg-neutral-700">
		{#if $isLoggedIn}
			<nav class="flex justify-end">
				<!-- <button class="border p-1 bg-stone-200" on:click={() => { token.set(null); user.set(null); localStorage.removeItem('token'); goto('/login'); }}>
					Logout
				</button> -->
			</nav>
		{/if}
	</header>
	<div class="h-full">
		<slot /> <!-- This renders the child page content -->
	</div>
	<footer class="text-center text-stone-400 mx-32 border-t">
  	&copy; 2025 Go File API
	</footer>
</main>


