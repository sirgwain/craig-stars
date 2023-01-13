<script lang="ts">
	import { PRT } from '$lib/types/Race';

	export let prt: PRT = PRT.JoaT;
</script>

{#if prt == PRT.HE}
	You must expand to survive. You are given a small and cheap colony hull and an engine which
	travels at Warp 6 using no fuel. Your race will grow at twice the growth rate you select in step
	four; however, the maximum population for a given planet is cut in half. The completely flexible
	Meta Morph hull will be available only to your race
{:else if prt == PRT.SS}
	You can sneak through enemy territory and execute stunning surprise attacks. You are given
	top-drawer cloaks and all your ships have 75% cloaking built in. Cargo does not decrease your
	cloaking abilities. The Stealth Bomber and Rogue are at your disposal, as are a scanner, shield
	and armor with stealthy properties. Two scanners which allow you to steal minerals from enemy
	fleets and planets are also available. You may safely travel through mine fields at one warp speed
	faster than the limits
{:else if prt == PRT.WM}
	You rule the battle field. Your colonists attack better, your ships are faster in battle, and you
	build weapons 25% cheaper than other races.You start the game with a knowledge of Tech 6 weapons
	and Tech 1 in enerqy and propulsion. Unfortunately, your race doesn't understand the necessity of
	building any but the most basic planetary defenses and no mine fields
{:else if prt == PRT.CA}
	You are an expert at fiddling with planetary environments. You start the game with Tech 6 in
	Biotechnology and a ship capable of terraforming planets from orbit. You can arm your ships with
	bombs that unterraform enemy worlds. Terraforming costs you nothing and planets you leave revert
	to their original environments. Planets you own have up to a 10% chance of permanently improving
	an environment variable by 1% per vear
{:else if prt == PRT.IS}
	You are strong and hard to defeat. Your colonists repel attacks better, your ships heal faster,
	you have special battle devices that protect your ships and can lay Speed Trap mine fields. You
	have a device that acts as both a shield and armor. Your peace-loving people refuse to build Smart
	Bombs. Planetary defenses cost you 40% less, though weapons cost you 25% more. Your colonists are
	able to reproduce while being transported by your fleets
{:else if prt == PRT.SD}
	You are an expert in laying mine fields. You have a vast array of mine types at your disposal and
	two unique hull design which are made for mine dispersal. Your mine fields act as scanners and you
	have the ability to remote detonate your own Standard mine fields. You may safely travel two warp
	speeds faster than the stated limits through enemy mine fields. You start the game with 2 mine
	laying ships and Tech 2 in Propulsion and BioTech
{:else if prt == PRT.PP}
	Your race excels at accelerating mineral packets to distant planets. You start with a Warp 5
	accelerator at your home starbase and Tech 4 Energy. You will eventually be able to fling packets
	at the mind numbing speed of Warp 13. You can fling smaller packets and all of your packets have
	penetrating scanners embedded in them. You will start the game owning a second planet some
	distance away if the universe size isn't tiny. Packets you fling that aren't fully caught have a
	chance of terraforming the target planet
{:else if prt == PRT.IT}
	Your race excels in building stargates. You start with Tech 5 in Propulsion and Construction. You
	start the game with a second planet if the universe size isn't tiny. Both planets have stargates.
	Eventually you may build stargates which have unlimited capabilities. Stargates cost you 25% less
	to build. Your race can automatically scan any enemy planet with a stargate which is in range of
	one of your stargates. Exceeding the safety limits of stargates is less likely to kill your ships
{:else if prt == PRT.AR}
	Your race developed in an alternate plane. Your people cannot survive on planets and live in orbit
	on your starbases, which are 20% cheaper to build. You cannot build planetary installations, but
	your people have an intrinsic ability to mine and scan for enemy fleets. You can remote mine your
	own worlds. If a starbase is destroyed, all your colonists orbiting that world are killed. Your
	population maximums are determined by the type of starbase you have. You will eventually be able
	to build the Death Star
{:else if prt == PRT.JoaT}
	Your race does not specialize in a single area. You start the game with Tech 3 in all areas and an
	assortment of ships. Your Scout, Destroyer, and Frigate hulls have a built-in penetrating scanner
	which grows more powerful as your Electronics tech increases. Your maximum planetary population is
	20% greater than other races
{:else}
	{prt.toString()}
{/if}
