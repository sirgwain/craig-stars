import { Code, createClient } from '@connectrpc/connect';
import { createConnectTransport } from '@connectrpc/connect-web';
import { expect, it, vi } from 'vitest';
import { PlayerService } from '#lib/types/cs-proto.js';
import { createRetryInterceptor } from './retryInterceptor';

const badGateway = () => new Response('<html>502 Bad Gateway</html>', { status: 502 });
const ok = () => new Response('{}', { headers: { 'Content-Type': 'application/json' } });

function clientFor(fetch: typeof globalThis.fetch) {
	return createClient(
		PlayerService,
		createConnectTransport({
			baseUrl: 'https://example.test/api/grpc',
			interceptors: [createRetryInterceptor([0, 0])],
			fetch
		})
	);
}

it('retries idempotent calls through a restart', async () => {
	const fetch = vi.fn().mockResolvedValueOnce(badGateway()).mockResolvedValueOnce(ok());
	await clientFor(fetch).getPlayer({});
	expect(fetch).toHaveBeenCalledTimes(2);
});

it('never retries calls without an idempotency level', async () => {
	const fetch = vi.fn().mockImplementation(async () => badGateway());
	await expect(clientFor(fetch).submitTurn({})).rejects.toMatchObject({ code: Code.Unavailable });
	expect(fetch).toHaveBeenCalledTimes(1);
});
