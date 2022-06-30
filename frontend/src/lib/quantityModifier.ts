import hotkeys from 'hotkeys-js';

let tenModifier = 1;
let hundredModifier = 1;

export const bindQuantityModifier = () => {
	hotkeys('*', { keyup: true }, function (event) {
		if (hotkeys.shift) {
			if (event.type === 'keydown') {
				tenModifier = 10;
			} else if (event.type === 'keyup') {
				tenModifier = 1;
			}
		}
		if (hotkeys.control || hotkeys.command) {
			if (event.type === 'keydown') {
				hundredModifier = 100;
			} else if (event.type === 'keyup') {
				hundredModifier = 1;
			}
		}
	});
};

export const unbindQuantityModifier = () => {
	hotkeys.unbind('*');
};

export const getQuantityModifier = () => tenModifier * hundredModifier;
