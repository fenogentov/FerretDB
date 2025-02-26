// Copyright 2021 FerretDB Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/FerretDB/FerretDB/integration/setup"
	"github.com/FerretDB/FerretDB/integration/shareddata"
)

func TestIndexesCompatList(t *testing.T) {
	t.Parallel()

	s := setup.SetupCompatWithOpts(t, &setup.SetupCompatOpts{
		Providers:                shareddata.AllProviders(),
		AddNonExistentCollection: true,
	})
	ctx, targetCollections, compatCollections := s.Ctx, s.TargetCollections, s.CompatCollections

	for i := range targetCollections {
		targetCollection := targetCollections[i]
		compatCollection := compatCollections[i]

		t.Run(targetCollection.Name(), func(t *testing.T) {
			t.Helper()
			t.Parallel()

			targetCur, targetErr := targetCollection.Indexes().List(ctx)
			compatCur, compatErr := compatCollection.Indexes().List(ctx)

			require.NoError(t, compatErr)
			require.NoError(t, targetErr)

			targetRes := FetchAll(t, ctx, targetCur)
			compatRes := FetchAll(t, ctx, compatCur)

			assert.Equal(t, compatRes, targetRes)

			// Also test specifications to check they are identical.
			targetSpec, targetErr := targetCollection.Indexes().ListSpecifications(ctx)
			compatSpec, compatErr := compatCollection.Indexes().ListSpecifications(ctx)

			require.NoError(t, compatErr)
			require.NoError(t, targetErr)

			assert.Equal(t, compatSpec, targetSpec)
		})
	}
}

