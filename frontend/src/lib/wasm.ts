import { addError } from './services/Errors';
import type {
	Cost,
	Fleet,
	PlayerIntels,
	QueueItemType,
	Race,
	Rules,
	ShipDesign,
	ShipDesignSpec,
	Tech,
	TechLevel,
	WaypointDest
} from './types/cs';
import { type Planet } from './types/cs';
import type { CommandedPlayer } from './types/Player';

export type CS = {
	enableDebug: () => void;
	setRules: (rules: Rules) => void;
	setPlayer: (player: CommandedPlayer) => void;
	setDesigns: (designs: ShipDesign[]) => void;
	setIntel: (intel: PlayerIntels) => void;
	calculateRacePoints: (race: Race) => number | undefined;
	getResearchCost: (techLevel: TechLevel) => number | undefined;
	computeShipDesignSpec: (design: ShipDesign) => ShipDesignSpec | undefined;
	starbaseUpgradeCost: (design: ShipDesign, newDesign: ShipDesign) => Cost | undefined;
	techCost: (tech: Tech) => Cost | undefined;
	estimateProduction: (planet: Planet) => Planet | undefined;
	maxBuildable: (planet: Planet, itemType: QueueItemType) => number | undefined;
	addWaypoint: (
		fleet: Fleet,
		dest: WaypointDest,
		currentSelectedWaypointIndex: number,
		fastestWaypoint: boolean
	) => { fleet: Fleet; result: number } | undefined;
	updateWaypoint: (
		fleet: Fleet,
		dest: WaypointDest,
		currentSelectedWaypointIndex: number,
		fastestWaypoint: boolean
	) => { fleet: Fleet; result: boolean } | undefined;
	updateResourcesAvailable: (planet: Planet) => number | undefined;
};

// load a wasm module and returns a wrapper for executing functions
export async function loadWasm(): Promise<CS> {
	// @ts-expect-error __go_wasm__ is set by cs.wasm
	if (typeof __go_wasm__ == 'undefined') {
		// @ts-expect-error __go_wasm__ is set by cs.wasm
		window.__go_wasm__ = {};
	}

	// after loading the wasm module, __go_wasm__ will be replaced with our module
	// it will contain our cs methods along with a ready boolean
	type Bridge = {
		__ready__?: boolean;
	};

	// @ts-expect-error __go_wasm__ is set by cs.wasm
	const bridge = __go_wasm__ as CSWasm & Bridge;

	// load the wasm and start it up
	const csWasmUrl = new URL('$lib/wasm/cs.wasm', import.meta.url).href;
	const go = new Go();
	const result = await WebAssembly.instantiateStreaming(fetch(csWasmUrl), go.importObject);
	go.run(result.instance);

	// wait until the wasm finishes initializing
	let readyCount = 0;
	while (bridge.__ready__ !== true) {
		if (readyCount > 100) {
			throw Error('wasm was never ready');
		}
		await new Promise<void>((res) => {
			requestAnimationFrame(() => res());
			setTimeout(() => {
				readyCount++;
				res();
			}, 50);
		});
	}

	// all done, ready to execute!
	const cs = new CSWasmWrapper(bridge);

	if (PKG.version == '0.0.0-develop') {
		cs.enableDebug();
	}

	return cs;
}

// create a wrapper to serialize requests and responses to/from JSON
class CSWasmWrapper implements CS {
	constructor(private wasm: CS) {}

	// checkError checks if the wasm code threw an error and if so adds it as a notification
	// and return true
	checkError(): boolean {
		if ('wasmError' in window) {
			addError(`${window['wasmError']}`);
			delete window.wasmError;
			return true;
		}
		return false;
	}

	async enableDebug() {
		this.wasm.enableDebug();
		this.checkError();
	}

	setRules(rules: Rules) {
		this.wasm.setRules(rules);
		this.checkError();
	}

	setPlayer(player: CommandedPlayer) {
		this.wasm.setPlayer(player);
		this.checkError();
	}

	setDesigns(designs: ShipDesign[]) {
		this.wasm.setDesigns(designs);
		this.checkError();
	}

	setIntel(intel: PlayerIntels) {
		this.wasm.setIntel(intel);
		this.checkError();
	}

	computeShipDesignSpec(design: ShipDesign): ShipDesignSpec | undefined {
		const result = this.wasm.computeShipDesignSpec(design);
		if (this.checkError()) {
			return undefined;
		}
		return result;
	}

	starbaseUpgradeCost(design: ShipDesign, newDesign: ShipDesign): Cost | undefined {
		const result = this.wasm.starbaseUpgradeCost(design, newDesign);
		if (this.checkError()) {
			return undefined;
		}
		return result;
	}

	techCost(tech: Tech): Cost | undefined {
		const result = this.wasm.techCost(tech);
		if (this.checkError()) {
			return undefined;
		}
		return result;
	}

	calculateRacePoints(race: Race): number | undefined {
		const result = this.wasm.calculateRacePoints(race);
		if (this.checkError()) {
			return undefined;
		}
		return result;
	}

	estimateProduction(planet: Planet): Planet | undefined {
		const result = this.wasm.estimateProduction(planet);
		if (this.checkError()) {
			return undefined;
		}
		return result;
	}

	getResearchCost(techLevel: TechLevel): number | undefined {
		const result = this.wasm.getResearchCost(techLevel);
		if (this.checkError()) {
			return undefined;
		}
		return result;
	}

	maxBuildable(planet: Planet, itemType: QueueItemType): number | undefined {
		const result = this.wasm.maxBuildable(planet, itemType);
		if (this.checkError()) {
			return undefined;
		}
		return result;
	}

	addWaypoint(
		fleet: Fleet,
		dest: WaypointDest,
		currentSelectedWaypointIndex: number,
		fastestWaypoint: boolean
	): { fleet: Fleet; result: number } | undefined {
		const result = this.wasm.addWaypoint(
			fleet,
			dest,
			currentSelectedWaypointIndex,
			fastestWaypoint
		);
		if (this.checkError()) {
			return undefined;
		}
		return result;
	}

	updateWaypoint(
		fleet: Fleet,
		dest: WaypointDest,
		currentSelectedWaypointIndex: number,
		fastestWaypoint: boolean
	): { fleet: Fleet; result: boolean } | undefined {
		const result = this.wasm.updateWaypoint(
			fleet,
			dest,
			currentSelectedWaypointIndex,
			fastestWaypoint
		);
		if (this.checkError()) {
			return undefined;
		}
		return result;
	}

	updateResourcesAvailable(planet: Planet): number | undefined {
		const result = this.wasm.updateResourcesAvailable(planet);
		if (this.checkError()) {
			return undefined;
		}
		return result;
	}
}
