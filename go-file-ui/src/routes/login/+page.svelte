<script lang="ts">
  import { goto } from '$app/navigation';
  import { writable } from 'svelte/store';

  let email = '';
  let password = '';
  let error = '';

  const loading = writable(false);

  async function handleLogin() {
    loading.set(true);
    error = '';

    try {
      const res = await fetch('http://127.0.0.1:3000/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password }),
        credentials: 'include',
      });

      if (!res.ok) {
        const text = await res.text();
        error = text || 'Login failed';
        loading.set(false);
        return;
      }

      goto('/dashboard');
    } catch (err) {
      error = 'Network error';
      console.error(err);
    } finally {
      loading.set(false);
    }
  }
</script>

<div class="h-full flex flex-col justify-center bg-neutral-700 items-center">
  <div class="relative flex flex-col px-6 w-md bg-stone-200 items-center rounded-sm border-3 p-4 gap-4">
    <h1 class="text-2xl font-bold">Login</h1>
    <form
      class="flex flex-col gap-4 items-center w-full"
      on:submit|preventDefault={handleLogin}
    >
    <div class="flex flex-col gap-2 w-full">
      <input class="border-3 p-1 font-semibold rounded-sm" type="email" bind:value={email} placeholder="Email" required />
      <input class="border-3 p-1 font-semibold rounded-sm" type="password" bind:value={password} placeholder="Password" required />
    </div>
      <button class="cursor-pointer font-semibold rounded-sm border-3 w-fit p-1" type="submit" disabled={$loading}>
        {$loading ? 'Logging in...' : 'Login'}
      </button>
    </form>
    {#if error}
      <p class="text-red-500 absolute -bottom-8">{error}</p>
    {/if}
  </div>
</div>
