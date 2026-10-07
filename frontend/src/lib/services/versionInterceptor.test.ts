import { Code, createClient } from '@connectrpc/connect';
import { createConnectTransport } from '@connectrpc/connect-web';
import { expect, it, vi } from 'vitest';
import { TechService } from '#lib/types/cs-proto.js';
import { versionInterceptor } from './versionInterceptor';

const { observeServerVersion } = vi.hoisted(() => ({ observeServerVersion: vi.fn() }));
vi.mock('./Version', () => ({ versionUpdate: { observeServerVersion } }));

it('reads the version from error metadata and rethrows the failure', async () => {
	const client = createClient(
		TechService,
		createConnectTransport({
			baseUrl: 'https://example.test/api/grpc',
			interceptors: [versionInterceptor],
			fetch: async () =>
				new Response('{"code":"unauthenticated","message":"sign in"}', {
					status: 401,
					headers: { 'Content-Type': 'application/json', 'X-App-Version': 'v1.43.0' }
				})
		})
	);
	await expect(client.getTechs({})).rejects.toMatchObject({ code: Code.Unauthenticated });
	expect(observeServerVersion).toHaveBeenCalledWith('v1.43.0');
});
