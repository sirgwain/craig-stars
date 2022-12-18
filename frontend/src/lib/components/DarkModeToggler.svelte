<script lang="ts">
	import { Icon } from '@steeze-ui/svelte-icon';
	import { Moon, Sun } from '@steeze-ui/heroicons';
	import { onMount } from 'svelte';

	const lightTheme = 'emerald'
	const darkTheme = 'business'
	let isDark = false;
	onMount(() => {
		isDark =
			localStorage.theme === darkTheme ||
			(!('theme' in localStorage) && window.matchMedia('(prefers-color-scheme: dark)').matches);
	});

	function toggleTheme() {
		isDark = !isDark;
		localStorage.theme = isDark ? darkTheme : lightTheme;
		if (isDark) {
			document.documentElement.dataset.theme = darkTheme;
			document.documentElement.className = 'dark';
		} else {
			document.documentElement.dataset.theme = lightTheme;
			document.documentElement.className = 'light';
		}
	}
</script>

<button on:click={toggleTheme} class="pt-2 pb-3" type="button">
	{#if isDark}
		<Icon src={Sun} size="24" class="hover:stroke-accent" />
	{:else}
		<Icon src={Moon} size="24" class="hover:stroke-accent" />
	{/if}
</button>
