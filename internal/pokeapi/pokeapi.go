package pokeapi

import (
	"encoding/json"
	"fmt"
	pokecache "github.com/avgra3/pokedexcli/internal/pokecache"
	"io"
	"log"
	"net/http"
)

func GetLocations(url string, cache *pokecache.Cache) LocationResult {
	// Check cache has data
	if cacheEntry, ok := cache.Get(url); ok {
		fmt.Println("We are using the cache...")
		// We get back a []bytes which we would want to convert to a location result
		cachedResult := LocationResult{}
		err := json.Unmarshal(cacheEntry, &cachedResult)
		if err != nil {
			log.Fatal(err)
		}
		return cachedResult

	}

	// Cache has no data, call the api
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	// Add new data to cache
	cache.Add(url, body)

	if res.StatusCode > 299 {
		statusCodeMessage := fmt.Sprintf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
		log.Fatal(statusCodeMessage)
	}
	if err != nil {
		log.Fatal(err)
	}
	locationResults := LocationResult{}
	err = json.Unmarshal(body, &locationResults)
	if err != nil {
		log.Fatal(err)
	}
	return locationResults
}

type Locations struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type LocationResult struct {
	Count    int         `json:"count"`
	Next     string      `json:"next"`
	Previous string      `json:"previous"`
	Results  []Locations `json:"results"`
}

type Generation struct {
}

type Pokedex struct {
	Id             int            `json:"id"`
	Name           string         `json:"name"`
	IsMainSeries   bool           `json:"is_main_series"`
	Descriptions   []Description  `json:"descriptions"`
	Names          []LanguageName `json:"names"`
	PokemonEntries []PokemonEntry `json:"pokemon_entries"`
	Region         Region         `json:"region"`
	VersionGroups  []VersionGroup `json:"version_groups"`
}

type Description struct {
	Description string       `json:"description"`
	Language    LanguageName `json:"language"`
}

type PokemonSpecies struct {
	Id                   int                     `json:"id"`
	Name                 string                  `json:"name"`
	Order                int                     `json:"order"`
	GenderRate           int                     `json:"gender_rate"`
	CaptureRate          int                     `json:"capture_rate"`
	BaseHappiness        int                     `json:"base_happiness"`
	IsBaby               bool                    `json:"is_baby"`
	IsLegendary          bool                    `json:"is_legendary"`
	IsMythical           bool                    `json:"is_mythical"`
	HatchEncounter       int                     `json:"hatch_encounter"`
	HasGenderDifferences bool                    `json:"has_gender_differences"`
	FormsSwitchable      bool                    `json:"froms_switchable"`
	GrowthRate           []GrowthRate            `json:"growth_rate"`
	EggGroups            []EggGroup              `json:"egg_groups"`
	Color                PokemonColor            `json:"color"`
	Shape                PokemonShape            `json:"pokemon_shape"`
	EvolvesFromSpecies   *PokemonSpecies         `json:"evolves_from_species"`
	EvolutionChain       EvolutionChain          `json:"evolution_chain"`
	Habitat              PokemonHabitat          `json:"habitat"`
	Generation           Generation              `json:"generation"`
	Names                []LanguageName          `json:"names"`
	PalParkEncounters    []PalParkEncounterArea  `json:"pal_park_encounters"`
	FlavorTextEntries    []FlavorText            `json:"flavor_text_entries"`
	FormDescriptions     []Description           `json:"form_descriptions"`
	Genera               []Genus                 `json:"genera"`
	Varieties            []PokemonSpeciesVariety `json:"varieties"`
}

type GrowthRateExperienceLevel struct {
	Level      int `json:"level"`
	Experience int `json:"experience"`
}

type GrowthRate struct {
	Id             int                         `json:"id"`
	Name           string                      `json:"name"`
	Formula        string                      `json:"formula"`
	Descriptions   []Description               `json:"description"`
	Levels         []GrowthRateExperienceLevel `json:"levels"`
	PokemonSpecies []PokemonSpecies            `json:"pokemon_species"`
}

