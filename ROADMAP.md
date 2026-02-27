# VANIA — Procedural Metroidvania

**Gameplay style**: Metroidvania — 2D side-scrolling platformer with ability-gated exploration, persistent progression, and interconnected backtracking worlds.

**Vision**: Ship an infinitely replayable, fully procedural Metroidvania — every sprite, sound, story, and room generated from a single seed — achieving feature parity with the reference-complete *venture* roguelike across all five setting genres.

---

## Genre Support

Every system must implement the `GenreSwitcher` interface to switch thematic presentation at runtime:

```go
type GenreSwitcher interface {
    SetGenre(genreID string) // genreID: "fantasy" | "scifi" | "horror" | "cyberpunk" | "postapoc"
}
```

This applies to all ECS Systems (renderer, audio, AI, physics-hazard, HUD, narrative). Components and Entities hold genre-tagged data; only Systems are required to implement `SetGenre()`.

| Genre ID    | World Concept            | Traversal Hazards                    | Ability Gates                        | Visual Theme                              |
|-------------|--------------------------|--------------------------------------|--------------------------------------|-------------------------------------------|
| `fantasy`   | Enchanted castle         | Magical barriers, crumbling bridges  | Vine-covered doorways, fairy locks   | Vine-covered towers, fairy-lit caverns    |
| `scifi`     | Derelict space hulk      | Zero-G sections, airlock pressure    | Hull-breach bulkheads, ion doors     | Exposed conduit corridors, star-field bg  |
| `horror`    | Haunted mansion          | Creaking collapsing floors, darkness | Spirit seals, sanity-bar barriers    | Candlelit halls, sanity-draining fog      |
| `cyberpunk` | Megastructure            | Data-stream platforms, voltage arcs  | Hacking-gated doors, firewall walls  | Neon corridors, holographic signage       |
| `postapoc`  | Collapsed bunker         | Rubble platforms, radiation zones    | Sealed blast doors, mutagen locks    | Concrete debris, makeshift rope bridges   |

---

## Phased Milestones

### v1.0 — Core Engine + Playable Single-Player

*Goal: ECS scaffold, seed-based PCG, rendering, and fully playable platforming in one genre (`fantasy` baseline).*

#### ECS Framework
- [ ] Component / Entity / System interfaces (`SetGenre(genreID string)` required on every **System**; see interface definition above)
- [ ] System execution ordering and dependency graph
- [ ] Entity lifecycle management (spawn, despawn, pooling)

#### Seed-Based Deterministic RNG
- [x] Master seed → subsystem seed derivation (`HashSeed` via SHA-256)
- [x] Per-subsystem isolated `math/rand` sources
- [x] Determinism test suite (same seed → same game)

#### Input System
- [x] Keyboard / gamepad mapping
- [ ] Rebindable controls (stored in config)
- [ ] Input buffering for responsive platformer feel

#### Platforming Physics
- [ ] Gravity, variable-height jump (hold-to-rise)
- [ ] Wall-slide and wall-jump
- [ ] Dash (horizontal burst with i-frames)
- [ ] Double-jump
- [ ] Glide (slow-fall toggle)
- [ ] Grapple hook (swing to anchor point)
- [ ] Coyote-time and jump-buffer tolerances

#### Rendering — Sprites & Animation
- [x] Ebiten sprite batcher
- [x] Frame-based animation state machine
- [x] Procedural pixel-art sprite generation (cellular automata + symmetry)
- [x] Tileset generation with biome themes
- [ ] `SetGenre()` on renderer to swap palette/tileset presets

#### Camera System
- [x] Smooth follow camera with room-lock
- [ ] Camera transition animations on room change
- [ ] Screen-shake on impact / explosion

#### Procedural Level Generation — Room Graph
- [x] Graph-based world (80–150 rooms, 4–6 biomes)
- [x] Critical path with ability gates every ~5 rooms
- [x] Side branches (optional exploration)
- [ ] Backtracking shortcuts that unlock as abilities are gained
- [ ] Per-genre room tile vocabulary via `SetGenre()`

