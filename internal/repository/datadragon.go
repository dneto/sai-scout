package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/sourcegraph/conc/pool"
)

const cdn = "dd.b.pvp.net"

type downloadAllParams struct {
	Version   string
	Sets      []string
	Languages []string
}

type SetBundle struct {
	LastModified time.Time
	Locale       string
	Set          string
	Version      string
	Cards        []*Card
}

type setBundleWrite struct {
	LastModified time.Time
	Locale       string
	Set          string
	Version      string
}

func downloadAll(ctx context.Context, params downloadAllParams, saveFunc func(context.Context, *SetBundle) error) error {
	setPool := pool.New().WithErrors().WithMaxGoroutines(10)
	for _, language := range params.Languages {
		for _, set := range params.Sets {
			setPool.Go(func() error {
				bundle, err := downloadSetBundle(ctx, params.Version, set, language)
				if err != nil {
					return err
				}
				return saveFunc(context.Background(), bundle)
			})
		}
	}
	return setPool.Wait()
}

func downloadSetBundle(ctx context.Context, version string, set string, locale string) (*SetBundle, error) {
	baseURL := fmt.Sprintf("https://%s/%s", cdn, strings.Replace(version, ".", "_", -1))
	url := baseURL + fmt.Sprintf("/%s/%s/data/%s-%s.json", set, locale, set, locale)

	client := http.DefaultClient
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode > 400 {
		return nil, fmt.Errorf("failed to retrieve data: status code %d", resp.StatusCode)
	}

	var cards []*Card
	err = json.NewDecoder(resp.Body).Decode(&cards)
	if err != nil {
		return nil, err
	}
	lastModifiedHeader := resp.Header.Get("Last-Modified")
	lastModified, err := http.ParseTime(lastModifiedHeader)
	if err != nil {
		lastModified = time.Now()
	}
	log.Debug().Str("set", set).Str("language", locale).Msg("download successful")
	return &SetBundle{Cards: cards, LastModified: lastModified, Locale: locale, Set: set, Version: version}, nil
}
