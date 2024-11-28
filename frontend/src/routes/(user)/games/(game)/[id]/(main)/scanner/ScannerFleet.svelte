<script lang="ts">
	import { getGameContext } from '$lib/services/GameContext';
	import { radiansToDegrees } from '$lib/services/Math';
	import type { Fleet } from '$lib/types/Fleet';
	import { ownedBy } from '$lib/types/MapObject';
	import MapObjectScaler from './MapObjectScaler.svelte';

	const { settings, player } = getGameContext();

	type Props = {
		fleet: Fleet;
		commanded?: boolean;
		color?: string;
		commandedColor?: string;
	};

	let { fleet, commanded = false, color = '#0000FF', commandedColor = '#FFFF00' }: Props = $props();

	const size = 8;

	// identity or default is rotated 90º, or pointing up and to the right
	const angleOffset = 225;

	function getAngle(fleet: Fleet): number {
		return radiansToDegrees(Math.atan2(fleet.heading.y, fleet.heading.x)) + angleOffset;
	}

	function getTokenCount(fleet: Fleet): number {
		return fleet.tokens ? fleet.tokens.reduce((count, t) => count + t.quantity, 0) : 0;
	}

	let textColor = $derived(
		ownedBy(fleet, $player.num)
			? 'fill-orbit'
			: $player.isFriend(fleet.playerNum)
				? 'fill-orbit-friends'
				: 'fill-orbit-enemies'
	);
</script>

<!-- ScannerFleet -->
<MapObjectScaler mapObject={fleet}>
	<polygon
		points={`0,0 0,${size} ${size},${size}`}
		fill={commanded ? commandedColor : color}
		transform={`rotate(${getAngle(fleet)}) translate(${-size / 2} ${-size / 2})`}
	/>
	{#if $settings.showFleetTokenCounts}
		<!-- position the text below the fleet -->
		<text transform={`translate(0 ${size * 2.5})`} text-anchor="middle" class={textColor}
			>{getTokenCount(fleet)}</text
		>
	{/if}
</MapObjectScaler>