type EggGroup struct {
	Id             int              `json:"id"`
	Name           string           `json:"name"`
	Names          []LanguageName   `json:"names"`
	PokemonSpecies []PokemonSpecies `json:"pokemon_species"`
}

type PokemonColor struct {
	Id             string           `json:"id"`
	Name           string           `json:"name"`
	Names          []LanguageName   `json:"names"`
	PokemonSpecies []PokemonSpecies `json:"pokemon_species"`
}

type PokemonShape struct {
	Id             int              `json:"id"`
	Name           string           `json:"name"`
	AwesomeNames   []AwesomeName    `json:"awesome_names"`
	Names          LanguageName     `json:"names"`
	PokemonSpecies []PokemonSpecies `json:"pokemon_species"`
}

type Language struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Official bool   `json:"official"`
	Iso639   string `json:"iso639"`
	Iso3166  string `json:"iso3166"`
	Names    []LanguageName
}

type AwesomeName struct {
	AwesomeName string   `json:"awesome_name"`
	Language    Language `json:"language"`
}

type EvolutionChain struct {
	Id              int        `json:"id"`
	BabyTriggerItem Item       `json:"baby_trigger_item"`
	Chain           *ChainLink `json:"chain"`
}

type Item struct {
	Id                int                      `json:"id"`
	Name              string                   `json:"name"`
	Cost              int                      `json:"cost"`
	FlingPower        int                      `json:"fling_power"`
	FlingEffect       ItemFlingEffect          `json:"fling_effect"`
	Attributes        []ItemAttribute          `json:"attributes"`
	Category          []ItemCategory           `json:"category"`
	EffectEntries     []VerboseEffect          `json:"effect_entries"`
	FlavorTextEntries []VersionGroupFlavorText `json:"flavor_text_entries"`
	GameIndices       []GenerationGameIndex    `json:"game_indicies"`
	Names             []LanguageName           `json:"names"`
	Sprites           ItemSprites              `json:"sprites"`
	HeldByPokemon     ItemHolderPokemon        `json:"held_by_pokemon"`
	BabyTriggerFor    *EvolutionChain          `json:"baby_trigger_for"`
	Machines          []MachineVersionDetail   `json:"machines"`
}
type ChainLink struct {
	IsBaby           bool            `json:"is_baby"`
	Species          PokemonSpecies  `json:"species"`
	EvolutionDetails EvolutionDetail `json:"evolution_details"`
	EvolvesTo        []ChainLink     `json:"evolves_to"`
}
type EvolutionDetail struct {
	Item                  Item             `json:"item"`
	Trigger               EvolutionTrigger `json:"trigger"`
	Gender                int              `json:"gender"`
	HeldItem              Item             `json:"held_item"`
	KnownMove             Move             `json:"known_move"`
	KnownMoveType         Type             `json:"known_move_type"`
	Location              Location         `json:"location"`
	MinLevel              int              `json:"min_level"`
	MinHappiness          int              `json:"min_happiness"`
	MinBeauty             int              `json:"min_beauty"`
	MinAffection          int              `json:"min_affection"`
	NeedsOverworldRain    bool             `json:"needs_overworld_rain"`
	PartySpecies          PokemonSpecies   `json:"party_species"`
	PartyType             Type             `json:"party_type"`
	RelativePhysicalStats int              `json:"relative_physical_stats"`
	TimeOfDay             string           `json:"time_of_day"`
	TradeSpecies          PokemonSpecies   `json:"trade_species"`
	TurnUpsideDown        bool             `json:"turn_upside_down"`
}

type TypeRelations struct {
	NoDamageTo       []Type `json:"no_damage_to"`
	HalfDamageTo     []Type `json:"half_damage_to"`
	DoubleDamageTo   []Type `json:"double_damage_to"`
	NoDamageFrom     []Type `json:"no_damage_from"`
	HalfDamageFrom   []Type `json:"half_damage_from"`
	DoubleDamageFrom []Type `json:"double_damage_from"`
}
type TypeRelationsPast struct {
	Generation      Generation    `json:"generation"`
	DamageRelations TypeRelations `json:"damage_relations"`
}

