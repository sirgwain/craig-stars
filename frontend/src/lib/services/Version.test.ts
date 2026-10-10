import { get } from 'svelte/store';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { createVersionStore, watchRelease } from './Version';

describe('release notifications', () => {
	it('ignores missing and matching versions and all signals during development', () => {
		const store = createVersionStore('v1.42.0');
		store.observeFrontendUpdate();
		store.observeServerVersion(null);
		store.observeServerVersion('');
		store.observeServerVersion('v1.42.0');
		expect(get(store).available).toBe(false);

		const develop = createVersionStore('0.0.0-develop');
		develop.observeServerVersion('v1.43.0');
		develop.observeFrontendUpdate();
		expect(get(develop).available).toBe(false);
	});

	it('waits for both the server and the frontend, in either order', () => {
		const serverFirst = createVersionStore('v1.42.0');
		serverFirst.observeServerVersion('v1.43.0');
		expect(get(serverFirst).available).toBe(false);
		serverFirst.observeFrontendUpdate();
		expect(get(serverFirst)).toMatchObject({ available: true, serverVersion: 'v1.43.0' });

		const frontendFirst = createVersionStore('v1.42.0');
		frontendFirst.observeFrontendUpdate();
		expect(get(frontendFirst).available).toBe(false);
		frontendFirst.observeServerVersion('v1.43.0');
		expect(get(frontendFirst).available).toBe(true);
	});

	it('stops waiting when a frontend loaded mid-deploy sees the server catch up', () => {
		const store = createVersionStore('v1.43.0');
		store.observeServerVersion('v1.42.0');
		expect(get(store).serverVersion).toBe('v1.42.0');
		store.observeServerVersion('v1.43.0');
		expect(get(store)).toMatchObject({ available: false, serverVersion: null });
	});

	it('keeps Later dismissed across polling and API responses, then prompts for another release', () => {
		const store = createVersionStore('v1.42.0');
		store.observeServerVersion('v1.43.0');
		store.observeFrontendUpdate();
		expect(get(store)).toMatchObject({ available: true, dismissed: false });
		store.dismiss();
		store.observeFrontendUpdate();
		store.observeServerVersion('v1.43.0');
		expect(get(store).dismissed).toBe(true);
		store.observeServerVersion('v1.44.0');
		expect(get(store)).toMatchObject({ serverVersion: 'v1.44.0', dismissed: false });
	});
});

describe('watching for the rest of a release', () => {
	beforeEach(() => vi.useFakeTimers());
	afterEach(() => vi.useRealTimers());

	it('checks for the missing half until it is released, backing off', async () => {
		const store = createVersionStore('v1.42.0');
		const frontend = vi.fn(async () => {});
		const server = vi.fn(async () => {
			throw new Error('502');
		});
		const stop = watchRelease(store, { frontend, server }, [100, 200]);

		// nothing has been released
		await vi.advanceTimersByTimeAsync(1000);
		expect(frontend).not.toHaveBeenCalled();
		expect(server).not.toHaveBeenCalled();

		// the frontend is out, but the server is still restarting
		store.observeFrontendUpdate();
		await vi.advanceTimersByTimeAsync(0);
		expect(server).toHaveBeenCalledTimes(1);
		await vi.advanceTimersByTimeAsync(100);
		expect(server).toHaveBeenCalledTimes(2);
		await vi.advanceTimersByTimeAsync(100);
		expect(server).toHaveBeenCalledTimes(2);
		await vi.advanceTimersByTimeAsync(100);
		expect(server).toHaveBeenCalledTimes(3);

		// once the server is up we stop checking
		store.observeServerVersion('v1.43.0');
		await vi.advanceTimersByTimeAsync(1000);
		expect(server).toHaveBeenCalledTimes(3);
		expect(frontend).not.toHaveBeenCalled();
		stop();
	});

	it('checks for the frontend when the server is released first', async () => {
		const store = createVersionStore('v1.42.0');
		const frontend = vi.fn(async () => store.observeFrontendUpdate());
		const server = vi.fn(async () => {});
		const stop = watchRelease(store, { frontend, server }, [100]);

		store.observeServerVersion('v1.43.0');
		await vi.advanceTimersByTimeAsync(1000);
		expect(frontend).toHaveBeenCalledTimes(1);
		expect(server).not.toHaveBeenCalled();
		expect(get(store).available).toBe(true);
		stop();
	});
});
