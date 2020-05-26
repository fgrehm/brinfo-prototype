package cmd

import (
	"github.com/fgrehm/brinfo/core"
	"github.com/fgrehm/brinfo/core/scrapers"
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
		ID:             "br-gov-pr",
		Host:           "www.aen.pr.gov.br",
		ArticleScraper: scrapers.DefaultArticleScraper,
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
}
