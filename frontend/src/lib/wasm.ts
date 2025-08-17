import { createClient, type Client, type Transport } from '@connectrpc/connect';
import { createConnectTransport } from '@connectrpc/connect-web';
import { WasmService } from './protogen/craig_stars/v1/wasm_pb';

export type WasmClient = Client<typeof WasmService>;
export type CS = {
	transport: Transport;
	wasmService: WasmClient;
	handleGrpc: (
		path: string,
		requestBytes: Uint8Array,
		callback: (err: unknown, responseBytes: Uint8Array) => void
	) => void;
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
		await cs.wasmService.enableDebug({ debug: true });
	}

	return cs;
}

// create a wrapper to serialize requests and responses to/from JSON
class CSWasmWrapper implements CS {
	constructor(private wasm: CS) {}

	// This wasm transport used for grpc wasm calls
	transport = createConnectTransport({
		baseUrl: '',
		useBinaryFormat: true,
		fetch: async (input, init) => {
			const reqBytes = new Uint8Array(await new Request(input, init).arrayBuffer());
			const resBytes = await new Promise<Uint8Array>((resolve, reject) => {
				this.wasm.handleGrpc(input.toString(), reqBytes, (err: unknown, res: Uint8Array) => {
					if (err) reject(err);
					else resolve(res);
				});
			});

			return new Response(resBytes.slice().buffer, {
				status: 200,
				headers: { 'Content-Type': 'application/proto' }
			});
		}
	});

	wasmService = createClient(WasmService, this.transport);
	async handleGrpc(
		path: string,
		requestBytes: Uint8Array,
		callback: (err: unknown, responseBytes: Uint8Array) => void
	) {
		this.wasm.handleGrpc(path, requestBytes, callback);
	}
}
