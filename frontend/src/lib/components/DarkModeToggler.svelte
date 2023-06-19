<script lang="ts">
	import { Moon, Sun } from 'radix-icons-svelte';
	import { onMount } from 'svelte';

	let isDark = false;
	onMount(() => {
		isDark =
			localStorage.theme === 'dark' ||
			(!('theme' in localStorage) && window.matchMedia('(prefers-color-scheme: dark)').matches);
	});

	function toggleTheme() {
		isDark = !isDark;
		localStorage.theme = isDark ? 'dark' : 'light';
		if (isDark) {
			document.documentElement.classList.add('dark');
		} else {
			document.documentElement.classList.remove('dark');
		}
	}
</script>

<button on:click={toggleTheme} class="pt-2 pb-3">
	{#if isDark}
		<Sun size={16} />
	{:else}
		<Moon size={16} />
	{/if}
</button>