type MoveDamageClass struct {
	Id           int            `json:"id"`
	Name         string         `json:"name"`
	Descriptions []Description  `json:"descriptions"`
	Moves        []Move         `json:"moves"`
	Names        []LanguageName `json:"names"`
}

type TypePokemon struct {
	Slot    int     `json:"slot"`
	Pokemon Pokemon `json:"pokemon"`
}

type Type struct {
	Id                  int                   `json:"id"`
	Name                string                `json:"name"`
	DamageRelations     TypeRelations         `json:"damage_relations"`
	PastDamageRelations TypeRelationsPast     `json:"past_damage_relations"`
	GameIndicies        []GenerationGameIndex `json:"game_indicies"`
	Generation          Generation            `json:"generation"`
	MoveDamageClass     MoveDamageClass       `json:"move_damage_class"`
	Names               []Language            `json:"names"`
	Pokemon             TypePokemon           `json:"pokemon"`
	Moves               []Move                `json:"moves"`
}
type EvolutionTrigger struct {
	Id             int            `json:"id"`
	Name           string         `json:"name"`
	Names          []LanguageName `json:"names"`
	PokemonSpecies PokemonSpecies `json:"pokemon_species"`
}

type Move struct {
	Id                 int                    `json:"id"`
	Name               string                 `json:"name"`
	Accuracy           int                    `json:"accuracy"`
	EffectChance       int                    `json:"effect_chance"`
	Pp                 int                    `json:"pp"`
	Priority           int                    `json:"priority"`
	Power              int                    `json:"power"`
	ContestCombos      ContestComboSets       `json:"contest_combos"`
	ContestType        ContestType            `json:"contest_type"`
	ContestEffect      ContestEffect          `json:"contest_effect"`
	DamageClass        MoveDamageClass        `json:"damage_class"`
	EffectEntries      VerboseEffect          `json:"effect_entries"`
	EffectChanges      AbilityEffectChange    `json:"effect_changes"`
	LearnedByPokemon   []AbilityEffectChange  `json:"learned_by_pokemon"`
	FlavorTextEntries  []MoveFlavorText       `json:"flavor_text_entries"`
	Generation         Generation             `json:"generation"`
	Machines           []MachineVersionDetail `json:"machines"`
	Meta               MoveMetaData           `json:"meta"`
	Names              []LanguageName         `json:"names"`
	PastValues         []PastMoveStatValues   `json:"past_values"`
	StatChanges        []MoveStatChange       `json:"stat_changes"`
	SuperContestEffect MoveStatChange         `json:"stat_contest_effect"`
	Target             *MoveTarget            `json:"target"`
	Type               Type                   `json:"type"`
}

type PastMoveStatValues struct {
	Accuracy      int           `json:"accuracy"`
	EffectChance  int           `json:"effect_chance"`
	Power         int           `json:"power"`
	Pp            int           `json:"pp"`
	EffectEntries VerboseEffect `json:"effect_entries"`
	Type          Type          `json:"type"`
	VersionGroup  VersionGroup  `json:"version_gruop"`
}

type MoveStatChange struct {
	Change int  `json:"change"`
	Stat   Stat `json:"stat"`
}

type Stat struct {
	Id               int                  `json:"id"`
	Name             string               `json:"name"`
	GameIndex        int                  `json:"game_index"`
	IsBattleOnly     bool                 `json:"is_battle_only"`
	AffectingMoves   MoveStatAffectSets   `json:"affecting_moves"`
	AffectingNatures NatureStatAffectSets `json:"affecting_natures"`
	Characteristics  []Characteristic     `json:"characteristic"`
	MoveDamageClass  MoveDamageClass      `json:"move_damage_class"`
	Names            []LanguageName       `json:"names"`
}

type MoveStatAffectSets struct {
	Increase []MoveStatAffect `json:"increase"`
	Decrease []MoveStatAffect `json:"decrease"`
}

type MoveStatAffect struct {
	Change int  `json:"change"`
	Move   Move `json:"move"`
}

type NatureStatAffectSets struct {
	Increase Nature `json:"increase"`
	Decrease Nature `json:"decrease"`
}

