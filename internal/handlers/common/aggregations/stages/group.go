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

package stages

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/FerretDB/FerretDB/internal/handlers/common"
	"github.com/FerretDB/FerretDB/internal/handlers/common/aggregations"
	"github.com/FerretDB/FerretDB/internal/handlers/common/aggregations/operators/accumulators"
	"github.com/FerretDB/FerretDB/internal/handlers/commonerrors"
	"github.com/FerretDB/FerretDB/internal/types"
	"github.com/FerretDB/FerretDB/internal/util/iterator"
	"github.com/FerretDB/FerretDB/internal/util/lazyerrors"
	"github.com/FerretDB/FerretDB/internal/util/must"
)

// group represents $group stage.
//
//	{ $group: {
//		_id: <groupExpression>,
//		<groupBy[0].outputField>: {accumulator0: expression0},
//		...
//		<groupBy[N].outputField>: {accumulatorN: expressionN},
//	}}
//
// $group uses group expression to group documents that have the same evaluated expression.
// The evaluated expression becomes the _id for that group of documents.
// For each group of documents, accumulators are applied.
type group struct {
	groupExpression any
	groupBy         []groupBy
}

// groupBy represents accumulation to apply on the group.
type groupBy struct {
	accumulate  func(iter types.DocumentsIterator) (any, error)
	outputField string
}

// newGroup creates a new $group stage.
func newGroup(stage *types.Document) (aggregations.Stage, error) {
	fields, err := common.GetRequiredParam[*types.Document](stage, "$group")
	if err != nil {
		return nil, commonerrors.NewCommandErrorMsgWithArgument(
			commonerrors.ErrStageGroupInvalidFields,
			"a group's fields must be specified in an object",
			"$group (stage)",
		)
	}

	var groupKey any
	var groups []groupBy

	iter := fields.Iterator()

	defer iter.Close()

	for {
		field, v, err := iter.Next()
		if errors.Is(err, iterator.ErrIteratorDone) {
			break
		}

		if err != nil {
			return nil, lazyerrors.Error(err)
		}

		if field == "_id" {
			if doc, ok := v.(*types.Document); ok {
				if err = validateDocumentExpression(doc); err != nil {
					return nil, err
				}
			}
			groupKey = v
			continue
		}

		accumulator, err := accumulators.NewAccumulator("$group", field, v)
		if err != nil {
			return nil, err
		}

		groups = append(groups, groupBy{
			outputField: field,
			accumulate:  accumulator.Accumulate,
		})
	}

	if groupKey == nil {
		return nil, commonerrors.NewCommandErrorMsgWithArgument(
			commonerrors.ErrStageGroupMissingID,
			"a group specification must include an _id",
			"$group (stage)",
		)
	}

	return &group{
		groupExpression: groupKey,
		groupBy:         groups,
	}, nil
}

// Process implements Stage interface.
func (g *group) Process(ctx context.Context, iter types.DocumentsIterator, closer *iterator.MultiCloser) (types.DocumentsIterator, error) { //nolint:lll // for readability
	docs, err := iterator.ConsumeValues(iterator.Interface[struct{}, *types.Document](iter))
	if err != nil {
		return nil, lazyerrors.Error(err)
	}

	groupedDocuments, err := g.groupDocuments(ctx, docs)
	if err != nil {
		return nil, err
	}

	var res []*types.Document

	for _, groupedDocument := range groupedDocuments {
		doc := must.NotFail(types.NewDocument("_id", groupedDocument.groupID))

		groupIter := iterator.Values(iterator.ForSlice(groupedDocument.documents))
		defer groupIter.Close()

		for _, accumulation := range g.groupBy {
			out, err := accumulation.accumulate(groupIter)
			if err != nil {
				return nil, err
			}

			if doc.Has(accumulation.outputField) {
				// document has duplicate key
				return nil, commonerrors.NewCommandErrorMsgWithArgument(
					commonerrors.ErrDuplicateField,
					fmt.Sprintf("duplicate field: %s", accumulation.outputField),
					"$group (stage)",
				)
			}

			doc.Set(accumulation.outputField, out)
		}

		res = append(res, doc)
	}

	return iterator.Values(iterator.ForSlice(res)), nil
}