#### Audio — Waveform Synthesis & SFX
- [x] Sine / square / sawtooth / triangle / noise waveforms
- [x] ADSR envelope system
- [x] SFX generation (jump, land, attack, hit, pickup, door, ambient)
- [ ] `SetGenre()` on audio to select thematic instrument packs

#### UI / HUD / Menus
- [x] Health bar, ability indicators, seed display
- [ ] Main menu, pause menu, options screen
- [ ] Genre-themed UI skin switchable via `SetGenre()`

#### Save / Load
- [x] Multiple save slots with checkpoint autosave
- [ ] Slot selection screen
- [ ] Seed embedded in save for reproducibility

#### Config / Settings
- [ ] Resolution, volume, key bindings persisted to disk
- [ ] CLI flags (`--seed`, `--play`, `--genre`)

---

### v2.0 — Core Systems (Combat, AI, Ability Progression, All 5 Genres)

*Goal: Complete the gameplay loop — fighting, dying, learning abilities, and exploring all five genre skins.*

#### Combat System — Melee & Ranged
- [x] Melee attack (swing hitbox, combo chain)
- [x] Ranged attack (projectile with falloff)
- [ ] Block / parry with timing window
- [ ] Knockback and stagger states
- [ ] Damage numbers HUD

#### Status Effects
- [ ] Burn, freeze, shock, poison, bleed, slow, haste
- [ ] Stack / duration management
- [ ] Genre-mapped variants (e.g., "irradiate" for `postapoc`, "haunt" for `horror`)
- [ ] `SetGenre()` renames and recolours status icons

#### AI — Behavior Trees (Platformer-Adapted)
- [x] Patrol, chase, flee, flying movement patterns
- [x] Group coordination (5 formation types)
- [x] Adaptive difficulty (learning behaviors)
- [ ] Platformer-aware pathfinding (ledge detection, wall awareness)
- [ ] Ranged enemy positioning (maintain preferred attack range)
- [ ] Flying enemy altitude management

#### Boss Gatekeeper Encounters
- [x] Multi-phase boss fights
- [ ] Boss guards new ability on defeat
- [ ] Ability-gate door behind boss unlocks after kill
- [ ] Genre-themed boss skins via `SetGenre()`

#### Ability Progression — Ability Tree with Gating
- [x] Ability unlock order (dash, double-jump, glide, grapple, wall-jump)
- [ ] Stat-upgrade nodes (max HP, attack, defense) interspersed with movement unlocks
- [ ] Ability tree UI with lock/unlock animations
- [ ] `SetGenre()` flavour text on abilities (e.g., "Phase Dash" for `scifi`, "Soul Step" for `horror`)

#### Inventory & Items
- [x] Item pickup, visible in treasure rooms
- [ ] Inventory screen (grid layout)
- [ ] Consumables (potions → genre-skinned: mana vials, stim-packs, holy water, etc.)
- [ ] Key items (gate openers, lore fragments)
- [ ] Equipment slots (weapon, charm, armour)

#### Loot / Drops
- [x] Enemy drop tables (items, currency, ability fragments)
- [ ] Drop rarity tiers (common, uncommon, rare, legendary)
- [ ] Drop VFX and audio feedback

#### Shops / Economy
- [ ] Merchant NPC in safe rooms (genre-skinned: wizard, trader, black-market dealer)
- [ ] Currency system (dropped by enemies, found in rooms)
- [ ] Item buyback

#### Quests / Objectives
- [ ] Primary quest (reach final boss, gain all abilities)
- [ ] Optional side objectives (find lore rooms, defeat bonus bosses)
- [ ] Objective tracker in HUD
- [ ] Genre-flavoured quest text via `SetGenre()`

#### All 5 Genres — Full Integration
- [ ] `SetGenre()` implemented on every system (renderer, audio, AI, HUD, physics hazards)
- [ ] Genre selection at game start (or seed-derived genre)
- [ ] Per-genre tileset / palette / SFX / music presets
- [ ] Per-genre platforming hazards (magic barriers, airlock vents, darkness, voltage floors, radiation clouds)