type Nature struct {
	Id                         int                         `json:"id"`
	Name                       string                      `json:"name"`
	DecreasedStat              *Stat                       `json:"decreased_stat"`
	IncreasedStat              *Stat                       `json:"increased_stat"`
	HatesFlavor                BerryFlavor                 `json:"hates_flavor"`
	LikesFlavor                BerryFlavor                 `json:"likes_flavor"`
	PokeathlonStatChanges      []NatureStatChange          `json:"pokeathlon_stat_changes"`
	MoveBattleStylePreferences []MoveBattleStylePreference `json:"move_battle_style_preferences"`
	Names                      []LanguageName              `json:"names"`
}

type MoveBattleStylePreference struct {
	LowHpPreference  int             `json:"low_hp_preference"`
	HighHpPreference int             `json:"high_hp_preference"`
	MoveBattleStyle  MoveBattleStyle `json:"move_battle_style"`
}

type MoveBattleStyle struct {
	Id    int            `json:"id"`
	Name  string         `json:"name"`
	Names []LanguageName `json:"names"`
}

type NatureStatChange struct {
	MaxChange      int            `json:"max_change"`
	PokeathlonStat PokeathlonStat `json:"pokeathlon_stat"`
}

type PokeathlonStat struct {
	Id               int                            `json:"id"`
	Name             string                         `json:"name"`
	Names            LanguageName                   `json:"names"`
	AffectingNatures NaturePokeathlonStatAffectSets `json:"affecting_natures"`
}

type NaturePokeathlonStatAffectSets struct {
	Increase NaturePokeathlonStatAffect `json:"increase"`
	Decrease NaturePokeathlonStatAffect `json:"decrease"`
}

type NaturePokeathlonStatAffect struct {
	MaxChange int    `json:"max_change"`
	Nature    Nature `json:"nature"`
}

type BerryFlavor struct {
	Id          int            `json:"id"`
	Name        string         `json:"name"`
	Berries     FlavorBerryMap `json:"berries"`
	ContestType *ContestType   `json:"contest_type"`
	Names       []LanguageName `json:"names"`
}

type FlavorBerryMap struct {
	Potency int   `json:"potency"`
	Berry   Berry `json:"berry"`
}

type Berry struct {
	Id               int             `json:"id"`
	Name             string          `json:"name"`
	GrowthTime       int             `json:"growth_time"`
	MaxHarvest       int             `json:"max_harvest"`
	NaturalGiftPower int             `json:"natural_gift_power"`
	Size             int             `json:"size"`
	Smoothness       int             `json:"smoothness"`
	SoilDryness      int             `json:"soil_dryness"`
	Firmness         BerryFirmness   `json:"firmness"`
	Flavors          *FlavorBerryMap `json:"flavors"`
	Item             Item            `json:"item"`
	NaturalGiftType  Type            `json:"natural_gift_type"`
}

type BerryFirmness struct {
	Id      int            `json:"id"`
	Name    string         `json:"name"`
	Berries []Berry        `json:"berries"`
	Names   []LanguageName `json:"names"`
}

type Characteristic struct {
	Id             int         `json:"id"`
	GeneModulo     int         `json:"gene_modulo"`
	PossibleValues []int       `json:"possible_values"`
	HighestStat    Stat        `json:"highest_stat"`
	Descriptions   Description `json:"descriptions"`
}

type MoveTarget struct {
	Id           int           `json:"id"`
	Name         string        `json:"name"`
	Descriptions []Description `json:"descriptions"`
	Moves        Move          `json:"move"`
	Names        LanguageName  `json:"names"`
}

type MoveMetaData struct {
	Ailment       MoveAilment  `json:"ailment"`
	Category      MoveCategory `json:"category"`
	MinHits       int          `json:"min_hits"`
	MaxHits       int          `json:"max_hits"`
	MinTurns      int          `json:"min_turns"`
	MaxTurns      int          `json:"max_turns"`
	Drain         int          `json:"drain"`
	Healing       int          `json:"healing"`
	CritRate      int          `json:"crit_rate"`
	AilmentChance int          `json:"ailment_chance"`
	FlinchChance  int          `json:"flinch_chance"`
	StatChance    int          `json:"stat_chance"`
}