// groupDocuments groups documents by group expression.
func (g *group) groupDocuments(ctx context.Context, in []*types.Document) ([]groupedDocuments, error) {
	groupKey, ok := g.groupExpression.(string)
	if !ok {
		// non-string key aggregates values of all `in` documents into one aggregated document.
		return []groupedDocuments{{
			groupID:   g.groupExpression,
			documents: in,
		}}, nil
	}

	expression, err := aggregations.NewExpression(groupKey)
	if err != nil {
		var exprErr *aggregations.ExpressionError
		if !errors.As(err, &exprErr) {
			return nil, lazyerrors.Error(err)
		}

		switch exprErr.Code() {
		case aggregations.ErrNotExpression:
			// constant value aggregates values of all `in` documents into one aggregated document.
			return []groupedDocuments{{
				groupID:   groupKey,
				documents: in,
			}}, nil
		case aggregations.ErrEmptyFieldPath:
			return nil, commonerrors.NewCommandErrorMsgWithArgument(
				// TODO
				commonerrors.ErrGroupInvalidFieldPath,
				"'$' by itself is not a valid Expression",
				"$group (stage)",
			)
		case aggregations.ErrInvalidExpression:
			return nil, commonerrors.NewCommandErrorMsgWithArgument(
				commonerrors.ErrFailedToParse,
				fmt.Sprintf("'%s' starts with an invalid character for a user variable name", types.FormatAnyValue(groupKey)),
				"$group (stage)",
			)
		case aggregations.ErrEmptyVariable:
			return nil, commonerrors.NewCommandErrorMsgWithArgument(
				commonerrors.ErrFailedToParse,
				"empty variable names are not allowed",
				"$group (stage)",
			)
		// TODO https://github.com/FerretDB/FerretDB/issues/2275
		case aggregations.ErrUndefinedVariable:
			return nil, commonerrors.NewCommandErrorMsgWithArgument(
				commonerrors.ErrGroupUndefinedVariable,
				fmt.Sprintf("Use of undefined variable: %s", types.FormatAnyValue(groupKey)),
				"$group (stage)",
			)
		default:
			panic(fmt.Sprintf("unhandled field path error %s", exprErr.Error()))
		}
	}

	var group groupMap

	for _, doc := range in {
		val := expression.Evaluate(doc)
		group.addOrAppend(val, doc)
	}

	return group.docs, nil
}

// groupedDocuments contains group key and the documents for that group.
type groupedDocuments struct {
	groupID   any
	documents []*types.Document
}

// groupMap holds groups of documents.
type groupMap struct {
	docs []groupedDocuments
}

// addOrAppend adds a groupID documents pair if the groupID does not exist,
// if the groupID exists it appends the documents to the slice.
func (m *groupMap) addOrAppend(groupKey any, docs ...*types.Document) {
	for i, g := range m.docs {
		// groupID is a distinct key and can be any BSON type including array and Binary,
		// so we cannot use structure like map.
		// Compare is used to check if groupID exists in groupMap, because
		// numbers are grouped for the same value regardless of their number type.
		if types.CompareForAggregation(groupKey, g.groupID) == types.Equal {
			m.docs[i].documents = append(m.docs[i].documents, docs...)
			return
		}
	}

	m.docs = append(m.docs, groupedDocuments{
		groupID:   groupKey,
		documents: docs,
	})
}

// validateDocumentExpression returns error when there is unsupported expression present.
func validateDocumentExpression(doc *types.Document) error {
	iter := doc.Iterator()
	defer iter.Close()

	for {
		k, v, err := iter.Next()
		if errors.Is(err, iterator.ErrIteratorDone) {
			return nil
		}

		if err != nil {
			return err
		}

		if strings.HasPrefix(k, "$") {
			// TODO: https://github.com/FerretDB/FerretDB/issues/2165
			return commonerrors.NewCommandErrorMsgWithArgument(
				commonerrors.ErrNotImplemented,
				fmt.Sprintf("%s operator is not implemented for $group key expression yet", k),
				"$group (stage)",
			)
		}

		if docVal, ok := v.(*types.Document); ok {
			if err = validateDocumentExpression(docVal); err != nil {
				return err
			}
		}
	}
}

// Type implements Stage interface.
func (g *group) Type() aggregations.StageType {
	return aggregations.StageTypeDocuments
}

// check interfaces
var (
	_ aggregations.Stage = (*group)(nil)
)
