import { AdminService } from '#lib/types/cs-proto.js';
import { FleetService } from '#lib/types/cs-proto.js';
import { GameService } from '#lib/types/cs-proto.js';
import { MinefieldService } from '#lib/types/cs-proto.js';
import { PlanetService } from '#lib/types/cs-proto.js';
import {
	BattlePlanService,
	PlayerService,
	ProductionPlanService,
	TransportPlanService
} from '#lib/types/cs-proto.js';
import { RaceService } from '#lib/types/cs-proto.js';
import { ShipDesignService } from '#lib/types/cs-proto.js';
import { TechService } from '#lib/types/cs-proto.js';
import { UserService } from '#lib/types/cs-proto.js';
import { createClient } from '@connectrpc/connect';
import { createConnectTransport } from '@connectrpc/connect-web';
import { asPlayerInterceptor } from './asPlayerInterceptor';
import { retryInterceptor } from './retryInterceptor';
import { revisionInterceptor } from './revisionInterceptor';
import { versionInterceptor } from './versionInterceptor';

// hit our grpc endpoint
const transport = createConnectTransport({
	baseUrl: '/api/grpc',
	// the revisionInterceptor is last so it sees the X-As-Player header
	interceptors: [retryInterceptor, versionInterceptor, asPlayerInterceptor, revisionInterceptor]
	// useBinaryFormat: true
});

export const adminClient = createClient(AdminService, transport);
export const battlePlanClient = createClient(BattlePlanService, transport);
export const fleetClient = createClient(FleetService, transport);
export const gameClient = createClient(GameService, transport);
export const minefieldClient = createClient(MinefieldService, transport);
export const planetClient = createClient(PlanetService, transport);
export const playerClient = createClient(PlayerService, transport);
export const productionPlanClient = createClient(ProductionPlanService, transport);
export const raceClient = createClient(RaceService, transport);
export const shipDesignClient = createClient(ShipDesignService, transport);
export const techClient = createClient(TechService, transport);
export const transportPlanClient = createClient(TransportPlanService, transport);
export const userClient = createClient(UserService, transport);