func TestIndexesCompatCreate(t *testing.T) {
	setup.SkipForTigrisWithReason(t, "Indexes creation is not supported for Tigris")

	t.Parallel()

	for name, tc := range map[string]struct { //nolint:vet // for readability
		models     []mongo.IndexModel
		resultType compatTestCaseResultType // defaults to nonEmptyResult
		skip       string                   // optional, skip test with a specified reason
	}{
		"Empty": {
			models:     []mongo.IndexModel{},
			resultType: emptyResult,
		},
		"SingleIndex": {
			models: []mongo.IndexModel{
				{Keys: bson.D{{"v", -1}}},
			},
		},
		"SingleIndexMultiField": {
			models: []mongo.IndexModel{
				{Keys: bson.D{{"foo", 1}, {"bar", -1}}},
			},
		},
		"DuplicateID": {
			models: []mongo.IndexModel{
				{
					Keys: bson.D{{"_id", 1}}, // this index is already created by default
				},
			},
			skip: "https://github.com/FerretDB/FerretDB/issues/2311",
		},
		"DescendingID": {
			models: []mongo.IndexModel{
				{Keys: bson.D{{"_id", -1}}},
			},
			resultType: emptyResult,
		},
		"NonExistentField": {
			models: []mongo.IndexModel{
				{Keys: bson.D{{"field-does-not-exist", 1}}},
			},
		},
		"DotNotation": {
			models: []mongo.IndexModel{
				{Keys: bson.D{{"v.foo", 1}}},
			},
		},
		"DangerousKey": {
			models: []mongo.IndexModel{
				{
					Keys: bson.D{
						{"v", 1},
						{"foo'))); DROP TABlE test._ferretdb_database_metadata; CREATE INDEX IF NOT EXISTS test ON test.test (((_jsonb->'foo", 1},
					},
				},
			},
		},
		"SameKey": {
			models: []mongo.IndexModel{
				{Keys: bson.D{{"v", -1}, {"v", 1}}},
			},
			resultType: emptyResult,
		},
		"CustomName": {
			models: []mongo.IndexModel{
				{
					Keys:    bson.D{{"foo", 1}, {"bar", -1}},
					Options: new(options.IndexOptions).SetName("custom-name"),
				},
			},
		},

		"MultiDirectionDifferentIndexes": {
			models: []mongo.IndexModel{
				{Keys: bson.D{{"v", -1}}},
				{Keys: bson.D{{"v", 1}}},
			},
		},
		"MultiOrder": {
			models: []mongo.IndexModel{
				{Keys: bson.D{{"foo", -1}}},
				{Keys: bson.D{{"v", 1}}},
				{Keys: bson.D{{"bar", 1}}},
			},
		},
		"MultiSameKeyUsed": {
			models: []mongo.IndexModel{
				{Keys: bson.D{{"foo", 1}}},
				{Keys: bson.D{{"foo", 1}, {"v", 1}}},
				{Keys: bson.D{{"bar", 1}}},
			},
		},
		"BuildSameIndex": {
			models: []mongo.IndexModel{
				{Keys: bson.D{{"v", 1}}},
				{Keys: bson.D{{"v", 1}}},
			},
			resultType: emptyResult,
			skip:       "https://github.com/FerretDB/FerretDB/issues/2311",
		},
		"MultiWithInvalid": {
			models: []mongo.IndexModel{
				{
					Keys: bson.D{{"foo", 1}, {"bar", 1}, {"v", -1}},
				},
				{
					Keys: bson.D{{"v", -1}, {"v", 1}},
				},
			},
			resultType: emptyResult,
		},
		"SameKeyDifferentNames": {
			models: []mongo.IndexModel{
				{
					Keys:    bson.D{{"v", -1}},
					Options: new(options.IndexOptions).SetName("foo"),
				},
				{
					Keys:    bson.D{{"v", -1}},
					Options: new(options.IndexOptions).SetName("bar"),
				},
			},
			resultType: emptyResult,
		},
		"SameNameDifferentKeys": {
			models: []mongo.IndexModel{
				{
					Keys:    bson.D{{"foo", -1}},
					Options: new(options.IndexOptions).SetName("index-name"),
				},
				{
					Keys:    bson.D{{"bar", -1}},
					Options: new(options.IndexOptions).SetName("index-name"),
				},
			},
			resultType: emptyResult,
		},
	} {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			if tc.skip != "" {
				t.Skip(tc.skip)
			}

			t.Helper()
			t.Parallel()

			// Use per-test setup because createIndexes modifies collection state,
			// however, we don't need to run index creation test for all the possible collections.
			s := setup.SetupCompatWithOpts(t, &setup.SetupCompatOpts{
				Providers:                []shareddata.Provider{shareddata.Composites},
				AddNonExistentCollection: true,
			})
			ctx, targetCollections, compatCollections := s.Ctx, s.TargetCollections, s.CompatCollections

			var nonEmptyResults bool
			for i := range targetCollections {
				targetCollection := targetCollections[i]
				compatCollection := compatCollections[i]

				t.Run(targetCollection.Name(), func(t *testing.T) {
					t.Helper()

					targetRes, targetErr := targetCollection.Indexes().CreateMany(ctx, tc.models)
					compatRes, compatErr := compatCollection.Indexes().CreateMany(ctx, tc.models)

					if targetErr != nil {
						t.Logf("Target error: %v", targetErr)
						t.Logf("Compat error: %v", compatErr)

						// error messages are intentionally not compared
						AssertMatchesCommandError(t, compatErr, targetErr)

						return
					}
					require.NoError(t, compatErr, "compat error; target returned no error")

					assert.Equal(t, compatRes, targetRes)

					if compatErr == nil {
						nonEmptyResults = true
					}

					// List indexes to check they are identical after creation.
					targetCur, targetErr := targetCollection.Indexes().List(ctx)
					compatCur, compatErr := compatCollection.Indexes().List(ctx)

					require.NoError(t, compatErr)
					assert.Equal(t, compatErr, targetErr)

					targetIndexes := FetchAll(t, ctx, targetCur)
					compatIndexes := FetchAll(t, ctx, compatCur)

					assert.Equal(t, compatIndexes, targetIndexes)

					// List specifications to check they are identical after creation.
					targetSpec, targetErr := targetCollection.Indexes().ListSpecifications(ctx)
					compatSpec, compatErr := compatCollection.Indexes().ListSpecifications(ctx)

					require.NoError(t, compatErr)
					require.NoError(t, targetErr)

					require.NotEmpty(t, compatSpec)
					assert.Equal(t, compatSpec, targetSpec)
				})
			}

			switch tc.resultType {
			case nonEmptyResult:
				assert.True(t, nonEmptyResults, "expected non-empty results (some documents should be modified)")
			case emptyResult:
				assert.False(t, nonEmptyResults, "expected empty results (no documents should be modified)")
			default:
				t.Fatalf("unknown result type %v", tc.resultType)
			}
		})
	}
}

