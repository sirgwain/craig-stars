<script lang="ts">
	import { Icon } from '@steeze-ui/svelte-icon';
	import { Moon, Sun } from '@steeze-ui/heroicons';
	import { onMount } from 'svelte';

	let isDark = false;
	onMount(() => {
		isDark =
			localStorage.theme === 'business' ||
			(!('theme' in localStorage) && window.matchMedia('(prefers-color-scheme: dark)').matches);
	});

	function toggleTheme() {
		isDark = !isDark;
		localStorage.theme = isDark ? 'business' : 'emerald';
		if (isDark) {
			document.documentElement.dataset.theme = 'business';
			document.documentElement.className = 'dark';
		} else {
			document.documentElement.dataset.theme = 'emerald';
			document.documentElement.className = 'light';
		}
	}
</script>

<button on:click={toggleTheme} class="pt-2 pb-3">
	{#if isDark}
		<Icon src={Sun} size="16" class="hover:stroke-accent" />
	{:else}
		<Icon src={Moon} size="16" class="hover:stroke-accent" />
	{/if}
</button>
