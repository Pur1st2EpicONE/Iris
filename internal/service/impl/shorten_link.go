package impl

import (
	"Iris/internal/errs"
	"Iris/internal/models"
	"context"
	"errors"

	"github.com/jxskiss/base62"
)

const offset = 1000000

func (s *Service) ShortenLink(ctx context.Context, link models.Link) (string, error) {

	if err := validateLink(link); err != nil {
		return "", err
	}

	if link.Alias != "" {
		return s.saveWithAlias(ctx, link)
	}

	id, err := s.storage.SaveOriginal(ctx, link.OriginalURL)
	if err != nil {
		s.logger.LogError("service — failed to save original link",
			err, "link", link.OriginalURL, "layer", "service.impl")
		return "", err
	}

	shortLink := encode(id + offset)

	if err := s.storage.SaveShort(ctx, id, shortLink); err != nil {
		s.logger.LogError("service — failed to save short link", err,
			"link", link.OriginalURL, "short link", shortLink, "layer", "service.impl")
		return "", err
	}

	return shortLink, nil

}

func (s *Service) saveWithAlias(ctx context.Context, link models.Link) (string, error) {

	err := s.storage.SaveWithAlias(ctx, link)
	if err == nil {
		return link.Alias, nil
	}

	if errors.Is(err, errs.ErrAliasExists) {
		return "", err
	}

	s.logger.LogError("service — failed to save link with alias", err,
		"link", link.OriginalURL, "alias", link.Alias, "layer", "service.impl")

	return "", err
}

func encode(id int64) string {
	return string(base62.FormatUint(uint64(id)))
}
