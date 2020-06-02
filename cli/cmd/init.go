package cmd

import (
	"github.com/fgrehm/brinfo/core"
	"github.com/fgrehm/brinfo/core/scrapers"
	xt "github.com/fgrehm/brinfo/core/scrapers/extractors"
	mem "github.com/fgrehm/brinfo/storage/inmemory"
)

var repo core.ContentSourceRepo

func init() {
	repo = mem.NewContentSourceRepo()
	repo.Register(&core.ContentSource{
		ID:             "br-gov-sp",
		Host:           "www.saopaulo.sp.gov.br",
		ArticleScraper: scrapers.DefaultArticleScraper,
	})
	repo.Register(&core.ContentSource{
		ID:             "br-gov-ac",
		Host:           "agencia.ac.gov.br",
		ArticleScraper: scrapers.DefaultArticleScraper,
	})
	repo.Register(&core.ContentSource{
		ID:             "br-gov-rs",
		Host:           "estado.rs.gov.br",
		ArticleScraper: scrapers.DefaultArticleScraper,
	})
	repo.Register(&core.ContentSource{
		ID:             "br-ses-rs",
		Host:           "saude.rs.gov.br",
		ArticleScraper: scrapers.DefaultArticleScraper,
	})
	repo.Register(&core.ContentSource{
		ID:             "br-gov-pb",
		Host:           "paraiba.pb.gov.br",
		ArticleScraper: scrapers.DefaultArticleScraper,
	})
	repo.Register(&core.ContentSource{
		ID:               "br-gov-pr",
		Host:             "www.aen.pr.gov.br",
		ForceContentType: `text/html; charset="UTF-8"`,
		ArticleScraper: core.CombinedArticleScraper(
			scrapers.DefaultArticleScraper,
			scrapers.CustomArticleScraper(scrapers.CustomArticleScraperConfig{
				PublishedAt: xt.TimeText("aside dl dd p"),
			}),
		),
	})
	repo.Register(&core.ContentSource{
		ID:             "br-gov-mg",
		Host:           "www.agenciaminas.mg.gov.br",
		ArticleScraper: scrapers.DefaultArticleScraper,
	})
	repo.Register(&core.ContentSource{
		ID:               "br-gov-pe",
		Host:             "www.pe.gov.br",
		ForceContentType: `text/html; charset="UTF-8"`,
		ArticleScraper:   scrapers.DefaultArticleScraper,
	})
	repo.Register(&core.ContentSource{
		ID:   "br-gov-ba",
		Host: "www.ba.gov.br",
		ArticleScraper: core.CombinedArticleScraper(
			scrapers.DefaultArticleScraper,
			scrapers.CustomArticleScraper(scrapers.CustomArticleScraperConfig{
				PublishedAt: xt.TimeText("#main-content .field--name-field-data-da-noticia"),
				Images:      xt.Images(".field--name-field-fotos-galeria-multimidia .field--item a", "href"),
			}),
		),
	})
	repo.Register(&core.ContentSource{
		ID:               "br-ses-ba",
		Host:             "www.saude.ba.gov.br",
		ForceContentType: `text/html; charset="ISO-8859-1"`,
		ArticleScraper: core.CombinedArticleScraper(
			scrapers.DefaultArticleScraper,
			scrapers.CustomArticleScraper(scrapers.CustomArticleScraperConfig{
				// TODO: Figure out why htmlinfo is not parsing the <head> properly here
				Title:       xt.Text("#titulo h1", false),
				PublishedAt: xt.TimeText("#conteudo .sc_event_date strong"),
			}),
		),
	})
}
