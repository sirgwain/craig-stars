<script lang="ts">
	import {
		commandMapObject,
		commandedFleet,
		commandedMapObject,
		commandedPlanet,
		selectMapObject,
		selectWaypoint,
		selectedMapObject,
		selectedWaypoint,
		zoomTarget
	} from '$lib/services/Context';
	import type { FullGame } from '$lib/services/FullGame';
	import { PlanetService } from '$lib/services/PlanetService';
	import { settings } from '$lib/services/Settings';
	import { WaypointTask } from '$lib/types/Fleet';
	import { MapObjectType, None, ownedBy, type MapObject } from '$lib/types/MapObject';
	import { emptyVector, equal, type Vector } from '$lib/types/Vector';
	import type { ScaleLinear } from 'd3-scale';
	import { scaleLinear } from 'd3-scale';
	import { select } from 'd3-selection';
	import { ZoomTransform, zoom, type D3ZoomEvent, type ZoomBehavior } from 'd3-zoom';
	import { Html, LayerCake, Svg } from 'layercake';
	import { setContext } from 'svelte';
	import { writable } from 'svelte/store';
	import MapObjectQuadTreeFinder from './MapObjectQuadTreeFinder.svelte';
	import ScannerFleets from './ScannerFleets.svelte';
	import ScannerMineFieldPattern from './ScannerMineFieldPattern.svelte';
	import ScannerMineFields from './ScannerMineFields.svelte';
	import ScannerMineralPackets from './ScannerMineralPackets.svelte';
	import ScannerNames from './ScannerNames.svelte';
	import ScannerPacketDests from './ScannerPacketDests.svelte';
	import ScannerPlanets from './ScannerPlanets.svelte';
	import ScannerScanners from './ScannerScanners.svelte';
	import ScannerWarpLine from './ScannerWarpLine.svelte';
	import ScannerWaypoints from './ScannerWaypoints.svelte';
	import SelectedMapObject from './SelectedMapObject.svelte';
	import SelectedWaypoint from './SelectedWaypoint.svelte';

	export let game: FullGame;

	const xGetter = (mo: MapObject) => mo?.position?.x;
	const yGetter = (mo: MapObject) => mo?.position?.y;

	let clientWidth = 100;
	let clientHeight = 100;
	let aspectRatio = 1;
	let transform: ZoomTransform;
	let zoomBehavior: ZoomBehavior<HTMLElement, any>;
	let root: HTMLElement;
	let commandedMapObjectIndex = 0;
	let padding = 20; // 20 px, used in zooming
	let scaleX: ScaleLinear<number, number, never>;
	let scaleY: ScaleLinear<number, number, never>;
	let movingWaypoint = false;

	const scale = writable(game.area.y / 400); // tiny games are at 1x starting zoom, the rest zoom in based on universe size
	const clampedScale = writable($scale);
	$: $clampedScale = Math.min(3, $scale); // don't let the scale used for scanner objects go more than 1/2th size
	// $: console.log('scale ', $scale, ' clampedScale', $clampedScale);

	// handle zoom in/out
	// this behavior controls how the zoom behaves
	// below we handle zooming events by updating a transform
	$: {
		if (root) {
			handleResize();

			zoomBehavior = zoom<HTMLElement, any>()
				.extent([
					[0, 0],
					[clientWidth, clientHeight]
				])
				.scaleExtent([0.75, 10])
				.translateExtent([
					[-20, -20],
					[clientWidth + padding, clientHeight + padding]
				])
				.on('zoom', handleZoom);

			enableDragAndZoom();
		}
	}

	$: setPacketDest = $settings.setPacketDest;
	$: addWaypoint = $settings.addWaypoint;

	// enable drag and zoom, but disable dblclick zoom events
	const enableDragAndZoom = () => select(root).call(zoomBehavior).on('dblclick.zoom', null);
	// disable drag and zoom temporarily
	const disableDragAndZoom = () => select(root).on('.zoom', null);

	const xRange = () => {
		if (aspectRatio > 1 && clientHeight > clientWidth) {
			// tall skinny viewport, wide map, so fully expand on the x
			// but shrink up height
			return [0, clientWidth];
		} else if (aspectRatio > 1 && clientWidth > clientHeight) {
			// wide viewport, wide map, so fully expand on the y
			// but shrink up width
			return [0, clientHeight * aspectRatio];
		}
		return [0, Math.min(clientWidth, clientHeight)];
	};
	const yRange = () => {
		if (aspectRatio > 1 && clientHeight > clientWidth) {
			// tall skinny viewport, wide map, so fully expand on the x
			// but shrink up height
			return [0, clientWidth / aspectRatio];
		} else if (aspectRatio > 1 && clientWidth > clientHeight) {
			// wide viewport, wide map, so fully expand on the y
			return [0, clientHeight];
		}
		return [0, Math.min(clientWidth, clientHeight)];
	};

	function handleResize() {
		clientWidth = root?.clientWidth ?? 100;
		clientHeight = root?.clientHeight ?? 100;
		aspectRatio = game.area.x / game.area.y;

		// compute scales
		scaleX = scaleLinear().range(xRange()).domain([0, game.area.x]);
		scaleY = scaleLinear().range(yRange()).domain([0, game.area.y]);
	}

	function handleZoom(e: D3ZoomEvent<HTMLElement, any>) {
		transform = e.transform;
		$scale = transform.k;
		// console.log('handleZoom', e, transform);
	}

	// zoom to the commanded map object every time it changes
	$: if (root && $zoomTarget) {
		translateViewport($zoomTarget.position);
	}

	// zoom the display to a point on the map
	function translateViewport(position: Vector, scaleTo?: number) {
		if (root) {
			select(root).call(zoomBehavior.scaleTo, $scale);
			const scaled: Vector = {
				x: scaleX(position.x),
				y: scaleY(position.y)
			};
			let localScale = $scale;
			if (scaleTo) {
				localScale = scaleTo;
			}
			select(root)
				.call(zoomBehavior.translateTo, scaled.x, scaled.y)
				.call(zoomBehavior.scaleTo, localScale);
		}
	}

	// if the shift key is held, add a waypoint instead of selecting a mapobject
	async function onAddWaypoint(options: { mo?: MapObject; position?: Vector }) {
		if (!$commandedFleet?.waypoints) {
			return;
		}

		let waypointIndex = $commandedFleet.waypoints.findIndex((wp) => $selectedWaypoint == wp);
		if (waypointIndex == -1) {
			waypointIndex = 0;
		}
		const currentWaypoint = $commandedFleet.waypoints[waypointIndex];
		let nextWaypoint = currentWaypoint;
		if (waypointIndex < $commandedFleet.waypoints.length - 1) {
			nextWaypoint = $commandedFleet.waypoints[waypointIndex + 1];
		}

		const position = options.mo ? options.mo.position : options.position ?? emptyVector;
		if (equal(position, currentWaypoint.position) || equal(position, nextWaypoint.position)) {
			// don't duplicate waypoints
			return;
		}

		if (options.mo) {
			const mo = options.mo;
			$commandedFleet.waypoints.splice(waypointIndex + 1, 0, {
				position: mo.position,
				warpSpeed: $commandedFleet.spec?.engine?.idealSpeed ?? 5,
				targetName: mo.name,
				targetPlayerNum: mo.playerNum,
				targetNum: mo.num,
				targetType: mo.type,
				task: WaypointTask.None,
				transportTasks: { fuel: {}, ironium: {}, boranium: {}, germanium: {}, colonists: {} }
			});
		} else if (options.position) {
			$commandedFleet.waypoints.splice(waypointIndex + 1, 0, {
				position: options.position,
				warpSpeed: $commandedFleet.spec?.engine?.idealSpeed ?? 5,
				task: WaypointTask.None,
				transportTasks: { fuel: {}, ironium: {}, boranium: {}, germanium: {}, colonists: {} }
			});
		}

		await game.updateFleetOrders($commandedFleet);
		// select the new waypoint
		selectWaypoint($commandedFleet.waypoints[$commandedFleet.waypoints.length - 1]);
		if ($selectedWaypoint && $selectedWaypoint.targetType && $selectedWaypoint.targetNum) {
			const mo = game.universe.getMapObject($selectedWaypoint);

			if (mo) {
				selectMapObject(mo);
			}
		}
	}
	/**
	 * When a mapobject is selected we go through a few steps.
	 * - We select it if it's a new selection
	 * - We cycle through our commandable objects at the same location if we own an object there
	 * @param mo
	 */
	async function mapobjectSelected(mo: MapObject) {
		// reset waypoint dragging
		enableDragAndZoom();
		movingWaypoint = false;

		let myMapObject = mo;
		if (ownedBy(mo, game.player.num) && mo.type === MapObjectType.Planet) {
			myMapObject = game.getPlanet(mo.num) ?? mo;
		}

		const commandedIntelObject = $commandedMapObject
			? game.universe.getMapObject({
					targetType: $commandedMapObject.type,
					targetNum: $commandedMapObject.num,
					targetPlayerNum: $commandedMapObject.playerNum
			  })
			: undefined;

		if (setPacketDest) {
			if (mo.type != MapObjectType.Planet) {
				return;
			} else {
				$settings.setPacketDest = false;
				// something went wrong, can't set dest on a planet without a massdriver
				if (!$commandedPlanet?.spec.hasMassDriver) {
					return;
				}
				$commandedPlanet.packetTargetNum = mo.num;
				const result = await PlanetService.update(game.id, $commandedPlanet);
				Object.assign($commandedPlanet, result);
				$commandedPlanet = $commandedPlanet;
				game.universe.updatePlanet($commandedPlanet);
				game.universe.planets = game.universe.planets;
			}
		}

		if ($selectedMapObject !== myMapObject) {
			// we selected a different object, so just select it
			selectMapObject(myMapObject);

			// if we selected a mapobject that is a waypoint, select the waypoint as well
			if ($commandedFleet?.waypoints) {
				const fleetWaypoint = $commandedFleet.waypoints.find((wp) =>
					equal(wp.position, myMapObject.position)
				);
				if (fleetWaypoint) {
					selectWaypoint(fleetWaypoint);
				}
			}
		} else {
			if ($commandedFleet?.waypoints) {
				const fleetWaypoint = $commandedFleet.waypoints.find((wp) =>
					equal(wp.position, myMapObject.position)
				);
				if (fleetWaypoint && fleetWaypoint === $selectedWaypoint) {
					movingWaypoint = true;
					disableDragAndZoom();
				}
			}

			// we selected the same mapobject twice
			const myMapObjectsAtPosition = game.universe.getMyMapObjectsByPosition(mo);
			if (myMapObjectsAtPosition?.length > 0) {
				// if our currently commanded map object is not at this location, reset the index
				if (!myMapObjectsAtPosition.find((mo) => mo == commandedIntelObject)) {
					commandedMapObjectIndex = 0;
				} else {
					// command the next one
					commandedMapObjectIndex =
						commandedMapObjectIndex >= myMapObjectsAtPosition.length - 1
							? 0
							: commandedMapObjectIndex + 1;
				}
				const nextMapObject = myMapObjectsAtPosition[commandedMapObjectIndex];
				if (nextMapObject.type === MapObjectType.Planet) {
					commandMapObject(game.getPlanet(nextMapObject.num) as MapObject);
					commandedMapObjectIndex = 0;
				} else {
					commandMapObject(nextMapObject);
				}
			}
		}
	}

	function dragWaypointMove(position: Vector, mo: MapObject | undefined) {
		if (
			$selectedWaypoint &&
			$commandedFleet &&
			movingWaypoint &&
			!equal($selectedWaypoint.position, $commandedFleet?.waypoints[0].position)
		) {
			if (mo) {
				$selectedWaypoint.position = mo.position;
				$selectedWaypoint.targetType = mo.type;
				$selectedWaypoint.targetNum = mo.num;
				$selectedWaypoint.targetPlayerNum = mo.playerNum;
				$selectedWaypoint.targetName = mo.name;
			} else {
				$selectedWaypoint.position = position;
				$selectedWaypoint.targetType = MapObjectType.None;
				$selectedWaypoint.targetNum = None;
				$selectedWaypoint.targetPlayerNum = None;
				$selectedWaypoint.targetName = '';
			}
		}
	}

	function dragWaypointDone(position: Vector, mo: MapObject | undefined) {
		if (
			$selectedWaypoint &&
			$commandedFleet &&
			movingWaypoint &&
			!equal($selectedWaypoint.position, $commandedFleet?.waypoints[0].position)
		) {
			game.updateFleetOrders($commandedFleet);
		}
	}

	let data: MapObject[] = [];
	$: {
		const waypoints: MapObject[] = [];
		if ($commandedFleet?.waypoints) {
			waypoints.push(
				...$commandedFleet.waypoints.map((wp) => {
					const mo = game.universe.getMapObject(wp);
					if (mo) {
						return mo;
					} else {
						return {
							position: wp.position,
							type: wp.targetType ?? MapObjectType.PositionWaypoint,
							name: wp.targetName ?? '',
							num: wp.targetNum ?? 0,
							playerNum: wp.targetPlayerNum ?? 0
						} as MapObject;
					}
				})
			);
		}
		data = [
			...waypoints,
			...game.universe.fleets.filter(
				(f) => f.orbitingPlanetNum === None || f.orbitingPlanetNum === undefined
			),
			...(game.player.fleetIntels?.filter((f) => !f.orbitingPlanetNum) ?? []),
			...game.universe.mineralPackets,
			...game.player.mineralPacketIntels,
			...game.universe.mineFields,
			...game.player.mineFieldIntels,
			...game.player.planetIntels
		];
	}

	setContext('game', game);
	setContext('scale', clampedScale);
