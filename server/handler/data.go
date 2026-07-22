package handler

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// DataHandler merges data/<appID>/metadata.json, tech-stack.json, and the
// per-product files referenced from metadata.json's product list into a
// single response — the same shape the old combine-wbs.js build step used
// to produce, but computed at request time so no Node/build step is needed.
// appID is resolved as: ?app= query override > the server's configured default.
func DataHandler(dataFS fs.FS, defaultApp string) gin.HandlerFunc {
	return func(c *gin.Context) {
		appID := c.DefaultQuery("app", defaultApp)

		merged, err := buildAppData(dataFS, appID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, merged)
	}
}

// AppsHandler lists the dataset folder names available under data/, so
// callers can discover what's available to pass as ?app=. A directory only
// counts as a dataset if it actually holds a metadata.json — that keeps
// support folders (e.g. data/schema/, which holds JSON Schema files) from
// being offered as selectable datasets that would 404 on /api/data.
func AppsHandler(dataFS fs.FS) gin.HandlerFunc {
	return func(c *gin.Context) {
		entries, err := fs.ReadDir(dataFS, ".")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		apps := []string{}
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			if _, err := fs.Stat(dataFS, e.Name()+"/metadata.json"); err != nil {
				continue
			}
			apps = append(apps, e.Name())
		}
		c.JSON(http.StatusOK, gin.H{"apps": apps})
	}
}

func readJSONObject(dataFS fs.FS, path string) (map[string]interface{}, error) {
	b, err := fs.ReadFile(dataFS, path)
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}
	return m, nil
}

func firstNonNil(vals ...interface{}) interface{} {
	for _, v := range vals {
		if v != nil {
			return v
		}
	}
	return nil
}

