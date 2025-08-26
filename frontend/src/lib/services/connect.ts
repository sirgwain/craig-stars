import { AdminService } from '$lib/types/cs-proto';
import { BattleService } from '$lib/types/cs-proto';
import { FleetService } from '$lib/types/cs-proto';
import { GameService } from '$lib/types/cs-proto';
import { MinefieldService } from '$lib/types/cs-proto';
import { PlanetService } from '$lib/types/cs-proto';
import {
	BattlePlanService,
	PlayerService,
	ProductionPlanService,
	TransportPlanService
} from '$lib/types/cs-proto';
import { RaceService } from '$lib/types/cs-proto';
import { ShipDesignService } from '$lib/types/cs-proto';
import { TechService } from '$lib/types/cs-proto';
import { UserService } from '$lib/types/cs-proto';
import { createClient } from '@connectrpc/connect';
import { createConnectTransport } from '@connectrpc/connect-web';

// hit our grpc endpoint
const transport = createConnectTransport({
	baseUrl: '/api/grpc'
	// useBinaryFormat: true
});

export const adminClient = createClient(AdminService, transport);
export const battleClient = createClient(BattleService, transport);
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
