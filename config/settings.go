package config

import (
	"html/template"
	"kproxy/metadata"
	"net/http"
	"strings"
)

func _populateRuleIntoSlice(cacheRules []metadata.CacheRule, rule string) []metadata.CacheRule {
	if cacheRules == nil {
		return nil
	}

	var newCacheRules []metadata.CacheRule
	for _, cacheRule := range cacheRules {
		cacheRule.Rule = rule
		newCacheRules = append(newCacheRules, cacheRule)
	}

	return newCacheRules
}

func writeTemplate(name string, data interface{}, res http.ResponseWriter) {
	tmpl := template.Must(template.ParseFiles("config/templates/" + name + ".html"))
	res.Header().Add("Cache-Control", "no-cache")
	_ = tmpl.Execute(res, data)
}

func getSettings(res http.ResponseWriter, req *http.Request) {
	settings := metadata.GetSettings(req)
	settings.NeverCache = _populateRuleIntoSlice(settings.NeverCache, "never")
	settings.AlwaysCache = _populateRuleIntoSlice(settings.AlwaysCache, "always")
	writeTemplate("settings", settings, res)
}

func saveSettings(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		res.Header().Add("Allow", "POST")
		res.WriteHeader(405)
		return
	}

	_ = req.ParseForm()
	rule := req.Form.Get("rule")
	glob := req.Form.Get("glob")
	onlyTypes := req.Form.Get("only-types")
	if glob == "" || onlyTypes == "" || (rule != "always" && rule != "never") {
		res.WriteHeader(400)
		return
	}

	cacheRule := metadata.CacheRule{
		Glob:      glob,
		OnlyTypes: strings.Split(onlyTypes, ","),
	}

	settings := metadata.GetSettings(req)
	if rule == "always" {
		settings.AlwaysCache = append(settings.AlwaysCache, cacheRule)
	} else {
		settings.NeverCache = append(settings.NeverCache, cacheRule)
	}

	settings.Save()

	res.Header().Add("Location", "/settings")
	res.WriteHeader(307)
}

func deleteCacheRule(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("Cache-Control", "no-cache")

	glob := req.URL.Query().Get("glob")
	rule := req.URL.Query().Get("rule")
	if glob == "" || (rule != "always" && rule != "never") {
		res.WriteHeader(400)
		return
	}

	settings := metadata.GetSettings(req)
	var list *[]metadata.CacheRule
	if rule == "always" {
		list = &settings.AlwaysCache
	} else {
		list = &settings.NeverCache
	}

	var newList []metadata.CacheRule
	for _, cacheRule := range *list {
		if cacheRule.Glob != glob {
			newList = append(newList, cacheRule)
		}
	}

	if rule == "always" {
		settings.AlwaysCache = newList
	} else {
		settings.NeverCache = newList
	}

	settings.Save()
	res.Header().Add("Location", "/settings")
	res.WriteHeader(307)
}
