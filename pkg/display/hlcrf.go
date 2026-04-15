package display

import (
	"os"
	"path/filepath"
	"strings"
)

func (s *Service) buildHLCRFComponents(pageURL string) string {
	loaded, err := s.loadManifestForOrigin(pageURL)
	if err != nil || loaded == nil {
		return ""
	}
	var scripts []string
	for _, component := range loaded.Manifest.HLCRF {
		templateBody := strings.TrimSpace(component.Template)
		if templateBody == "" && strings.TrimSpace(component.Name) != "" {
			body, readErr := os.ReadFile(filepath.Join(loaded.BaseDir, component.Name))
			if readErr == nil {
				templateBody = string(body)
			}
		}
		if templateBody == "" {
			continue
		}
		tag := strings.TrimSpace(component.Tag)
		if tag == "" {
			tag = defaultHLCRFTag(component.Name)
		}
		scripts = append(scripts, renderHLCRFComponent(tag, templateBody))
	}
	return strings.Join(scripts, "\n")
}

func renderHLCRFComponent(tag, templateBody string) string {
	return `(function(){if(customElements.get(` + quoteJS(tag) + `)){return;}const tpl=document.createElement('template');tpl.innerHTML=` +
		quoteJS(templateBody) +
		`;class CoreHLCRFElement extends HTMLElement{connectedCallback(){if(this.shadowRoot){return;}const root=this.attachShadow({mode:'open'});root.appendChild(tpl.content.cloneNode(true));}}customElements.define(` +
		quoteJS(tag) +
		`,CoreHLCRFElement);})();`
}

func defaultHLCRFTag(name string) string {
	name = strings.TrimSpace(strings.ToLower(name))
	name = strings.TrimSuffix(name, filepath.Ext(name))
	name = strings.ReplaceAll(name, "_", "-")
	name = strings.ReplaceAll(name, " ", "-")
	if !strings.Contains(name, "-") {
		name = "core-" + name
	}
	return name
}
