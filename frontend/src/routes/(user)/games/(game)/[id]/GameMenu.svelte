<script lang="ts">
	import { game } from '$lib/services/Context';
	import { Menu, Upload } from '@steeze-ui/heroicons';
	import { Icon } from '@steeze-ui/svelte-icon';
	import { createEventDispatcher } from 'svelte';
	import DarkModeToggler from '$lib/components/DarkModeToggler.svelte';

	const updateTitle = () => (document.title = `${$game?.name} - ${$game?.year}`);
	$: $game && updateTitle();
	const dispatch = createEventDispatcher();
</script>

{#if $game}
	<div class="navbar bg-base-100 flex flex-row">
		<div class="flex-1">
			<a class="btn btn-ghost text-xl text-primary" href={`/`}>cs</a>
			<div class="md:block">
				<a class="btn btn-ghost text-lg text-accent" href={`/games/${$game.id}`}
					>{$game.name} - {$game.year}</a
				>
			</div>
		</div>
		<div class="flex-initial">
			<div class="hidden md:inline-block">
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
			</div>

			<div class="tooltip tooltip-bottom" data-tip="submit turn">
				<button on:click={() => dispatch('submit-turn')} class="btn btn-primary" title="submit turn"
					><span class="hidden md:inline-block mr-1">Submit Turn</span><Icon
						src={Upload}
						size="16"
					/></button
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
					<li class="md:hidden menu-title">
						<span>Reports</span>
					</li>
					<li class="md:hidden"><a href={`/games/${$game.id}/planets`}>Planets</a></li>
					<li class="md:hidden"><a href={`/games/${$game.id}/fleets`}>Fleets</a></li>
					<li class="md:hidden"><a href={`/games/${$game.id}/designs`}>Designs</a></li>
					<li class="md:hidden"><a href={`/games/${$game.id}/messages`}>Messages</a></li>

					<li class="md:hidden menu-title">
						<span>Game</span>
					</li>
					<li>
						<a href={`/games/${$game.id}/race`} class="justify-between">Race</a>
					</li>
					<li>
						<a href={`/games/${$game.id}/techs`} class="justify-between">Techs</a>
					</li>
					<li><div class="divider" /></li>
					<li><a href="/logout">Logout</a></li>
				</ul>
			</div>
		</div>
	</div>
{/if}