type MoveAilment struct {
	Id    int            `json:"id"`
	Name  string         `json:"name"`
	Moves []Move         `json:"moves"`
	Names []LanguageName `json:"names"`
}

type MoveCategory struct {
	Id           int           `json:"id"`
	Name         string        `json:"name"`
	Moves        []Move        `json:"moves"`
	Descriptions []Description `json:"descriptions"`
}

type MoveFlavorText struct {
	FlavorText   string       `json:"flavor_text"`
	Language     Language     `json:"language"`
	VersionGroup VersionGroup `json:"version_group"`
}

type AbilityEffectChange struct {
	EffectEntries []Effect     `json:"effect_entries"`
	VersionGroup  VersionGroup `json:"version_group"`
}

type Effect struct {
	Effect   string   `json:"effect"`
	Language Language `json:"language"`
}

type ContestEffect struct {
	Id                int          `json:"id"`
	Appeal            int          `json:"appeal"`
	Jam               int          `json:"jam"`
	EffectEntries     []Effect     `json:"effect_entries"`
	FlavorTextEntries []FlavorText `json:"flavor_text_entries"`
}

type ContestType struct {
	Id          int           `json:"id"`
	Name        string        `json:"name"`
	BerryFlavor BerryFlavor   `json:"berry_flavor"`
	Names       []ContestName `json:"names"`
}

type ContestName struct {
	Name     string   `json:"name"`
	Color    string   `json:"color"`
	Language Language `json:"language"`
}

type ContestComboSets struct {
	Normal ContestComboDetail `json:"normal"`
	Super  ContestComboDetail `json:"super"`
}

type ContestComboDetail struct {
	UseBefore []Move `json:"use_before"`
	UseAfter  []Move `json:"use_after"`
}

type ItemSprites struct {
	Default string `json:"default"`
}

type ItemFlingEffect struct {
	Id            int      `json:"id"`
	Name          string   `json:"name"`
	EffectEntries []Effect `json:"effect_entries"`
	Items         []Item   `json:"items"`
}

type ItemAttribute struct {
	Id           int            `json:"id"`
	Name         string         `json:"name"`
	Items        []Item         `json:"items"`
	Names        []LanguageName `json:"names"`
	Descriptions []Description  `json:"descriptions"`
}

type ItemCategory struct {
	Id     int            `json:"id"`
	Name   string         `json:"name"`
	Items  []Item         `json:"items"`
	Names  []LanguageName `json:"names"`
	Pocket ItemPocket     `json:"pocket"`
}

type ItemPocket struct {
	Id         int            `json:"id"`
	Name       string         `json:"name"`
	Categories []ItemCategory `json:"categories"`
	Names      []LanguageName `json:"names"`
}

type VerboseEffect struct {
	Effect      string   `json:"effect"`
	ShortEffect string   `json:"short_effect"`
	Language    Language `json:"language"`
}

type VersionGroupFlavorText struct {
	Text         string       `json:"text"`
	Language     Language     `json:"language"`
	VersionGroup VersionGroup `json:"version_group"`
}

type ItemHolderPokemon struct {
	Pokemon        Pokemon                          `json:"pokemon"`
	VersionDetails []ItemHolderPokemonVersionDetail `json:"version_details"`
}

type ItemHolderPokemonVersionDetail struct {
	Rarity  int     `json:"rarity"`
	Version Version `json:"version"`
}

type MachineVersionDetail struct {
	Machine      Machine      `json:"machine"`
	VersionGroup VersionGroup `json:"version_group"`
}

type Machine struct {
	Id           int          `json:"id"`
	Item         Item         `json:"item"`
	Move         Move         `json:"move"`
	VersionGroup VersionGroup `json:"version_group"`
}

type PokemonHabitat struct {
	Id             int              `json:"id"`
	Name           string           `json:"name"`
	Names          []LanguageName   `json:"names"`
	PokemonSpecies []PokemonSpecies `json:"pokemon_species"`
}