func buildAppData(dataFS fs.FS, appID string) (map[string]interface{}, error) {
	base := appID

	meta, err := readJSONObject(dataFS, base+"/metadata.json")
	if err != nil {
		return nil, fmt.Errorf("unknown dataset %q: %w", appID, err)
	}

	projectConfig, _ := meta["projectConfig"].(map[string]interface{})
	if projectConfig == nil {
		projectConfig = meta
	}
	fpConfig, _ := meta["fpConfig"].(map[string]interface{})
	if fpConfig == nil {
		fpConfig = map[string]interface{}{}
	}
	effortConfig, _ := meta["effortConfig"].(map[string]interface{})
	if effortConfig == nil {
		effortConfig = map[string]interface{}{}
	}
	cfgBlock, _ := meta["config"].(map[string]interface{})

	var productsList []interface{}
	if pl, ok := projectConfig["products"].([]interface{}); ok {
		productsList = pl
	} else if pl, ok := meta["products"].([]interface{}); ok {
		productsList = pl
	}

	var techStack map[string]interface{}
	if ts, err := readJSONObject(dataFS, base+"/tech-stack.json"); err == nil {
		techStack = ts
	}

	statusCounts := map[string]int{
		"completed": 0, "partial": 0, "in-progress": 0,
		"roadmap": 0, "deprecated": 0, "beta": 0, "onhold": 0,
	}
	totalFeatures, totalCapabilities := 0, 0

	products := []interface{}{}
	for _, pmRaw := range productsList {
		pm, ok := pmRaw.(map[string]interface{})
		if !ok {
			continue
		}
		dataFile, _ := pm["dataFile"].(string)
		if dataFile == "" {
			continue
		}
		productData, err := readJSONObject(dataFS, base+"/"+dataFile)
		if err != nil {
			continue // matches combine-wbs.js's warn-and-skip behavior
		}
		if sc, ok := pm["shortCode"]; ok {
			productData["shortCode"] = sc
		}
		if d, ok := pm["description"]; ok {
			productData["description"] = d
		}

		if features, ok := productData["features"].([]interface{}); ok {
			totalFeatures += len(features)
			for _, fRaw := range features {
				f, ok := fRaw.(map[string]interface{})
				if !ok {
					continue
				}
				caps, ok := f["capabilities"].([]interface{})
				if !ok {
					continue
				}
				totalCapabilities += len(caps)
				for _, capRaw := range caps {
					cap, ok := capRaw.(map[string]interface{})
					if !ok {
						continue
					}
					status, _ := cap["status"].(string)
					if status == "" {
						status = "completed"
					}
					if _, known := statusCounts[status]; known {
						statusCounts[status]++
					}
				}
			}
		}

		products = append(products, productData)
	}

	metadata := map[string]interface{}{
		"title":           firstNonNil(projectConfig["title"], meta["title"]),
		"version":         meta["version"],
		"createdDate":     meta["createdDate"],
		"generatedDate":   time.Now().Format("2006-01-02"),
		"description":     meta["description"],
		"mainProductName": firstNonNil(projectConfig["mainProductName"], meta["mainProductName"]),
		"brandPrefix":     firstNonNil(projectConfig["brandPrefix"], meta["brandPrefix"]),
		"organization":    firstNonNil(projectConfig["organization"], meta["organization"]),
		"levels":          firstNonNil(projectConfig["levels"], meta["levels"]),
	}

	var cfgDefaultRate, cfgCurrency, cfgHoursPerDay, cfgDaysPerMonth, cfgDefaultPDR interface{}
	if cfgBlock != nil {
		cfgDefaultRate = cfgBlock["defaultHourlyRate"]
		cfgCurrency = cfgBlock["currency"]
		cfgHoursPerDay = cfgBlock["hoursPerDay"]
		cfgDaysPerMonth = cfgBlock["daysPerMonth"]
		cfgDefaultPDR = cfgBlock["defaultPDR"]
	}
	config := map[string]interface{}{
		"defaultHourlyRate": firstNonNil(projectConfig["defaultHourlyRate"], cfgDefaultRate, 75.0),
		"currency":          firstNonNil(projectConfig["currency"], cfgCurrency, "USD"),
		"hoursPerDay":       firstNonNil(projectConfig["hoursPerDay"], cfgHoursPerDay, 8.0),
		"daysPerMonth":      firstNonNil(projectConfig["daysPerMonth"], cfgDaysPerMonth, 20.0),
		// Base hours per Function Point — the unadjusted human baseline for the
		// stack. Tech-stack and productivity factors compose on top of it in the
		// app (effective PDR = base x tech x productivity), so this must NOT
		// already have an AI discount baked in. 8.0 is the conventional
		// human-driven default. Previously this was read by the app but never
		// passed through here, so any dataset value was silently discarded.
		"defaultPDR": firstNonNil(projectConfig["defaultPDR"], cfgDefaultPDR, 8.0),
	}

	result := map[string]interface{}{
		"appId":                appID,
		"metadata":             metadata,
		"projectSummary":       meta["projectSummary"],
		"config":               config,
		"statusDefinitions":    meta["statusDefinitions"],
		"techStackFactors":     firstNonNil(effortConfig["techStackFactors"], meta["techStackFactors"]),
		"productivityFactors":  firstNonNil(effortConfig["productivityFactors"], meta["productivityFactors"]),
		"sdlcPhases":           firstNonNil(effortConfig["sdlcPhases"], []interface{}{}),
		"fpWeights":      firstNonNil(fpConfig["fpWeights"], meta["fpWeights"]),
		"gscDefinitions": firstNonNil(fpConfig["gscDefinitions"], []interface{}{}),
		// VAF = vafBase + (vafIncrement x TDI). The IFPUG standard is
		// 0.65 + 0.01 x TDI, which bounds the adjustment to +/-35%. Exposed as
		// config so nothing is hardcoded, but changing these means the result is
		// no longer IFPUG-conformant — see docs/ai-dlc-estimation-model.md.
		"vafBase":      firstNonNil(fpConfig["vafBase"], 0.65),
		"vafIncrement": firstNonNil(fpConfig["vafIncrement"], 0.01),
		"glossary":             meta["glossary"],
		"products":             products,
		"summary": map[string]interface{}{
			"totalProducts":     len(products),
			"totalFeatures":     totalFeatures,
			"totalCapabilities": totalCapabilities,
			"statusCounts":      statusCounts,
		},
	}
	if techStack != nil {
		result["techStack"] = techStack["products"]
		result["techStackSummary"] = techStack["summary"]
	}

	return result, nil
}