// TestIndexesCreateRunCommand tests specific behavior for index creation that can be only provided through RunCommand.
func TestIndexesCompatCreateRunCommand(t *testing.T) {
	setup.SkipForTigrisWithReason(t, "Indexes creation is not supported for Tigris")

	t.Parallel()

	ctx, targetCollections, compatCollections := setup.SetupCompat(t)
	targetCollection := targetCollections[0]
	compatCollection := compatCollections[0]

	for name, tc := range map[string]struct { //nolint:vet // for readability
		collectionName any
		indexName      any
		key            any
		resultType     compatTestCaseResultType // defaults to nonEmptyResult
		skip           string                   // optional, skip test with a specified reason
	}{
		"invalid-collection-name": {
			collectionName: 42,
			key:            bson.D{{"v", -1}},
			indexName:      "custom-name",
			resultType:     emptyResult,
		},
		"nil-collection-name": {
			collectionName: nil,
			key:            bson.D{{"v", -1}},
			indexName:      "custom-name",
			resultType:     emptyResult,
		},
		"index-name-not-set": {
			collectionName: "test",
			key:            bson.D{{"v", -1}},
			indexName:      nil,
			resultType:     emptyResult,
			skip:           "https://github.com/FerretDB/FerretDB/issues/2311",
		},
		"empty-index-name": {
			collectionName: "test",
			key:            bson.D{{"v", -1}},
			indexName:      "",
			resultType:     emptyResult,
			skip:           "https://github.com/FerretDB/FerretDB/issues/2311",
		},
		"non-string-index-name": {
			collectionName: "test",
			key:            bson.D{{"v", -1}},
			indexName:      42,
			resultType:     emptyResult,
		},
		"existing-name-different-key-length": {
			collectionName: "test",
			key:            bson.D{{"_id", 1}, {"v", 1}},
			indexName:      "_id_", // the same name as the default index
			skip:           "https://github.com/FerretDB/FerretDB/issues/2311",
		},
		"invalid-key": {
			collectionName: "test",
			key:            42,
			resultType:     emptyResult,
		},
		"empty-key": {
			collectionName: "test",
			key:            bson.D{},
			resultType:     emptyResult,
		},
		"key-not-set": {
			collectionName: "test",
			resultType:     emptyResult,
			skip:           "https://github.com/FerretDB/FerretDB/issues/2311",
		},
	} {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			if tc.skip != "" {
				t.Skip(tc.skip)
			}

			t.Helper()
			t.Parallel()

			indexesDoc := bson.D{}

			if tc.key != nil {
				indexesDoc = append(indexesDoc, bson.E{Key: "key", Value: tc.key})
			}

			if tc.indexName != nil {
				indexesDoc = append(indexesDoc, bson.E{"name", tc.indexName})
			}

			var targetRes bson.D
			targetErr := targetCollection.Database().RunCommand(
				ctx, bson.D{
					{"createIndexes", tc.collectionName},
					{"indexes", bson.A{indexesDoc}},
				},
			).Decode(&targetRes)

			var compatRes bson.D
			compatErr := compatCollection.Database().RunCommand(
				ctx, bson.D{
					{"createIndexes", tc.collectionName},
					{"indexes", bson.A{indexesDoc}},
				},
			).Decode(&compatRes)

			if targetErr != nil {
				t.Logf("Target error: %v", targetErr)
				t.Logf("Compat error: %v", compatErr)

				// error messages are intentionally not compared
				AssertMatchesCommandError(t, compatErr, targetErr)

				return
			}
			require.NoError(t, compatErr, "compat error; target returned no error")

			if tc.resultType == emptyResult {
				require.Nil(t, targetRes)
				require.Nil(t, compatRes)
			}

			assert.Equal(t, compatRes, targetRes)

			targetErr = targetCollection.Database().RunCommand(
				ctx, bson.D{{"listIndexes", tc.collectionName}},
			).Decode(&targetRes)

			compatErr = compatCollection.Database().RunCommand(
				ctx, bson.D{{"listIndexes", tc.collectionName}},
			).Decode(&targetRes)

			require.Nil(t, targetRes)
			require.Nil(t, compatRes)

			if targetErr != nil {
				t.Logf("Target error: %v", targetErr)
				t.Logf("Compat error: %v", compatErr)

				// error messages are intentionally not compared
				AssertMatchesCommandError(t, compatErr, targetErr)

				return
			}
			require.NoError(t, compatErr, "compat error; target returned no error")
		})
	}
}

