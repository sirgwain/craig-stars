# Deviations from original game

`craig-stars`, as a clone of the original 1995 Stars! game, strives to be as close to the original as possible in terms of base game mechanics. However, certain elements either do not translate well to a modern web game or are simply not worth keeping.

Here is a non-exhaustive list of differences between `craig-stars` and the original game. Features yet to be implemented are not included.

## Medium-ish changes

- Upgrading a starbase now checks all slots in _both_ designs when calculating refunds, rather than only the corresponding slot in the new design[^1]. Additionally, part refund/transfer checks still occur when swapping hulls (so adding a component and swapping hulls on consecutive turns costs the same regardless of order).

[^1]: For refernce, Stars! charges extra for merely moving components to a different slot.

<!-- TODO: Check and fix more bugs -->

## Bugs fixed

- [0.2% min damage bug](https://wiki.starsautohost.org/wiki/Known_Bugs#0.2%_Minimum_Damage)
- [Colonization Module Check](https://wiki.starsautohost.org/wiki/Known_Bugs#Colonization_Module_Check)
- [Cheap Starbase](https://wiki.starsautohost.org/wiki/Known_Bugs#Cheap_Starbase)

## Minor changes

- The `DeltaPopulation` variable responsible for tracking pop growths below multiples of 100 did not reset when abandoning or invading a planet. This is fixed in `craig-stars`.
