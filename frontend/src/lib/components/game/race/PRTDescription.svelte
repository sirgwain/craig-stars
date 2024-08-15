<script lang="ts">
	import { PRT } from '$lib/types/Race';

	export let prt: PRT = PRT.JoaT;
</script>

{#if prt == PRT.HE}
	You must expand to survive. You are given a small and cheap colony hull with an engine that
	travels at Warp 6 using no fuel. Your race will grow at twice the growth rate selected in the
	Habitability Settings menu; however, the maximum population for a given planet is cut in half.
        The completely flexible Meta Morph hull and the beam-empowering Flux Capacitor will be available 
        only to your race, though your expansionist mentality renders you unable to build stargates. 
{:else if prt == PRT.SS}
	You can sneak through enemy territory and execute stunning surprise attacks. You are given access
	to top-drawer cloaks and all your ships have 75% cloaking built in. Cargo does not decrease your
	cloaking abilities. The Stealth Bomber and Rogue hulls are at your disposal, as are a trifecta of
        components with stealthy properties and a pair of scanners which allow you to steal minerals from 
        enemy fleets and planets. You may safely travel through enemy mine fields one warp speed faster 
        than the stated limits, and passively gain research by covertly spying on others' accomplishments. 
{:else if prt == PRT.WM}
	You rule the battle field. Your colonists attack better, your ships are faster in battle, and you
	build weapons 25% cheaper than other races. You start the game with Tech 6 in Weapons and Tech 1 
        in Energy and Propulsion, and can eventually build the powerful Battle Cruiser and Dreadnought hulls. 
        Unfortunately, your race doesn't understand the necessity of building any but the most basic 
        planetary defenses, and cannot build mine layers.
{:else if prt == PRT.CA}
	You are an expert at fiddling with planetary environments. You start the game with Tech 6 in
	Biotechnology and a ship capable of terraforming planets from orbit. You can arm your ships with
	bombs that de-terraform enemy worlds. Terraforming costs you nothing and planets you leave revert
	to their original environments. Planets you own have up to a 10% chance of permanently improving 
        a random environment variable by 1% each year.
{:else if prt == PRT.IS}
	You are strong and hard to defeat. Your colonists repel attacks better, your ships heal faster,
	you have special battle devices that protect your ships, and can lay Speed Trap mine fields to 
        slow attackers. You have a pair of devices that act as both shield and armor. Planetary defenses 
        cost you 40% less, though weapons cost you 25% more. Your colonists are able to reproduce while 
        being transported by your fleets at 50% speed, beaming down excess colonists when orbiting planets 
        you control. However, your peace-loving citizens refuse to build Smart Bombs of any kind.
{:else if prt == PRT.SD}
	You are an expert in laying mine fields. You start the game with 2 mine laying ships and Tech 2 
        in Propulsion and Biotechnology. You have access to a vast array of mine layers, alongside two 
        unique hull designs specializing in mine dispersal. Your mine fields decay slower, detect enemy 
        ships and can be laid while moving, and you have the ability to remotely detonate your own Standard 
        mine fields. You may safely travel two warp speeds faster than the stated limits through enemy mine fields. 
{:else if prt == PRT.PP}
	Your race excels at accelerating mineral packets to distant planets. You start the game with Tech 4 in 
        Energy, a Warp 5 mass driver on your home starbase, and a second planet some distance away if the universe 
        size isn't tiny. You will eventually be able to build accelerators that can fling packets at the 
        mind-numbing speed of warp 16. Your packets are smaller, cheaper, decay slower and contain penetrating 
        scanners whose range depends on their speed. Packets you fling that aren't fully caught have a 
        chance of terraforming the target planet, and may permanently alter it if enough minerals are sent.
{:else if prt == PRT.IT}
	Your race excels in building stargates. You start with Tech 5 in Propulsion and Construction,  
        a small assortment of ships, and a second planet if the universe size isn't tiny. Both planets 
        have stargates. You can send ships with cargo through stargates unimpeded, are less likely to 
        lose ships from exceeding stargate safety limits, and will eventually have access to stargates 
        with unlimited capabilities. Stargates cost you 25% less to build, and automatically scan any 
        enemy planet in range that also has a stargate. However, your people are ill-suited to using 
        mass drivers; your drivers are worse at sending and receiving mineral packets, and any packets 
        built by you will cost more minerals and decay quickly, even ones sent at normally safe speeds.
{:else if prt == PRT.AR}
	Your race developed in an alternate plane. Your people cannot survive on planets and live in orbit
	on your starbases, which are 20% cheaper to build. If a starbase is destroyed, all your colonists 
        living on it are killed. You cannot build planetary installations, but your people have an intrinsic 
        ability to mine and scan for enemy fleets, and you can remote mine your own worlds for more minerals. 
        Your population maximums are determined solely by the type of starbase you have, while your population 
        resource production depends on both habitability and your Energy Tech level. You will eventually be able 
        to build the Death Star, the biggest, baddest starbase to ever exist.
{:else if prt == PRT.JoaT}
	Your race does not specialize in a single area. You start the game with Tech 3 in all fields and a wide
	assortment of ships. Your Scout, Destroyer, and Frigate hulls have a built-in penetrating scanner
	that grows more powerful as your Electronics tech increases. Your maximum planetary population is 
        20% greater than other races. 
{:else}
	{prt.toString()}
{/if}
