package object

import (
	"bytes"
	"text/template"

	"github.com/quailyquaily/goldmark-enclave/core"
)

const podbeanTpl = `
<iframe title="Podbean episode player" allowtransparency="true" height="150" width="100%" style="border: none; min-width: min(100%, 430px);height:150px;" scrolling="no" data-name="pb-iframe-player" src="https://www.podbean.com/player-v2/?from=embed&i={{.EmbedID}}&share=1&download=0&fonts=Arial&skin={{.Skin}}&font-color=auto&rtl=0&logo_link=episode_page&btn-skin=c73a3a&size=150" loading="lazy"></iframe>
`

func GetPodbeanHtml(enc *core.Enclave) (string, error) {
	// Theme handling mirrors other providers: default light
	skin := "f6f6f6"
	if enc.Theme == "dark" {
		skin = "1b1b1b"
	}

	var err error
	ret := ""
	if enc.IframeDisabled {
		ret, err = GetNoIframeTplHtml(enc, string(enc.Image.Destination))
		if err != nil {
			return "", err
		}
	} else {
		t, err := template.New("podbean").Parse(podbeanTpl)
		if err != nil {
			return "", err
		}

		buf := bytes.Buffer{}
		if err = t.Execute(&buf, map[string]string{
			"EmbedID": enc.ObjectID,
			"Skin":    skin,
		}); err != nil {
			return "", err
		}
		ret = buf.String()
	}

	return ret, nil
}
