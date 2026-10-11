import { Code, createClient } from '@connectrpc/connect';
import { createConnectTransport } from '@connectrpc/connect-web';
import { get } from 'svelte/store';
import { beforeEach, expect, it } from 'vitest';
import { FleetService, GameService } from '#lib/types/cs-proto.js';
import { resetPlayerRevisions, revisionInterceptor } from './revisionInterceptor';
import { staleData } from './Stale';

// a fake server tracking a revision like the real one
let serverRevision = 0;
let sent: (string | null)[] = [];

const transport = createConnectTransport({
	baseUrl: 'https://example.test/api/grpc',
	interceptors: [revisionInterceptor],
	fetch: async (input, init) => {
		const url = `${input}`;
		const headers = new Headers(init?.headers);
		const clientRevision = headers.get('X-Player-Revision');
		sent.push(clientRevision);
		// let other requests start before this one responds
		await new Promise((resolve) => setTimeout(resolve, 1));
		const respond = (body: string, status = 200) =>
			new Response(body, {
				status,
				headers: { 'Content-Type': 'application/json', 'X-Player-Revision': `${serverRevision}` }
			});
		if (url.endsWith('GetGame')) {
			return respond('{}');
		}
		if (clientRevision !== null && clientRevision !== `${serverRevision}`) {
			return respond('{"code":"aborted","message":"newer changes"}', 409);
		}
		serverRevision++;
		return respond('{}');
	}
});
const games = createClient(GameService, transport);
const fleets = createClient(FleetService, transport);

beforeEach(() => {
	serverRevision = 0;
	sent = [];
	resetPlayerRevisions();
});

it('sends changes one at a time with the latest revision', async () => {
	await games.getGame({ gameId: 1n });
	await Promise.all([
		fleets.renameFleet({ gameId: 1n, fleetNum: 1, name: 'a' }),
		fleets.renameFleet({ gameId: 1n, fleetNum: 1, name: 'b' }),
		games.getGame({ gameId: 1n })
	]);
	await games.getGame({ gameId: 1n });
	expect(sent.filter((revision) => revision !== null)).toEqual(['0', '1']);
	expect(get(staleData)).toBe(false);
});

it('notices changes from another device on the next read and rejects further changes', async () => {
	await games.getGame({ gameId: 1n });
	serverRevision++;
	// requests for other games don't matter
	await games.getGame({ gameId: 2n });
	expect(get(staleData)).toBe(false);

	await games.getGame({ gameId: 1n });
	expect(get(staleData)).toBe(true);
	await expect(fleets.renameFleet({ gameId: 1n, fleetNum: 1, name: 'a' })).rejects.toMatchObject({
		code: Code.Aborted
	});

	// reloading the game starts from the server's revision
	resetPlayerRevisions();
	await games.getGame({ gameId: 1n });
	await fleets.renameFleet({ gameId: 1n, fleetNum: 1, name: 'a' });
	expect(get(staleData)).toBe(false);
});
