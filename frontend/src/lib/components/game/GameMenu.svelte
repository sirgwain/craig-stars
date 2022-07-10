<script lang="ts">
	import { Icon } from '@steeze-ui/svelte-icon';
	import { Menu, Sun } from '@steeze-ui/heroicons';
	import { game } from '$lib/services/Context';
	import { createEventDispatcher } from 'svelte';
	import DarkModeToggler from '../DarkModeToggler.svelte';

	const updateTitle = () => (document.title = `${$game.name} - ${$game.year}`);
	$: $game && updateTitle();
	const dispatch = createEventDispatcher();
</script>

{#if $game}
	<div class="navbar">
		<div class="flex-1">
			<a class="btn btn-ghost text-xl text-accent" href={`/`}>cs</a>
			<a class="btn btn-ghost text-xl text-primary" href={`/games/${$game.id}`}
				>{$game.name} - {$game.year}</a
			>
		</div>
		<div class="flex-none">
			<div class="dropdown">
				<label for="reports" tabindex="0" class="btn btn-ghost">Reports</label>
				<ul
					id="reports"
					tabindex="0"
					class=" menu menu-compact dropdown-content mt-3 p-2 shadow bg-base-100 rounded-box w-52"
				>
					<li><a href={`/games/${$game.id}/planets`}>Planets</a></li>
					<li><a href={`/games/${$game.id}/fleets`}>Fleets</a></li>
					<li><a href={`/games/${$game.id}/designs`}>Designs</a></li>
					<li><a href={`/games/${$game.id}/messages`}>Messages</a></li>
				</ul>
			</div>
			<div class="dropdown">
				<button on:click={() => dispatch('submit-turn')} class="btn btn-primary">Submit Turn</button
				>
			</div>
			<div class="dropdown dropdown-end">
				<label for="menu" tabindex="0" class="btn btn-ghost">
					<div id="menu">
						<Icon src={Menu} size="24" />
					</div>
				</label>
				<ul
					tabindex="0"
					class="menu menu-compact dropdown-content mt-3 p-2 shadow bg-base-100 rounded-box w-52"
				>
					<li>
						<DarkModeToggler />
					</li>
					<li>
						<a href={`/games/${$game.id}/race`} class="justify-between">Race</a>
					</li>
					<li>
						<a href={`/games/${$game.id}/techs`} class="justify-between">Techs</a>
					</li>
					<li><a href="/logout">Logout</a></li>
				</ul>
			</div>
		</div>
	</div>
{/if}