#### Narrative Generation — Full Pipeline
- [x] Theme/story/lore generation
- [x] Character and faction generation
- [ ] Room descriptions surfaced in HUD on entry
- [ ] Environmental text (signs, terminals, gravestones) as diegetic lore
- [ ] `SetGenre()` re-skins narrative vocabulary

---

### v3.0 — Visual Polish (Lighting, Particles, Weather, Post-Processing)

*Goal: Make the procedurally generated world feel alive and distinct per genre.*

#### Dynamic Lighting
- [ ] Point lights (torches, glowing pickups, spell effects)
- [ ] Shadow casting for solid tiles
- [ ] Ambient light level per biome / room
- [ ] Genre presets via `SetGenre()` (warm candlelight for `horror`, blue neon glow for `cyberpunk`)

#### Particle Effects (Enhanced)
- [x] Combat hit sparks, blood/impact splats
- [x] Movement dust / landing puff
- [ ] Environmental particles (embers, snow, floating data bits, spores)
- [ ] Genre-specific particle themes via `SetGenre()`

#### Weather System (13 Types)
- [ ] Rain, snow, fog, sandstorm, acid rain, electrical storm, blizzard, heatwave, meteor shower, void mist, data storm, spore cloud, ash fall
- [ ] Weather affects platforming (ice slips, wind pushback, visibility)
- [ ] Genre-appropriate subset per theme via `SetGenre()`

#### Enhanced Sprite Generation
- [ ] Animated tiles (water flow, lava bubble, flickering lights)
- [ ] Multi-frame enemy procedural animations (idle, walk, attack, death, hurt)
- [x] Player animation frames per ability state
- [ ] Genre palette overlays via `SetGenre()`

#### Music — Adaptive Multi-Layer
- [x] Dynamic layer system (exploration, combat, boss)
- [x] Biome-specific tracks
- [ ] Smooth cross-fade between states
- [ ] Genre instrument mapping via `SetGenre()` (lute → synthesizer → pipe organ → glitch-bass → distorted guitar)

#### Audio — Positional SFX
- [ ] Distance attenuation for offscreen sounds
- [ ] Left/right stereo panning
- [ ] Reverb preset per room size / material

#### Genre Post-Processing Presets
- [ ] `fantasy` — warm desaturated vignette, bloom on magic
- [ ] `scifi` — cool scanline overlay, chromatic aberration at edges
- [ ] `horror` — desaturate + red-tint at low sanity, film grain
- [ ] `cyberpunk` — neon bloom, CRT curvature, glitch artifacts
- [ ] `postapoc` — sepia wash, dust overlay, vignette

---

### v4.0 — Gameplay Expansion (Advanced Abilities, Companions, Storytelling, Map)

*Goal: Deepen the Metroidvania loop — richer exploration, narrative payoff, and companion mechanics.*

#### Map System
- [ ] Automap that reveals explored rooms
- [ ] Ability-gate annotations (locked doors shown with required ability icon)
- [ ] Map screen with zoom / pan
- [ ] Room type icons (save, shop, boss, lore)
- [ ] Genre-themed map aesthetic via `SetGenre()`

#### Advanced Movement Abilities
- [ ] Swim (fluid sections, required for underwater rooms)
- [ ] Climb (vertical surfaces — vines for `fantasy`, pipes for `scifi`, ropes for `postapoc`)
- [ ] Ground-pound / stomp (breaks weak floors, stuns enemies)
- [ ] Blink (short-range teleport through barriers)

#### Backtracking World — Full Implementation
- [ ] Previously locked rooms accessible after ability unlock
- [ ] Shortcut doors that open from inside
- [ ] Room state persistence (destroyed breakables, collected items stay collected)

#### Destructible Environments
- [ ] Breakable walls (grapple / bomb / strength ability required)
- [ ] Collapsible floors (weight trigger or damage)
- [ ] Hidden passages revealed by scan ability
- [ ] Genre-themed destruction FX via `SetGenre()`

#### Fluid Dynamics
- [ ] Water / lava / acid liquid volumes
- [ ] Buoyancy and swim controls in liquid
- [ ] Lava damage, acid status effect, water resets jumps
- [ ] Genre-mapped fluids via `SetGenre()` (water, coolant, toxic sludge, data stream, irradiated water)

