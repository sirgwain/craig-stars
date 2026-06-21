import { expect, submitTurn, test } from './setup';

test('shows homeworld minerals and status', async ({ testGamePage }) => {
	const { gamePage } = await testGamePage('Kitchen Sink');

	await expect(gamePage.tile('Minerals on Hand')).toBeVisible();
	await gamePage.expectTileCargo('Minerals on Hand', {
		ironium: '1000kT',
		boranium: '1000kT',
		germanium: '1000kT'
	});
	await gamePage.expectTileText('Status', ['Population 250,000', 'Defense Type SDI']);
});

test('selects minefield from scanner context menu', async ({ testGamePage }) => {
	const { gamePage, universe } = await testGamePage('Kitchen Sink');

	await gamePage.rightClickMapObject(universe.planets[0]);
	await gamePage.clickScannerContextButton('Humanoids Standard Minefield #1');
	await gamePage.expectSummaryText([
		'Humanoids Standard Minefield #1',
		'Location: (0, 0)',
		'Field Type: Standard',
		'Field Radius: 32 l.y. (1000 mines)',
		'Maximum Safe Speed: Warp 4'
	]);
});

test('shows mineral packet summary', async ({ testGamePage }) => {
	const { gamePage, universe } = await testGamePage('Kitchen Sink');

	await gamePage.selectMapObject(universe.mineralPackets[0]);
	await gamePage.expectSummaryText([
		'Humanoids Mineral Packet #1',
		'Location: (50, 0)',
		'Traveling at Warp: 5',
		'Destination: Planet 1',
		'ETA: 2 years',
		'Ironium 50kT',
		'Boranium 50kT',
		'Germanium 50kT'
	]);
});

test('shows foreign fleet and planet summaries', async ({ testGamePage }) => {
	const { gamePage, universe } = await testGamePage('Kitchen Sink');

	const otherPlayerFleet = universe.fleets.find((fleet) => fleet.mapObject?.playerNum === 2);
	if (!otherPlayerFleet?.mapObject) {
		throw new Error("other player's fleet not found");
	}
	await gamePage.selectMapObject(otherPlayerFleet);
	await gamePage.expectSummaryText([
		otherPlayerFleet.mapObject.name,
		'Ship Count: 1',
		`Mass: ${otherPlayerFleet.spec?.shipDesignSpec?.mass}`,
		`Warp Speed: ${otherPlayerFleet.warpSpeed}`,
		universe.playerIntels[1].racePluralName
	]);

	const otherPlayerPlanet = universe.planets.find((planet) => planet.mapObject?.playerNum === 2);
	if (!otherPlayerPlanet?.mapObject) {
		throw new Error("other player's planet not found");
	}
	await gamePage.selectMapObject(otherPlayerPlanet);
	await gamePage.expectSummaryText(['Report is current', universe.playerIntels[1].racePluralName]);
});

test('shows mystery trader summary', async ({ testGamePage }) => {
	const { gamePage, universe } = await testGamePage('Kitchen Sink');

	await gamePage.zoomOut(3);

	const mysteryTrader = universe.mysteryTraders[0];
	if (!mysteryTrader?.mapObject) {
		throw new Error('mysteryTrader not found');
	}
	await gamePage.selectMapObject(mysteryTrader);
	await gamePage.expectSummaryText([
		'The trader requests interested parties to send it a feet with at least 5000kT of minerals on board to be absorbed into the trader. It offers technological assistance in return. Trader is traveling at Warp 7.'
	]);
});

test('discovers paired wormhole after scout travel', async ({ testGamePage }) => {
	const { page, gamePage, universe } = await testGamePage('Kitchen Sink');

	const wormhole1 = universe.wormholes[0];
	if (!wormhole1?.mapObject) {
		throw new Error('wormhole not found');
	}
	await gamePage.selectMapObject(wormhole1);
	await gamePage.expectSummaryText([
		'Location: (10, 10)',
		'Stability: Rock Solid',
		'Destination: unknown'
	]);

	await expect(gamePage.tile('Fleets In Orbit')).toBeVisible();
	await gamePage.clickTileButton('Fleets In Orbit', 'Goto');
	await gamePage.selectMapObject(wormhole1, { modifiers: ['Meta'] });

	const { universe: updatedUniverse } = await submitTurn(page);

	expect(updatedUniverse?.wormholes.length).toBe(2);

	await gamePage.zoomOut(3);
	await gamePage.selectMapObject(updatedUniverse?.fleets[0], { clickCount: 2 });
	await expect(gamePage.tile('Fleet Waypoints')).toContainText('Space: (60, 60)');

	const wormhole2 = updatedUniverse?.wormholes[1];
	if (!wormhole2?.mapObject) {
		throw new Error('wormhole2 not found');
	}
	await gamePage.selectMapObject(wormhole2);
	await gamePage.expectSummaryNotText('Destination: unknown');
});