type PalParkEncounterArea struct {
	BaseScore int         `json:"base_score"`
	Rate      int         `json:"rate"`
	Area      PalParkArea `json:"area"`
}

type PalParkArea struct {
	Id                int                       `json:"id"`
	Name              string                    `json:"name"`
	Names             []LanguageName            `json:"names"`
	PokemonEncounters []PalParkEncounterSpecies `json:"pokemon_encounters"`
}

type PalParkEncounterSpecies struct {
	BaseScore      int            `json:"base_score"`
	Rate           int            `json:"rate"`
	PokemonSpecies PokemonSpecies `json:"pokemon_species"`
}

type FlavorText struct {
	FlavorText string   `json:"flavor_text"`
	Language   Language `json:"language"`
	Version    Version  `json:"version"`
}

type Genus struct {
	Genus    string   `json:"genus"`
	Language Language `json:"language"`
}

type PokemonSpeciesVariety struct {
	IsDefault bool    `json:"is_default"`
	Pokemon   Pokemon `json:"pokemon"`
}

type PokemonEntry struct {
	EntryNumber    int            `json:"entry_number"`
	PokemonSpecies PokemonSpecies `json:"pokemon_species"`
}

type VersionGroup struct{}

type Region struct {
	Id             int            `json:"id"`
	Locations      []Locations    `json:"locations"`
	Name           string         `json:"name"`
	Names          []LanguageName `json:"names"`
	MainGeneration Generation     `json:"main_generation"`
	Pokedexes      []Pokedex      `json:"pokedexes"`
	VersionGroups  []VersionGroup `json:"version_groups"`
}

type GenerationGameIndex struct{}

type Location struct {
	Id           int                   `json:"id"`
	Name         string                `json:"name"`
	Region       Region                `json:"region"`
	Names        []LanguageName        `json:"names"`
	GameIndicies []GenerationGameIndex `json:"game_indicies"`
	Areas        []LocationArea        `json:"areas"`
}

type LocationArea struct {
	Id                   int    `json:"id"`
	Name                 string `json:"name"`
	GameIndex            int    `json:"game_index"`
	EncounterMethodRates []EncounterMethodRates
	Location             Location
}

type SpecificPokemon struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type PokemonEncounters struct {
	EncounterMethodRates []EncounterMethodRates `json:"encounter_method_rates"`
	GameIndex            int                    `json:"game_index"`
	Id                   int                    `json:"id"`
	Location             Locations              `json:"location"`
	LocationName         string                 `json:"name"`
	Names                []LanguageName         `json:"names"`
	PokemonEncounters    []PokemonEncounter     `json:"pokemon_encounters"`
}

type PokemonEncounter struct {
	Pokemon        Pokemon          `json:"pokemon"`
	VersionDetails []VersionDetails `json:"version_details"`
}

type Pokemon struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type LanguageName struct {
	Name     string `json:"name"`
	Language string `json:"language"`
}

type EncounterMethodRates struct {
	EncounterMethod         EncounterMethod          `json:"encounter_method"`
	VersionEncounterDetails []VersionEncounterDetail `json:"version_details"`
}

type VersionEncounterDetail struct {
	Version          Version            `json:"version"`
	MaxChance        int                `json:"chance"`
	EncounterDetails []EncounterDetails `json:"encoutner_details"`
}

type EncounterDetails struct {
	MinLevel        int                     `json:"min_level"`
	MaxLevel        int                     `json:"max_level"`
	Conditionvalues EncounterConditionValue `json:"condition_values"`
	Chance          int                     `json:"chance"`
	Method          EncounterMethod         `json:"method"`
}

type EncounterMethod struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type EncounterCondition struct {
	Id     int                       `json:"id"`
	Name   string                    `json:"name"`
	Names  []LanguageName            `json:"names"`
	Values []EncounterConditionValue `json:"values"`
}

type EncounterConditionValue struct {
	Id        int                `json:"id"`
	Name      string             `json:"name"`
	Condition EncounterCondition `json:"condition"`
	Names     []LanguageName     `json:"names"`
}

type VersionDetails struct {
	Rate    int     `json:"rate"`
	Version Version `json:"version"`
}

type Version struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}
