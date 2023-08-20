<script lang="ts">
	import { getLabelForLRT, LRT, lrts, type Race } from '$lib/types/Race';

	export let race: Race;

	const getDescriptionForLRT = (lrt: LRT) => {
		switch (lrt) {
			case LRT.IFE:
				return 'This gives you the Fuel Mizer and Galaxy Scoop engines and increases your starting Propulsion tech by 1 level. All engines use 15% less fuel.';
			case LRT.TT:
				return 'You begin the game with the ability to adjust each of a planet’s environment attributes by up to 3% in either direction. Throughout the game, additional terraforming technologies not available to other players will be achievable, up to 30% terraforming. Total Terraforming requires 30% fewer resources.';
			case LRT.ARM:
				return 'Gives you three additional mining hulls and two new robots. You will start the game with two Midget Miners.';
			case LRT.ISB:
				return 'Gives you two new starbase designs. The Space Dock hull allows you to build starbases which can in turn build small to medium ships. The Ultra-Station is much larger than a standard Starbase. Your starbases are automatically cloaked by 20%. Starbases will cost you 20% less to build.';
			case LRT.GR:
				return 'Your race takes a holistic approach to research. Only half of the resources dedicated to research will be applied to the current field of research. 15% of the total will be applied to each of the fields.';
			case LRT.UR:
				return 'When you scrap a fleet at a starbase, you recover 90% of the minerals and 70% of the resources used to produce the fleet. Scrapping at a planet gives you 45% of the minerals and 35% of the resources.';
			case LRT.NRSE:
				return 'You will not be able to build the Radiating Hydro-Ram Scoop, Sub-Galactic Fuel Scoop, Trans-Galactic Fuel Scoop, Trans-Galactic Super Scoop, Trans- Galactic Mizer Scoop or the Galaxy Scoop. You will be able to build the Interspace-10 engine, which can travel warp 10 without taking damage.';
			case LRT.OBRM:
				return 'No Robo-Miner, Robo-Maxi-Miner or Robo-Super-Miner robots.';
			case LRT.NAS:
				return 'You will not have any standard scanners that can scan planets from a distance and see fleets hiding behind planets. All ranges for conventional scanners are doubled.';
			case LRT.LSP:
				return 'Instead of 25000 people, you start with 17500 (30% fewer). It takes a long time to overcome a lower starting population - it helps to have a high growth rate, but even then it can be painful.';
			case LRT.BET:
				return 'New technologies initially cost twice as much to build. As soon as you exceed all of the tech requirements by one level, the cost drops back to normal. Miniaturization, the lowering of production costs, occurs at 5% per level, up to 80%. Without this trait miniaturization occurs at 4% and tops out at 75%.';
			case LRT.RS:
				return 'All shields are 40% stronger than the listed rating. Shields regenerate at 10% of the maximum strength after every round of battle. All armors are at 50% of their rated strength.';
			case LRT.MA:
				return 'You will be able to turn resources into minerals more efficiently. One instance of mineral alchemy will use 25 resources to produce one kT of each mineral. Without this trait it takes 100 resources to produce one kT of each mineral.';
			case LRT.CE:
				return 'Engines cost 50% less to build, your ship engines aren’t entirely reliable. When attempting to travel at speeds above warp 6, there is a 10% chance the engines will refuse to engage.';
			default:
				return lrt.toString();
		}
	};
</script>

<div class="flex flex-row flex-wrap justify-between gap-1">
	{#each lrts as lrt}
		{#if race.lrts & lrt}
			<div class="card bg-base-200 shadow w-full">
				<div class="card-body">
					<div class="card-title text-lg">
						{getLabelForLRT(lrt)}
					</div>
					<div>{getDescriptionForLRT(lrt)}</div>
				</div>
			</div>
		{/if}
	{/each}
</div>
