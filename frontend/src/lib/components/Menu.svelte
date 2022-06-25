<script lang="ts">
	import type { User } from '$lib/types/User';
	import DarkModeToggler from './DarkModeToggler.svelte';

	export let user: User;

	export let menuItems = [
		['Home', '/'],
		['Games', '/games'],
		['Races', '/races'],
		['Techs', '/techs'],
		[`Welcome, ${user.username}`, '/logout']
	];
</script>

<div class="flex">
	<div class="grow">
		<div class="justify-start ">
			{#each menuItems as item}
				<a
					href={item[1]}
					class="last:float-right rounded-lg px-3 py-2 text-slate-700 dark:text-gray-100 font-medium hover:bg-slate-100 hover:text-slate-900 dark:hover:bg-slate-700 dark:hover:text-gray-300"
					>{item[0]}</a
				>
			{/each}
		</div>
	</div>
	<div class="flex-none">
		<DarkModeToggler />
	</div>
</div>
