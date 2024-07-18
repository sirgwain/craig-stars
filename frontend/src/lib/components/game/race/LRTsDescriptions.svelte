<script lang="ts">
	import { getLabelForLRT, LRT, lrts, type Race } from '$lib/types/Race';

	export let race: Race;

	const getDescriptionForLRT = (lrt: LRT) => {
		switch (lrt) {
			case LRT.IFE:
				return 'This gives you the Fuel Mizer and Galaxy Scoop engines and increases your starting Propulsion tech by 1 level. All engines use 15% less fuel.';
			case LRT.TT:
				return 'You will gain access to additional terraforming technologies by researching Biotechnology, starting with a 3% adjustment to all 3 environmental attributes and reaching a maximum of 30%. Terraforming requires 30% fewer resources.'
			case LRT.ARM:
				return 'Gives you three additional mining hulls and two new robots. You will start the game with two Midget Miners.';
			case LRT.ISB:
				return 'Gives you two new starbase designs. The Space Dock is a cheap starbase with a small space dock that can refuel ships and build small to medium sized vessels. The Ultra-Station is much larger and better fortified than a standard Starbase. Your starbases are automatically cloaked by 20% and cost 20% less to build.';
			case LRT.GR:
				return 'Your race takes a holistic approach to research. Only half of the resources dedicated to research will be applied to the currently selected field, but 15% of the amount will be applied to all other fields (for a total of 125%).';
			case LRT.UR:
				return 'When you scrap a fleet at a starbase, you will recover 90% of the minerals and up to 70% of the resources used to produce the fleet, avaliable for use the following year. Scrapping at a planet will yield half the amount of minerals and resources.';
			case LRT.NRSE:
				return 'You will not be able to build the Radiating Hydro-Ram Scoop, Sub-Galactic Fuel Scoop, Trans-Galactic Fuel Scoop, Trans-Galactic Super Scoop, Trans-Galactic Mizer Scoop or the Galaxy Scoop. You will be able to build the Interspace-10 engine, which can safely travel at warp 10 and has modest technology requirements.';
			case LRT.OBRM:
				return "The only mining ship available to you will be the Mini-Miner. This trait overrides 'Advanced Remote Mining'. Your maximum population per planet is increased by 10%.";
			case LRT.NAS:
				return 'You will not have access to any penetrating scanners that can scan planets from a distance and see fleets hiding behind planets. However, all ranges for conventional scanners are doubled.';
			case LRT.LSP:
				return 'Instead of 25000 people, you start with 17500 (30% fewer). It takes a long time to overcome a lower starting population - it helps to have a high growth rate, but even then it can be painful.';
			case LRT.BET:
				return 'New technologies initially cost twice as much to build. As soon as you exceed all of the tech requirements by one level, the cost drops back to normal. Miniaturization (the lowering of production costs at higher tech levels) occurs at 5% per level, up to 80%. (Without this trait, miniaturization occurs at 4% per level and caps at 75%.)';
			case LRT.RS:
				return 'All shields are 40% stronger than their listed rating. Shields regenerate 10% of their maximum strength after each round of battle until destroyed. However, any armor components added to your ships will only have 50% of their rated strength.';
			case LRT.MA:
				return 'You will be able to turn resources into minerals more efficiently. One instance of mineral alchemy will use 25 resources (instead of 100) to produce 1kT of each mineral.';
			case LRT.CE:
				return 'Engines cost you 50% less to build, but aren’t entirely reliable. When attempting to travel at speeds above warp 6, there is a 10% chance the engines will refuse to engage. Your starting Propulsion tech level will be increased by 1.';
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
