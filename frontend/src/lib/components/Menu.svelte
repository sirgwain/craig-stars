<script lang="ts">
	import type { User } from '$lib/types/User';
	import { Menu } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import DarkModeToggler from './DarkModeToggler.svelte';

	export let user: User | undefined;
</script>

<div class="navbar bg-base-100 flex flex-row">
	<div class="flex-1">
		<a class="btn btn-ghost text-xl text-accent" href="/"
			><span class="hidden md:block">craig-stars</span><span class="md:hidden">cs</span></a
		>
	</div>
	<div class="flex-none">
		<div class="hidden md:block flex-none">
			<ul class="menu menu-horizontal p-0">
				{#if user}
					<li><a href="/games">Games</a></li>
					<li><a href="/races">Races</a></li>
				{/if}
				<li><a href="/techs">Techs</a></li>
			</ul>
		</div>

		<div class="dropdown dropdown-end">
			<label for="menu" tabindex="0" class="btn btn-ghost">
				<div id="menu">
					<Icon src={Menu} size="24" />
				</div>
			</label>
			<ul
				tabindex="0"
				class="menu menu-compact dropdown-content mt-3 p-2 shadow bg-base-100 rounded-box w-30"
			>
				{#if user}
					<li>{user.username}</li>
				{/if}
				<li>
					<DarkModeToggler />
				</li>
				{#if user}
					<li class="md:hidden"><a href="/games">Games</a></li>
					<li class="md:hidden"><a href="/races">Races</a></li>
					<li class="md:hidden"><a href="/techs">Techs</a></li>
					<li><a href="/logout">Logout</a></li>
				{:else}
					<li class="md:hidden"><a href="/techs">Techs</a></li>
				{/if}
			</ul>
		</div>
	</div>
</div>
