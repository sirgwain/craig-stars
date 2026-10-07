import { writable } from 'svelte/store';

type VersionUpdate = {
	available: boolean;
	serverVersion: string | null;
	dismissed: boolean;
};

export function createVersionStore(clientVersion: string) {
	const { subscribe, update } = writable<VersionUpdate>({
		available: false,
		serverVersion: null,
		dismissed: false
	});
	const enabled = clientVersion !== '0.0.0-develop';

	return {
		subscribe,
		observeServerVersion(serverVersion: string | null) {
			if (!enabled || !serverVersion || serverVersion === clientVersion) return;
			update((state) => {
				if (state.serverVersion === serverVersion) return state;
				return {
					available: true,
					serverVersion,
					// A first API signal can add a label to an already-dismissed frontend update.
					dismissed: state.serverVersion === null ? state.dismissed : false
				};
			});
		},
		observeFrontendUpdate() {
			if (!enabled) return;
			update((state) => (state.available ? state : { ...state, available: true }));
		},
		dismiss() {
			update((state) => ({ ...state, dismissed: true }));
		}
	};
}

export const versionUpdate = createVersionStore(PKG.version);
