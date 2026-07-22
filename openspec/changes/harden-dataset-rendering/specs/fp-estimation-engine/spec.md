# FP Estimation Engine

## ADDED Requirements

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
