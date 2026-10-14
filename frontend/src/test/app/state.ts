// Stub for $app/state in browser component tests, which run without SvelteKit.
export const page = {
	url: new URL(globalThis.location?.href ?? 'http://localhost/'),
	params: {} as Record<string, string>,
	route: { id: null as string | null },
	status: 200,
	error: null,
	data: {} as Record<string, unknown>,
	form: undefined,
	state: {}
};

export const navigating = { current: null };

export const updated = {
	current: false,
	check: async () => false
};
