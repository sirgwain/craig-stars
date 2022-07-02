<script lang="ts">
	import type { User } from '$lib/types/User';
	import DarkModeToggler from './DarkModeToggler.svelte';

	export let user: User;
</script>

<div class="navbar bg-base-100">
	<div class="navbar-start">
		<div class="flex-1">
			<a class="btn btn-ghost text-xl text-accent" href="/">craig-stars</a>
		</div>
		<div class="flex-none">
			<ul class="menu menu-horizontal p-0">
				<li><a href="/games">Games</a></li>
				<li><a href="/races">Races</a></li>
				<li><a href="/techs">Techs</a></li>
			</ul>
		</div>
	</div>
	<div class="navbar-end">
		{#if user}
			<div class="mr-2"><a class="link" href="/logout">{user.username}</a></div>
		{/if}
		<DarkModeToggler />
	</div>
</div>