func TestIndexesCompatDrop(t *testing.T) {
	setup.SkipForTigrisWithReason(t, "Indexes are not supported for Tigris")

	t.Parallel()

	for name, tc := range map[string]struct { //nolint:vet // for readability
		dropIndexName string                   // name of a single index to drop
		dropAll       bool                     // set true for drop all indexes, if true dropIndexName must be empty.
		resultType    compatTestCaseResultType // defaults to nonEmptyResult
		toCreate      []mongo.IndexModel       // optional, if not nil create indexes before dropping
	}{
		"DropAllCommand": {
			toCreate: []mongo.IndexModel{
				{Keys: bson.D{{"v", 1}}},
				{Keys: bson.D{{"foo", -1}}},
				{Keys: bson.D{{"bar", 1}}},
				{Keys: bson.D{{"pam.pam", -1}}},
			},
			dropAll: true,
		},
		"ID": {
			dropIndexName: "_id_",
			resultType:    emptyResult,
		},
		"AscendingValue": {
			toCreate: []mongo.IndexModel{
				{Keys: bson.D{{"v", 1}}},
			},
			dropIndexName: "v_1",
		},
		"DescendingValue": {
			toCreate: []mongo.IndexModel{
				{Keys: bson.D{{"v", -1}}},
			},
			dropIndexName: "v_-1",
		},
		"AsteriskWithDropOne": {
			toCreate: []mongo.IndexModel{
				{Keys: bson.D{{"v", -1}}},
			},
			dropIndexName: "*",
			resultType:    emptyResult,
		},
		"NonExistent": {
			dropIndexName: "nonexistent_1",
			resultType:    emptyResult,
		},
	} {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Helper()
			t.Parallel()

			if tc.dropAll {
				require.Empty(t, tc.dropIndexName, "index name must be empty when dropping all indexes")
			}

			// It's enough to use a single provider for drop indexes test as indexes work the same for different collections.
			s := setup.SetupCompatWithOpts(t, &setup.SetupCompatOpts{
				Providers:                []shareddata.Provider{shareddata.Composites},
				AddNonExistentCollection: true,
			})
			ctx, targetCollections, compatCollections := s.Ctx, s.TargetCollections, s.CompatCollections

			var nonEmptyResults bool
			for i := range targetCollections {
				targetCollection := targetCollections[i]
				compatCollection := compatCollections[i]

				t.Run(targetCollection.Name(), func(t *testing.T) {
					t.Helper()

					if tc.toCreate != nil {
						_, targetErr := targetCollection.Indexes().CreateMany(ctx, tc.toCreate)
						_, compatErr := compatCollection.Indexes().CreateMany(ctx, tc.toCreate)
						require.NoError(t, compatErr)
						require.NoError(t, targetErr)
					}

					var targetRes, compatRes bson.Raw
					var targetErr, compatErr error

					if tc.dropAll {
						targetRes, targetErr = targetCollection.Indexes().DropAll(ctx)
						compatRes, compatErr = compatCollection.Indexes().DropAll(ctx)
					} else {
						targetRes, targetErr = targetCollection.Indexes().DropOne(ctx, tc.dropIndexName)
						compatRes, compatErr = compatCollection.Indexes().DropOne(ctx, tc.dropIndexName)
					}

					require.Equal(t, compatErr, targetErr)
					require.Equal(t, compatRes, targetRes)

					if targetErr == nil {
						nonEmptyResults = true
					}

					// List indexes to see they are identical after drop.
					targetCur, targetErr := targetCollection.Indexes().List(ctx)
					compatCur, compatErr := compatCollection.Indexes().List(ctx)

					require.NoError(t, compatErr)
					require.Equal(t, compatErr, targetErr)

					targetIndexes := FetchAll(t, ctx, targetCur)
					compatIndexes := FetchAll(t, ctx, compatCur)

					require.Equal(t, compatIndexes, targetIndexes)
				})
			}

			switch tc.resultType {
			case nonEmptyResult:
				require.True(t, nonEmptyResults, "expected non-empty results (some documents should be modified)")
			case emptyResult:
				require.False(t, nonEmptyResults, "expected empty results (no documents should be modified)")
			default:
				t.Fatalf("unknown result type %v", tc.resultType)
			}
		})
	}
}