</script>

<svelte:window on:resize={handleResize} />

<div class={`grow bg-black overflow-hidden p-[${padding}px]`}>
	<LayerCake
		{data}
		x={xGetter}
		y={yGetter}
		xDomain={[0, game.area.x]}
		yDomain={[0, game.area.y]}
		{xRange}
		{yRange}
		yReverse={true}
		bind:element={root}
	>
		<!-- <Svg viewBox={`0 0 ${game.area.x} ${game.area.y}`}> -->
		<Svg>
			<g preserveAspectRatio="true" transform={transform?.toString()}>
				<ScannerScanners />
				<ScannerMineFieldPattern />
				<ScannerMineFields />
				<ScannerPacketDests />
				<ScannerWaypoints />
				<SelectedWaypoint />
				<ScannerPlanets />
				<ScannerMineralPackets />
				<ScannerFleets />
				<ScannerWarpLine />
				<SelectedMapObject />
			</g>
		</Svg>
		<Html>
			{#if transform}
				<ScannerNames {transform} />

				<MapObjectQuadTreeFinder
					on:mapobject-selected={(mo) => mapobjectSelected(mo.detail)}
					on:add-waypoint={(mo) => onAddWaypoint(mo.detail)}
					on:drag-waypoint-move={(e) => dragWaypointMove(e.detail.position, e.detail.mo)}
					on:drag-waypoint-done={(e) => dragWaypointDone(e.detail.position, e.detail.mo)}
					searchRadius={20}
					let:x
					let:y
					let:found
					{transform}
				/>
			{/if}
		</Html>
	</LayerCake>
</div>
