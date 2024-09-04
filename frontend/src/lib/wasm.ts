import type { Race } from '$lib/types/Race';
import type { Fleet } from './types/Fleet';
import { type Planet } from './types/Planet';
import type { Player } from './types/Player';
import type { Rules } from './types/Rules';
import type { ShipDesign } from './types/ShipDesign';

export type CS = {
	enableDebug: () => void;
	setRules: (rules: Rules) => void;
	setPlayer: (player: Player) => void;
	setDesigns: (designs: ShipDesign[]) => void;
	calculateRacePoints: (race: Race) => number;
	estimateProduction: (planet: Planet) => Promise<Planet>;
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

// our wasm calls actually take json strings as params for easier serializing between go/typescript
// this type represents the actual cs.wasm calls
type CSWasm = {
	calculateRacePoints: (race: Race) => number;
	enableDebug: () => void;
	setRules: (rules: Rules) => void;
	setPlayer: (player: Player) => void;
	setDesigns: (designs: ShipDesign[]) => void;
	estimateProduction: (planet: Planet) => Promise<Planet>;
};

// create a wrapper to serialize requests and responses to/from JSON
class CSWasmWrapper implements CS {
	constructor(private wasm: CSWasm) {}

	async enableDebug() {
		this.wasm.enableDebug();
	}

	setRules(rules: Rules) {
		this.wasm.setRules(rules);
	}

	setPlayer(player: Player) {
		this.wasm.setPlayer(player);
	}

	setDesigns(designs: ShipDesign[]) {
		this.wasm.setDesigns(designs);
	}

	async estimateProduction(planet: Planet): Promise<Planet> {
		return await this.wasm.estimateProduction(planet);
	}
	calculateRacePoints(race: Race): number {
		return this.wasm.calculateRacePoints(race);
	}
}
