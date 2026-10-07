import { get } from 'svelte/store';
import { describe, expect, it } from 'vitest';
import { createVersionStore } from './Version';

describe('release notifications', () => {
	it('ignores missing and matching versions and all signals during development', () => {
		const store = createVersionStore('v1.42.0');
		store.observeServerVersion(null);
		store.observeServerVersion('');
		store.observeServerVersion('v1.42.0');
		expect(get(store).available).toBe(false);

		const develop = createVersionStore('0.0.0-develop');
		develop.observeServerVersion('v1.43.0');
		develop.observeFrontendUpdate();
		expect(get(develop).available).toBe(false);
	});

	it('keeps Later dismissed across polling and API responses, then prompts for another release', () => {
		const store = createVersionStore('v1.42.0');
		store.observeServerVersion('v1.43.0');
		expect(get(store)).toMatchObject({ available: true, dismissed: false });
		store.dismiss();
		store.observeFrontendUpdate();
		store.observeServerVersion('v1.43.0');
		expect(get(store).dismissed).toBe(true);
		store.observeServerVersion('v1.44.0');
		expect(get(store)).toMatchObject({ serverVersion: 'v1.44.0', dismissed: false });
	});

	it('adds the version label without reopening a dismissed idle-tab prompt', () => {
		const store = createVersionStore('v1.42.0');
		store.observeFrontendUpdate();
		expect(get(store)).toMatchObject({ available: true, serverVersion: null });
		store.dismiss();
		store.observeServerVersion('v1.43.0');
		expect(get(store)).toMatchObject({ serverVersion: 'v1.43.0', dismissed: true });
	});
});
