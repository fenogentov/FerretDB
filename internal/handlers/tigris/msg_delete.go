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

package tigris

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/tigrisdata/tigris-client-go/driver"

	"github.com/FerretDB/FerretDB/internal/handlers/common"
	"github.com/FerretDB/FerretDB/internal/handlers/commonerrors"
	"github.com/FerretDB/FerretDB/internal/handlers/tigris/tigrisdb"
	"github.com/FerretDB/FerretDB/internal/handlers/tigris/tjson"
	"github.com/FerretDB/FerretDB/internal/types"
	"github.com/FerretDB/FerretDB/internal/util/iterator"
	"github.com/FerretDB/FerretDB/internal/util/lazyerrors"
	"github.com/FerretDB/FerretDB/internal/util/must"
	"github.com/FerretDB/FerretDB/internal/wire"
)

// MsgDelete implements HandlerInterface.
func (h *Handler) MsgDelete(ctx context.Context, msg *wire.OpMsg) (*wire.OpMsg, error) {
	dbPool, err := h.DBPool(ctx)
	if err != nil {
		return nil, lazyerrors.Error(err)
	}

	document, err := msg.Document()
	if err != nil {
		return nil, lazyerrors.Error(err)
	}

	params, err := common.GetDeleteParams(document, h.L)
	if err != nil {
		return nil, err
	}

	qp := tigrisdb.QueryParams{
		DB:         params.DB,
		Collection: params.Collection,
	}

	var deleted int32
	var delErrors commonerrors.WriteErrors

	// process every delete filter
	for i, deleteParams := range params.Deletes {
		qp.Filter = deleteParams.Filter

		del, err := execDelete(ctx, &execDeleteParams{
			dbPool,
			&qp,
			h.DisableFilterPushdown,
			deleteParams.Limited,
		})
		if err == nil {
			deleted += del
			continue
		}

		delErrors.Append(err, int32(i))

		if params.Ordered {
			break
		}
	}

	replyDoc := must.NotFail(types.NewDocument(
		"ok", float64(1),
	))

	if delErrors.Len() > 0 {
		replyDoc = delErrors.Document()
	}

	replyDoc.Set("n", deleted)

	var reply wire.OpMsg
	must.NoError(reply.SetSections(wire.OpMsgSection{
		Documents: []*types.Document{replyDoc},
	}))

	return &reply, nil
}

// execDeleteParams contains parameters for execDelete function.
type execDeleteParams struct {
	dbPool                *tigrisdb.TigrisDB
	qp                    *tigrisdb.QueryParams
	disableFilterPushdown bool
	limited               bool
}

// execDelete fetches documents, filter them out and limiting with the given limit value.
// It returns the number of deleted documents or an error.
func execDelete(ctx context.Context, dp *execDeleteParams) (int32, error) {
	var err error

	resDocs := make([]*types.Document, 0, 16)

	var deleted int32

	// filter is used to filter documents on the FerretDB side,
	// qp.Filter is used to filter documents on the Tigris side (query pushdown).
	filter := dp.qp.Filter

	if dp.disableFilterPushdown {
		dp.qp.Filter = nil
	}

	// fetch current items from collection
	iter, err := dp.dbPool.QueryDocuments(ctx, dp.qp)
	if err != nil {
		return 0, err
	}

	defer iter.Close()

	// iterate through every document and delete matching ones
	for {
		var doc *types.Document

		_, doc, err = iter.Next()
		if err != nil {
			if errors.Is(err, iterator.ErrIteratorDone) {
				break
			}

			return 0, lazyerrors.Error(err)
		}

		// fetch current items from collection
		matches, err := common.FilterDocument(doc, filter)
		if err != nil {
			return 0, err
		}

		if !matches {
			continue
		}

		resDocs = append(resDocs, doc)

		// if limit is set, no need to fetch all the documents
		if dp.limited {
			break
		}
	}

	// if no field is matched in a row, go to the next one
	if len(resDocs) == 0 {
		return 0, nil
	}

	res, err := deleteDocuments(ctx, dp.dbPool, dp.qp, resDocs)
	if err != nil {
		return 0, err
	}

	deleted += int32(res)

	return deleted, nil
}

// deleteDocuments deletes documents by _id.
func deleteDocuments(ctx context.Context, dbPool *tigrisdb.TigrisDB, qp *tigrisdb.QueryParams, docs []*types.Document) (int, error) { //nolint:lll // argument list is too long
	ids := make([]map[string]any, len(docs))
	for i, doc := range docs {
		id := must.NotFail(tjson.Marshal(must.NotFail(doc.Get("_id"))))
		ids[i] = map[string]any{"_id": json.RawMessage(id)}
	}

	var f driver.Filter
	switch len(ids) {
	case 0:
		f = driver.Filter(`{}`)
	case 1:
		f = must.NotFail(json.Marshal(ids[0]))
	default:
		f = must.NotFail(json.Marshal(map[string]any{"$or": ids}))
	}

	_, err := dbPool.Driver.UseDatabase(qp.DB).Delete(ctx, qp.Collection, f)
	if err != nil {
		return 0, lazyerrors.Error(err)
	}

	return len(ids), nil
}
