# FP Estimation Engine

## Requirements

### Requirement: Generic WBS rendering from dataset content only
The app layer SHALL render its title, branding, product list, features, capabilities, FP weights, GSC configuration, and status labels entirely from the dataset served at `GET /api/data`. The app layer SHALL NOT contain any hardcoded reference to a specific project's name, brand, or domain vocabulary.

#### Scenario: Switching datasets changes all visible branding
- **GIVEN** the server is configured to serve the `tripma` dataset instead of `ai-agents-provly`
- **WHEN** the page loads
- **THEN** the title, footer text, and export filename prefix all reflect `tripma`'s `metadata.brandPrefix`
- **AND** no text from the previous dataset remains visible

### Requirement: Status legend labels are dataset-driven
The system SHALL render each status legend label from `wbsData.statusDefinitions[key].label`, falling back to a generic English default only when a dataset omits that field.

#### Scenario: Dataset relabels a status
- **GIVEN** a dataset's `metadata.json` sets `statusDefinitions.roadmap.label` to `"Build New"`
- **WHEN** the status legend renders
- **THEN** the roadmap filter shows "Build New", not a hardcoded default
- **AND** no app-layer code change was required to achieve this

#### Scenario: Dataset omits statusDefinitions
- **GIVEN** a dataset's `metadata.json` has no `statusDefinitions` block
- **WHEN** the status legend renders
- **THEN** each label falls back to its generic default ("Completed", "Partial", "In Progress", "Roadmap", "Beta")

### Requirement: Capability-level inclusion toggle
The system SHALL provide checkboxes at the capability (3rd level) in the WBS tree to include or exclude individual capabilities from FP estimates, cascading from product/feature exclusion.

#### Scenario: Toggle individual capability
- **GIVEN** a user is viewing the WBS tree with a product and feature expanded
- **WHEN** the user clicks the checkbox next to a capability
- **THEN** that capability's inclusion state toggles
- **AND** the UFP/AFP totals update immediately

#### Scenario: Parent exclusion cascades and disables children
- **GIVEN** a feature is excluded (unchecked)
- **WHEN** the user views the capabilities under that feature
- **THEN** all capability checkboxes under it are disabled
- **AND** their FP is excluded from all totals regardless of their own stored state

### Requirement: Inclusion and filter state persist per dataset
The system SHALL persist product/feature/capability inclusion state, status filters, and GSC values to `localStorage`, namespaced by the active dataset's identity so multiple datasets can be used in the same browser without interfering with each other.

#### Scenario: State survives page reload for the same dataset
- **GIVEN** a user has excluded specific capabilities while viewing the `ai-agents-provly` dataset
- **WHEN** the page is reloaded with the same dataset active
- **THEN** the same capabilities remain excluded

#### Scenario: Two datasets don't share state
- **GIVEN** a user has excluded capabilities while viewing `ai-agents-provly`
- **WHEN** the user (or server config) switches to the `tripma` dataset in the same browser
- **THEN** `tripma`'s inclusion state is independent — unaffected by what was excluded under `ai-agents-provly`

### Requirement: FP calculation respects all inclusion levels
The FP calculation engine SHALL check inclusion state at product, feature, and capability levels before including a capability's function points in totals, and SHALL respect the active status filter.

#### Scenario: Excluded capability not counted
- **GIVEN** a capability with 10 FP is excluded
- **WHEN** the system calculates total UFP
- **THEN** that capability's 10 FP is NOT included in the total

#### Scenario: Export respects inclusion state
- **GIVEN** some capabilities are excluded
- **WHEN** the user exports to Excel or PDF
- **THEN** the exported totals match the on-screen calculated totals

### Requirement: A malformed or missing optional field degrades one section, not the whole app
Each top-level UI section (project summary, stats, WBS tree, FP types, productivity factors, tech factors, tech stack, glossary, GSC factors) SHALL render independently of the others. A rendering failure in one section SHALL be logged to the console naming the section, and SHALL NOT prevent any other section from rendering.

#### Scenario: One malformed journey doesn't blank the page
- **GIVEN** a dataset's `projectSummary.keyJourneys` contains one entry missing its `steps` field
- **WHEN** the page loads
- **THEN** that journey renders a "No steps defined for this journey" placeholder
- **AND** the WBS tree, stats, and every other tab render normally

### Requirement: GSC definitions are dataset-driven
The GSC (General System Characteristics) sliders SHALL use the active dataset's `fpConfig.gscDefinitions` when present and non-empty, including each factor's authored default rating. Only when the dataset omits or empties this field SHALL the app fall back to the standard IFPUG 14-factor list.

#### Scenario: Dataset-authored GSC default is respected
- **GIVEN** a dataset's `fpConfig.gscDefinitions` sets a factor's `default` to `5`
- **WHEN** the GSC Factors tab renders
- **THEN** that factor's slider initializes to `5`, not a hardcoded fallback value

#### Scenario: Dataset omits gscDefinitions
- **GIVEN** a dataset's `fpConfig.gscDefinitions` is absent, null, or `[]`
- **WHEN** the GSC Factors tab renders
- **THEN** the standard IFPUG 14-factor list is used, with no error
