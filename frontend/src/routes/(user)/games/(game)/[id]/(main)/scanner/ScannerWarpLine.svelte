<script lang="ts">
	import { MapObjectType, type MovingMapObject } from '$lib/types/MapObject';

	import { getGameContext } from '$lib/services/GameContext';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';
	import type { Writable } from 'svelte/store';

	const { player, universe, selectedMapObject } = getGameContext();
	const scale = getContext<Writable<number>>('scale');
	const { data, xGet, yGet, xScale, yScale, width, height } = getContext<LayerCake>('LayerCake');

	type Line = {
		path: string;
		props: any;
	};

	let line: Line | undefined;
	let color = '#ffffff';

	$: strokeWidth = 3 / $scale;

	$: {
		line = undefined;

		// show the warp line for other player fleets, or mystery traders or mineral packets
		if (
			$selectedMapObject &&
			($selectedMapObject.type == MapObjectType.MineralPacket ||
				$selectedMapObject.type == MapObjectType.MysteryTrader ||
				($selectedMapObject.type == MapObjectType.Fleet &&
					$selectedMapObject.playerNum != $player.num))
		) {
			const mo = $selectedMapObject as MovingMapObject;
			const heading = mo.heading ?? { x: 0, y: 0 };
			const warpSpeed = mo.warpSpeed ?? 0;
			const distPerLy = warpSpeed * warpSpeed;
			if (mo.playerNum) {
				color = $universe.getPlayerColor(mo.playerNum);
			} else if (mo.type == MapObjectType.MysteryTrader) {
				color = '#00FFFF'
			}

			if (warpSpeed) {
				const coords = [-5, -4, -3, -2, -1, 1, 2, 3, 4, 5].map((dist: number) => ({
					position: {
						x: mo.position.x + heading.x * distPerLy * dist,
						y: mo.position.y + heading.y * distPerLy * dist
					}
				}));

				line = {
					path: 'M' + coords.map((coord) => `${$xGet(coord)}, ${$yGet(coord)}`).join('L'),
					props: {
						'stroke-width': strokeWidth,
						stroke: color
					}
				};
			}
		}
	}
</script>

{#if line}
	<svg>
		<defs>
			<marker
				id="warp-arrow"
				class="warpline-arrow"
				viewBox="0 0 10 10"
				refX="13"
				refY="5"
				markerUnits="strokeWidth"
				markerWidth="3"
				markerHeight="3"
				orient="auto"
			>
				<path d="M 3 0 L 7 5 L 3 10" stroke={color} fill="context-fill" stroke-width={2} />
			</marker>
		</defs>
	</svg>
	<path d={line.path} {...line.props} marker-mid="url(#warp-arrow)" />
{/if}
