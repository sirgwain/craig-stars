import { writable, type Readable } from 'svelte/store';

type VersionUpdate = {
	// true once both the frontend and the server have been released
	available: boolean;
	serverVersion: string | null;
	frontendUpdated: boolean;
	dismissed: boolean;
};

export function createVersionStore(clientVersion: string) {
	const { subscribe, update } = writable<VersionUpdate>({
		available: false,
		serverVersion: null,
		frontendUpdated: false,
		dismissed: false
	});
	const enabled = clientVersion !== '0.0.0-develop';

	// A deploy replaces the server and the frontend one after the other. Reloading before both are
	// done would run mismatched versions, so wait to tell the player until we've seen both.
	const withAvailable = (state: VersionUpdate): VersionUpdate => ({
		...state,
		available: state.frontendUpdated && state.serverVersion !== null
	});

	return {
		subscribe,
		observeServerVersion(serverVersion: string | null) {
			if (!enabled || !serverVersion) return;
			if (serverVersion === clientVersion) {
				// we loaded a new frontend mid-deploy and the server has caught up to us
				update((state) =>
					state.serverVersion === null ? state : withAvailable({ ...state, serverVersion: null })
				);
				return;
			}
			update((state) => {
				if (state.serverVersion === serverVersion) return state;
				return withAvailable({
					...state,
					serverVersion,
					// prompt again if another release goes out after this one was dismissed
					dismissed: state.serverVersion === null ? state.dismissed : false
				});
			});
		},
		observeFrontendUpdate() {
			if (!enabled) return;
			update((state) =>
				state.frontendUpdated ? state : withAvailable({ ...state, frontendUpdated: true })
			);
		},
		dismiss() {
			update((state) => ({ ...state, dismissed: true }));
		}
	};
}

export const versionUpdate = createVersionStore(PKG.version);

type ReleaseChecks = {
	// ask for the frontend's version, i.e. SvelteKit's updated.check()
	frontend: () => Promise<unknown>;
	// make any api call so the server's version header is observed
	server: () => Promise<unknown>;
};

// After seeing half of a release, check for the other half until it shows up so the player is
// told even if they aren't doing anything. Returns a function to stop watching.
export function watchRelease(
	store: Readable<VersionUpdate>,
	checks: ReleaseChecks,
	delaysMs = [15_000, 30_000, 60_000, 5 * 60_000]
): () => void {
	let timer: ReturnType<typeof setTimeout> | undefined;
	let attempt = 0;
	let waitingFor: keyof ReleaseChecks | undefined;

	const check = async () => {
		timer = undefined;
		if (!waitingFor) return;
		try {
			await checks[waitingFor]();
		} catch {
			// the server is probably restarting, we'll check again
		}
		// the check may have completed the release
		if (waitingFor && !timer) {
			timer = setTimeout(check, delaysMs[Math.min(attempt++, delaysMs.length - 1)]);
		}
	};

	const unsubscribe = store.subscribe((state) => {
		const seenAny = state.frontendUpdated || state.serverVersion !== null;
		waitingFor =
			!seenAny || state.available ? undefined : state.frontendUpdated ? 'server' : 'frontend';
		if (!waitingFor) {
			clearTimeout(timer);
			timer = undefined;
			attempt = 0;
		} else if (!timer) {
			timer = setTimeout(check, 0);
		}
	});

	return () => {
		unsubscribe();
		waitingFor = undefined;
		clearTimeout(timer);
	};
}