#### Companion AI
- [ ] Follower NPC that assists in combat (ranged support)
- [ ] Companion ability (genre-skinned: fairy guide, drone, ghost, AI partner, survivor)
- [ ] Companion dialogue during exploration
- [ ] `SetGenre()` swaps companion model, voice lines, ability

#### Magic / Spell System (mapped to Ability System)
- [ ] Spell slots driven by ability unlocks
- [ ] MP / mana resource (genre: mana, energy, sanity charge, data, fuel)
- [ ] 8 combat spells as a separate offensive layer (e.g., fireball, ice spike, chain lightning) distinct from movement abilities — spells consume mana, movement abilities do not
- [ ] `SetGenre()` reskins spell visuals and names

#### Reputation / Alignment
- [ ] Moral choice events in lore rooms
- [ ] Alignment axis (e.g., hope ↔ despair for `horror`, order ↔ chaos for `cyberpunk`)
- [ ] Alignment affects NPC dialogue and ending text

#### Environmental Storytelling
- [ ] Scripted environmental vignettes (body with journal, broken terminal, collapsed altar)
- [ ] Audio logs / readable items surface in lore overlay
- [ ] Genre-flavoured props via `SetGenre()`

#### Books / Lore Codex
- [ ] In-game codex screen with discovered lore entries
- [ ] Procedurally generated lore texts per genre
- [ ] Unlock via exploration (hidden rooms, boss drops)

#### Mini-Games (Optional Challenge Rooms)
- [ ] Speedrun gauntlet room (timed obstacle course)
- [ ] Survival arena (wave-based with timer)
- [ ] Puzzle room (pressure-plate / switch logic)
- [ ] Genre-appropriate mini-game skins via `SetGenre()`

#### Tutorial System
- [ ] Contextual pop-up tutorials on first encounter (jump, dash, combat)
- [ ] Skippable for veteran players
- [ ] Genre-voiced hint text via `SetGenre()`

#### Achievement System (Enhanced)
- [x] 19 achievements across 6 categories
- [ ] In-game achievement notification toasts
- [ ] Achievement viewer screen

#### Crafting System
- [ ] Combine drops to craft upgrades at workbench rooms
- [ ] Recipe discovery (genre: spell tome, schematic, ritual, blueprint, manual)
- [ ] Limited crafting slots to encourage decision-making

#### Skill / Talent Trees
- [ ] Secondary talent grid (passive bonuses: move speed, crit %, cooldown reduction)
- [ ] Talent points from leveling / boss kills
- [ ] Genre-flavoured talent names via `SetGenre()`

#### XP / Leveling → Stat Upgrades
- [ ] XP from enemy kills and exploration
- [ ] Level-up grants stat point and talent point
- [ ] Stat increases (HP, ATK, DEF, SPD)

#### Character Archetypes (replaces venture's 35-class system: 15 base + 20 prestige)
- [ ] 5 starting archetypes (Warrior, Rogue, Mage, Ranger, Survivor) — reduced to match metroidvania convention where player class shapes starting ability and talent bias rather than branching class trees; ability-gate progression provides depth instead
- [ ] Each archetype has unique starting ability and talent bias
- [ ] Genre-skinned archetype names via `SetGenre()`

#### World Events (Timed Challenges)
- [ ] Periodic world events (invasion rooms, cursed zones, bonus loot windows)
- [ ] In-HUD event timer and description
- [ ] Genre-appropriate event flavour via `SetGenre()`

---

### v5.0+ — Multiplayer Co-Op, Social Features, Production Polish

*Goal: Co-op play, cross-server portals, and release-quality production.*