func TestIndexesCompatDropRunCommand(t *testing.T) {
	setup.SkipForTigrisWithReason(t, "Indexes are not supported for Tigris")

	t.Parallel()

	for name, tc := range map[string]struct { //nolint:vet // for readability
		toCreate []mongo.IndexModel // optional, if set, create the given indexes before drop is called
		toDrop   any                // required, index to drop

		resultType compatTestCaseResultType // optional, defaults to nonEmptyResult
		skip       string                   // optional, skip test with a specified reason
	}{
		"MultipleIndexesByName": {
			toCreate: []mongo.IndexModel{
				{Keys: bson.D{{"v", -1}}},
				{Keys: bson.D{{"v", 1}, {"foo", 1}}},
				{Keys: bson.D{{"v.foo", -1}}},
			},
			toDrop: bson.A{"v_-1", "v_1_foo_1"},
		},
		"DocumentIndex": {
			toCreate: []mongo.IndexModel{
				{Keys: bson.D{{"v", -1}}},
			},
			toDrop: bson.D{{"v", -1}},
		},
		"DropAllExpression": {
			toCreate: []mongo.IndexModel{
				{Keys: bson.D{{"v", -1}}},
				{Keys: bson.D{{"foo.bar", 1}}},
				{Keys: bson.D{{"foo", 1}, {"bar", 1}}},
			},
			toDrop: "*",
		},
		"MultipleKeyIndex": {
			toCreate: []mongo.IndexModel{
				{Keys: bson.D{{"_id", -1}, {"v", 1}}},
			},
			toDrop: bson.D{
				{"_id", -1},
				{"v", 1},
			},
		},
	} {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			if tc.skip != "" {
				t.Skip(tc.skip)
			}

			t.Helper()
			t.Parallel()

			require.NotNil(t, tc.toDrop, "toDrop must not be nil")

			// It's enough to use a single provider for drop indexes test as indexes work the same for different collections.
			s := setup.SetupCompatWithOpts(t, &setup.SetupCompatOpts{
				Providers:                []shareddata.Provider{shareddata.Composites},
				AddNonExistentCollection: true,
			})
			ctx, targetCollections, compatCollections := s.Ctx, s.TargetCollections, s.CompatCollections

			var nonEmptyResults bool
			for i := range targetCollections {
				targetCollection := targetCollections[i]
				compatCollection := compatCollections[i]

				t.Run(targetCollection.Name(), func(t *testing.T) {
					t.Helper()

					if tc.toCreate != nil {
						_, targetErr := targetCollection.Indexes().CreateMany(ctx, tc.toCreate)
						_, compatErr := compatCollection.Indexes().CreateMany(ctx, tc.toCreate)
						require.NoError(t, compatErr)
						require.NoError(t, targetErr)

						// List indexes to see they are identical after creation.
						targetCur, targetListErr := targetCollection.Indexes().List(ctx)
						compatCur, compatListErr := compatCollection.Indexes().List(ctx)

						require.NoError(t, compatListErr)
						require.NoError(t, targetListErr)

						targetList := FetchAll(t, ctx, targetCur)
						compatList := FetchAll(t, ctx, compatCur)

						require.Equal(t, compatList, targetList)
					}

					targetCommand := bson.D{
						{"dropIndexes", targetCollection.Name()},
						{"index", tc.toDrop},
					}

					compatCommand := bson.D{
						{"dropIndexes", compatCollection.Name()},
						{"index", tc.toDrop},
					}

					var targetRes bson.D
					targetErr := targetCollection.Database().RunCommand(ctx, targetCommand).Decode(&targetRes)

					var compatRes bson.D
					compatErr := compatCollection.Database().RunCommand(ctx, compatCommand).Decode(&compatRes)

					if targetErr != nil {
						t.Logf("Target error: %v", targetErr)
						t.Logf("Compat error: %v", compatErr)

						// error messages are intentionally not compared
						AssertMatchesCommandError(t, compatErr, targetErr)

						return
					}
					require.NoError(t, compatErr, "compat error; target returned no error")

					if tc.resultType == emptyResult {
						require.Nil(t, targetRes)
						require.Nil(t, compatRes)
					}

					require.Equal(t, compatRes, targetRes)

					if compatErr == nil {
						nonEmptyResults = true
					}

					// List indexes to see they are identical after deletion.
					targetCur, targetListErr := targetCollection.Indexes().List(ctx)
					compatCur, compatListErr := compatCollection.Indexes().List(ctx)

					require.NoError(t, compatListErr)
					assert.Equal(t, compatListErr, targetListErr)

					targetList := FetchAll(t, ctx, targetCur)
					compatList := FetchAll(t, ctx, compatCur)

					assert.Equal(t, compatList, targetList)
				})
			}

			switch tc.resultType {
			case nonEmptyResult:
				require.True(t, nonEmptyResults, "expected non-empty results (some indexes should be deleted)")
			case emptyResult:
				require.False(t, nonEmptyResults, "expected empty results (no indexes should be deleted)")
			default:
				t.Fatalf("unknown result type %v", tc.resultType)
			}
		})
	}
}
