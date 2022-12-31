import { PlayerSettings } from '$lib/types/PlayerSettings';
import { writable } from 'svelte/store';

export const settings = writable<PlayerSettings>(loadSettingsOrDefault());

settings.subscribe((value) => {
	localStorage.setItem('playerSettings', JSON.stringify(value));
});

function loadSettingsOrDefault(): PlayerSettings {
	const json = localStorage.getItem('playerSettings');
	if (json) {
		const settings = JSON.parse(json) as PlayerSettings;
		if (settings) {
			return settings;
		}
	}

	return new PlayerSettings();
}