#### Multiplayer Co-Op (Client-Server)
- [ ] Authoritative server with client prediction
- [ ] 200–5000 ms latency tolerance (lag compensation)
- [ ] 2–4 player co-op in shared world instance
- [ ] Shared ability gates (any player's ability unlocks for the party)

#### E2E Encrypted Chat
- [ ] In-game text chat with NaCl/TLS encryption
- [ ] Emote system excluded (not meaningful in single-player-primary platformer)

#### Trading
- [ ] Player-to-player item trade screen
- [ ] Scam protection (both-confirm before swap)

#### Co-Op Parties (replaces venture's guilds)
- [ ] Party creation, invite, kick, promote
- [ ] Shared seed for reproducible co-op sessions
- [ ] Party HUD overlay

#### Federation / Cross-Server Co-Op
- [ ] Inter-server party invites via federated identity
- [ ] Cross-server portal rooms for cross-instance play

#### Mail System
- [ ] Async in-game mail (items, notes) between players

#### CI/CD — Multi-Platform Builds
- [ ] GitHub Actions: Linux / macOS / Windows / WASM / mobile (iOS + Android)
- [ ] Binary signing and notarization
- [ ] Docker image for headless generation testing

#### Cross-Platform Builds
- [ ] Linux (amd64, arm64)
- [ ] macOS (universal binary)
- [ ] Windows (amd64)
- [ ] WASM (browser playable demo)
- [ ] Mobile (iOS + Android via Ebiten mobile target)

#### Test Coverage ≥ 82%
*82% matches the venture reference target — sufficient to catch regressions in deterministic PCG without requiring exhaustive coverage of generated-content edge cases.*
- [x] PCG core: seed, cache, validation (100%)
- [x] Graphics generators: sprite, tileset, palette
- [x] Audio generators: waveform, SFX, music
- [ ] World generator: room graph, biome, ability gates
- [ ] Entity generator: enemy, boss, item, ability
- [ ] Engine integration: combat, physics, input
- [ ] Multiplayer: netcode, chat, trading

#### Documentation
- [ ] `CHANGELOG.md` — version history
- [ ] `CONTROLS.md` — default key bindings per platform
- [ ] `FAQ.md` — seed sharing, co-op setup, genre selection
- [ ] API docs via `go doc`

#### Mod Framework
- [ ] Plugin interface for custom generators (implements `Generator` interface)
- [ ] Mod loading from `mods/` directory
- [ ] Documentation for mod authors

#### Performance Optimisation
- [ ] Viewport culling (don't render / update off-screen entities)
- [ ] Sprite cache (deduplicate seed-identical sprites across rooms)
- [ ] Draw call batching per tileset
- [ ] Parallel subsystem generation (graphics / audio / narrative concurrently)

---

## Excluded Features

The following *venture* features are intentionally omitted from vania:

| Venture Feature        | Rationale                                                                                   |
|------------------------|---------------------------------------------------------------------------------------------|
| Vehicles               | Not style-appropriate — side-scrolling platformer traversal replaced by movement abilities  |
| Emotes                 | Not meaningful in a primarily single-player platformer without persistent social spaces      |
| Building / Housing     | Not style-appropriate — world is a fixed procedural dungeon, not player-built settlements   |
| Furniture              | Excluded with housing — no player-owned spaces to furnish                                   |
| Territory Control      | Excluded — no persistent open world or faction map; dungeon rooms don't support ownership   |

---

## Shared Infrastructure

The following portable packages from venture can be reused or ported directly:

| Package / Module              | Reuse Plan                                                                              |
|-------------------------------|-----------------------------------------------------------------------------------------|
| `pkg/engine` (ECS core)       | Direct port — component/entity/system interfaces with `SetGenre()` hook                |
| `pkg/procgen/genre`           | Direct reuse — `SetGenre()` dispatch table and genre ID constants                       |
| Audio synthesis engine        | Already implemented in `internal/audio`; extend with genre instrument presets           |
| Save / load serialisation     | Port venture's JSON state serialiser; add ability-unlock and room-explored bitmasks     |
| Networking / netcode          | Port client-server loop and lag-compensation for v5.0 co-op                            |
| CI/CD workflow templates      | Reuse venture's multi-platform GitHub Actions matrix (Linux/macOS/Win/WASM/mobile)     |
| Post-processing shader stubs  | Port genre preset definitions; implement Ebiten shader equivalents                      |
| Behavior tree framework       | Port and adapt for platformer AI (add ledge-detection, altitude nodes)                 |
| Federation / cross-server     | Port identity and portal-room protocol unchanged                                        |
