import { writable } from 'svelte/store';

export const isLoggedIn = writable(false);
export const user = writable(null);

export async function checkAuth() {
  const res = await fetch('http://127.0.0.1:3000/auth/me', {
    credentials: 'include'
  });

  if (res.ok) {
    user.set(await res.json());
    isLoggedIn.set(true);
  } else {
    user.set(null);
    isLoggedIn.set(false);
  }
}