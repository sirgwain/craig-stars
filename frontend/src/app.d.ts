/// <reference types="@sveltejs/kit" />
/// <reference types="golang-wasm-exec" />

// See https://kit.svelte.dev/docs/types#app
// for information about these interfaces
declare namespace App {
	// interface Locals {}
	// interface Platform {}
	// interface Session {}
	// interface Stuff {}
}

declare const PKG: {
	version: string;
};
