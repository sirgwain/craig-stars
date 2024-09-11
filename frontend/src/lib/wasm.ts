import type { Race } from '$lib/types/Race';
import { addError } from './services/Errors';
import { type Planet } from './types/Planet';
import type { Player } from './types/Player';
import type { Rules } from './types/Rules';
import type { ShipDesign, Spec as ShipDesignSpec } from './types/ShipDesign';

export type CS = {
	enableDebug: () => void;
	setRules: (rules: Rules) => void;
	setPlayer: (player: Player) => void;
	setDesigns: (designs: ShipDesign[]) => void;
	calculateRacePoints: (race: Race) => number | undefined;
	computeShipDesignSpec: (design: ShipDesign) => ShipDesignSpec | undefined;
	estimateProduction: (planet: Planet) => Planet | undefined;
};

// load a wasm module and returns a wrapper for executing functions
export async function loadWasm(): Promise<CS> {
	// @ts-expect-error
	if (typeof __go_wasm__ == 'undefined') {
		// @ts-expect-error
		window.__go_wasm__ = {};
	}

	// after loading the wasm module, __go_wasm__ will be replaced with our module
	// it will contain our cs methods along with a ready boolean
	type Bridge = {
		__ready__?: boolean;
	};

	// @ts-expect-error
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

	setPlayer(player: Player) {
		this.wasm.setPlayer(player);
		this.checkError();
	}

	setDesigns(designs: ShipDesign[]) {
		this.wasm.setDesigns(designs);
		this.checkError();
	}

	estimateProduction(planet: Planet): Planet | undefined {
		const result = this.wasm.estimateProduction(planet);
		if (this.checkError()) {
			return undefined;
		}
		return result;
	}

	computeShipDesignSpec(design: ShipDesign): ShipDesignSpec | undefined {
		const result = this.wasm.computeShipDesignSpec(design);
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
}
