<script lang="ts">
	import { type MovingMapObject } from '$lib/types/MapObject';

	import { getGameContext } from '$lib/services/GameContext';
	import type { LayerCake } from 'layercake';
	import { getContext } from 'svelte';
	import type { SVGAttributes } from 'svelte/elements';
	import { MapObjectType } from '$lib/types/cs-proto';

	const { player, universe, selectedMapObject } = getGameContext();
	const { xGet, yGet } = getContext<LayerCake>('LayerCake');

	type Line = {
		color: string;
		path: string;
		props: SVGAttributes<SVGPathElement>;
	};

	let line: Line | undefined = $derived.by(() => {
		let color = '#ffffff';
		let strokeWidth = 1;

		// show the warp line for other player fleets, or mystery traders or mineral packets
		if (
			$selectedMapObject &&
			($selectedMapObject.mapObject?.type === MapObjectType.MINERAL_PACKET ||
				$selectedMapObject.mapObject?.type === MapObjectType.MYSTERY_TRADER ||
				($selectedMapObject.mapObject?.type === MapObjectType.FLEET &&
					$selectedMapObject.mapObject.playerNum != $player.num))
		) {
			const mo = $selectedMapObject as MovingMapObject;
			const heading = mo.heading;
			const warpSpeed = mo.warpSpeed;
			const distPerLy = warpSpeed * warpSpeed;
			if (mo.mapObject?.playerNum) {
				color = $universe.getPlayerColor(mo.mapObject.playerNum);
			} else if (mo.mapObject?.type === MapObjectType.MYSTERY_TRADER) {
				color = '#00FFFF';
			}

			if (warpSpeed) {
				const coords = [-5, -4, -3, -2, -1, 1, 2, 3, 4, 5].map((dist: number) => ({
					position: {
						x: (mo.mapObject?.position?.x ?? 0) + heading.x * Math.ceil(distPerLy * dist),
						y: (mo.mapObject?.position?.y ?? 0) + heading.y * Math.ceil(distPerLy * dist)
					}
				}));

				return {
					color,
					path: 'M' + coords.map((coord) => `${$xGet(coord)}, ${$yGet(coord)}`).join('L'),
					props: {
						'stroke-width': strokeWidth,
						stroke: color
					}
				};
			}
		}
	});
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
				<path d="M 3 0 L 7 5 L 3 10" stroke={line.color} fill="context-fill" stroke-width={2} />
			</marker>
		</defs>
	</svg>
	<path d={line.path} {...line.props} marker-mid="url(#warp-arrow)" />
{/if}
