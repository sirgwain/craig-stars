<script lang="ts">
	import type { Planet } from '$lib/types/cs-proto';
	import { Grav, Rad, Temp } from '$lib/types/Hab';
	import type { CommandedPlayer } from '$lib/types/Player';
	import PlanetHabBar from './PlanetHabBar.svelte';

	type Props = {
		planet: Planet;
		player: CommandedPlayer;
		onGravityTooltip?: (e: PointerEvent) => void;
		onTemperatureTooltip?: (e: PointerEvent) => void;
		onRadiationTooltip?: (e: PointerEvent) => void;
	};

	let { planet, player, onGravityTooltip, onTemperatureTooltip, onRadiationTooltip }: Props =
		$props();
</script>

<PlanetHabBar
	habType={Grav}
	value={planet.hab?.grav ?? 0}
	baseValue={planet.baseHab?.grav ?? 0}
	terraformValue={planet.spec?.terraformAmount?.grav ?? 0}
	high={player.race.habHigh?.grav ?? 0}
	low={player.race.habLow?.grav ?? 0}
	immune={!!player.race.immuneGrav}
	onTooltip={onGravityTooltip}
/>
<PlanetHabBar
	habType={Temp}
	value={planet.hab?.temp ?? 0}
	baseValue={planet.baseHab?.temp ?? 0}
	terraformValue={planet.spec?.terraformAmount?.temp ?? 0}
	high={player.race.habHigh?.temp ?? 0}
	low={player.race.habLow?.temp ?? 0}
	immune={!!player.race.immuneTemp}
	onTooltip={onTemperatureTooltip}
/>
<PlanetHabBar
	habType={Rad}
	value={planet.hab?.rad ?? 0}
	baseValue={planet.baseHab?.rad ?? 0}
	terraformValue={planet.spec?.terraformAmount?.rad ?? 0}
	high={player.race.habHigh?.rad ?? 0}
	low={player.race.habLow?.rad ?? 0}
	immune={!!player.race.immuneRad}
	onTooltip={onRadiationTooltip}
/>
