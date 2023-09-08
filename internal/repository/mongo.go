package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/dneto/sai-scout/internal/i18n"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var sets = []string{"set1", "set2", "set3", "set4", "set5", "set6", "set6cde", "set7", "set7b", "set8"}

const (
	collectionBundles = "bundles"
	collectionCards   = "cards"
	database          = "sai_scout"
)

func UpdateSetBundles(ctx context.Context, bundleVersion string, saveFunc func(context.Context, *SetBundle) error) error {
	err := downloadAll(ctx, downloadAllParams{
		Version:   bundleVersion,
		Sets:      sets,
		Languages: i18n.AsStringSlice(i18n.Locales),
	}, saveFunc)

	if err != nil {
		return fmt.Errorf("failed to download bundles from cdn: %w", err)
	}

	return nil
}

func createIndexes(ctx context.Context, coll *mongo.Collection) error {
	_, err := coll.Indexes().CreateOne(ctx, mongo.IndexModel{Keys: bson.D{{Key: "name", Value: "text"}}})
	if err != nil {
		return err
	}
	_, err = coll.Indexes().CreateOne(ctx, mongo.IndexModel{Keys: bson.D{{Key: "cardcode", Value: 1}}})
	return err
}

func InsertBuilder(cli *mongo.Client) func(context.Context, *SetBundle) error {
	db := cli.Database(database)
	return func(ctx context.Context, bundle *SetBundle) error {
		bundleCollection := db.Collection(collectionBundles)

		setFilter := bson.D{
			{Key: "set", Value: bundle.Set},
			{Key: "language", Value: bundle.Locale},
		}

		var oldBundle *SetBundle
		if err := bundleCollection.FindOne(ctx, setFilter).Decode(&oldBundle); err != nil {
			if !errors.Is(err, mongo.ErrNoDocuments) {
				return err
			}
		}

		if oldBundle == nil || bundle.LastModified.After(oldBundle.LastModified) {
			opts := options.Replace().SetUpsert(true)
			b := setBundleWrite{LastModified: bundle.LastModified, Locale: bundle.Locale, Set: bundle.Set, Version: bundle.Version}
			if _, err := bundleCollection.ReplaceOne(ctx, setFilter, b, opts); err != nil {
				return err
			}

			replaceOnes := make([]mongo.WriteModel, len(bundle.Cards))
			for i, c := range bundle.Cards {
				cardFilter := bson.D{
					{Key: "cardcode", Value: c.CardCode},
				}
				replaceOnes[i] = mongo.NewReplaceOneModel().
					SetFilter(cardFilter).
					SetReplacement(c).
					SetUpsert(true)
			}

			cardsCollection := db.Collection(cardCollection(bundle.Locale))
			if err := createIndexes(ctx, cardsCollection); err != nil {
				return fmt.Errorf("failed to create index: %w", err)
			}

			if _, err := cardsCollection.BulkWrite(ctx, replaceOnes); err != nil {
				return fmt.Errorf("failed to upsert cards: %w", err)
			}
		}
		return nil
	}
}

func FindCardsBuilder(cli *mongo.Client) func(context.Context, string, ...string) ([]*Card, error) {
	db := cli.Database(database)
	return func(ctx context.Context, language string, codes ...string) ([]*Card, error) {
		inCards := make(bson.A, len(codes))
		for i, c := range codes {
			inCards[i] = c
		}
		pipeline := bson.A{
			bson.D{{
				Key: "$match", Value: bson.D{{
					Key: "cardcode", Value: bson.D{{
						Key:   "$in",
						Value: inCards,
					}},
				}},
			}},
		}
		for _, i := range customFieldsPipeline() {
			pipeline = append(pipeline, i)
		}

		coll := db.Collection(cardCollection(language))
		c, err := coll.Aggregate(ctx, pipeline)
		if err != nil {
			return nil, err
		}
		var cards []*Card
		err = c.All(ctx, &cards)
		return cards, err
	}
}

func SearchByNameBuilder(cli *mongo.Client) func(context.Context, string, string) ([]*Card, error) {
	db := cli.Database(database)
	return func(ctx context.Context, language string, name string) ([]*Card, error) {
		if len(strings.Split(name, " ")) > 1 {
			name = fmt.Sprintf(`"%s"`, name)
		}
		pipeline := bson.A{
			bson.D{{Key: "$match",
				Value: bson.D{{
					Key: "$text", Value: bson.D{{
						Key: "$search", Value: name,
					}},
				}},
			}},
			bson.D{{Key: "$limit", Value: 25}},
		}

		for _, i := range customFieldsPipeline() {
			pipeline = append(pipeline, i)
		}

		coll := db.Collection(cardCollection(language))
		c, err := coll.Aggregate(ctx, pipeline)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve card: %w", err)
		}

		var cards []*Card
		err = c.All(ctx, &cards)

		cards = lo.UniqBy(cards, func(c *Card) string {
			return c.Name
		})

		return cards, err
	}
}

func cardCollection(lang string) string {
	return fmt.Sprintf("%s_%s", collectionCards, lang)
}

func customFieldsPipeline() []bson.D {
	return []bson.D{
		{{Key: "$lookup",
			Value: bson.D{
				{Key: "from", Value: "cards_en_us"},
				{Key: "localField", Value: "cardcode"},
				{Key: "foreignField", Value: "cardcode"},
				{Key: "as", Value: "en_us"},
				{Key: "pipeline",
					Value: bson.A{
						bson.D{
							{Key: "$project",
								Value: bson.D{
									{Key: "type", Value: 1},
									{Key: "supertype", Value: 1},
								},
							},
						},
					},
				},
			},
		}},
		{{Key: "$unwind", Value: "$en_us"}},
		{
			{Key: "$addFields",
				Value: bson.D{
					{Key: "supertyperef", Value: "$en_us.supertype"},
					{Key: "typeref", Value: "$en_us.type"},
				},
			},
		},
		{{Key: "$project", Value: bson.D{{Key: "result", Value: 0}}}},
		{{Key: "$sort", Value: bson.D{{Key: "cardcode", Value: 1}}}},
	}
}
