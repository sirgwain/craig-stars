package cs

import (
	"encoding/json"
	"slices"
)

// a list of all TechTags that benefit ships in combat
var CombatTechTags = []TechTag{
	TechTagArmor,
	TechTagBeamCapacitor,
	TechTagBeamDeflector,
	TechTagBeamWeapon,
	TechTagCapitalShipMissile,
	TechTagEngine,
	TechTagGatlingGun,
	TechTagInitiativeBonus,
	TechTagManeuveringJet,
	TechTagShieldSapper,
	TechTagShield,
	TechTagTorpedo,
	TechTagTorpedoJammer,
	TechTagTorpedoBonus,
}

// A collection of an object's TechTags (like on a tech part)
type TechTags map[TechTag]struct{}

// Create a new TechTags map from a list of TechTag items, or an empty map if none are specified
func newTechTags(tags ...TechTag) TechTags {
	newTechTags := TechTags{}
	for _, t := range tags {
		newTechTags[t] = struct{}{}
	}
	return newTechTags
}

// MarshalJSON implements json.Marshaler to marshal TechTags as an array.
func (tt TechTags) MarshalJSON() ([]byte, error) {
	tags := make([]TechTag, 0, len(tt))
	for tag := range tt {
		tags = append(tags, tag)
	}
	slices.Sort(tags)
	return json.Marshal(tags)
}

// UnmarshalJSON implements json.Unmarshaler.
func (tt TechTags) UnmarshalJSON(data []byte) error {
	var tags []any
	err := json.Unmarshal(data, &tags)
	if err != nil {
		return err
	}
	for _, tag := range tags {
		s, ok := tag.(string)
		if !ok {
			panic("could not convert json data back to TechTags")
		}
		tt[TechTag(s)] = struct{}{}
	}
	return nil
}

// returns true if tt has at least 1 of the specified TechTags
// and none of the tags in tagsToExclude
//
// Tag(s) contained in both lists will not be banned
func (tt TechTags) hasTags(tagsToInclude []TechTag, tagsToExclude ...TechTag) bool {
	blacklist := newTechTags(tagsToExclude...)
	whitelist := newTechTags(tagsToInclude...)
	hasTag := false

	for _, tag := range tt.GetTags() {
		switch {
		case whitelist.HasTag(tag):
			hasTag = true
		case blacklist.HasTag(tag):
			// our TechTags has a blacklisted tag not
			// also in our whitelist; automatic fail
			return false
		}
	}
	return hasTag
}

// return true if tt has this tag
func (tt TechTags) HasTag(tag TechTag) bool {
	_, ok := tt[tag]
	return ok
}

// return unsorted list of all tags in tt
func (tt TechTags) GetTags() []TechTag {
	list := make([]TechTag, 0, len(tt))
	for k := range tt {
		list = append(list, k)
	}
	return list
}

// return number of unique tags in tt
func (tt TechTags) Count() int {
	return len(tt)
}
