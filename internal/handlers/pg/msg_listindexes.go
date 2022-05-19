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

package pg

import (
	"context"
	"errors"
	"fmt"

	"github.com/FerretDB/FerretDB/internal/handlers/common"
	"github.com/FerretDB/FerretDB/internal/util/lazyerrors"
	"github.com/FerretDB/FerretDB/internal/wire"
	"github.com/jackc/pgx/v4"
)

// MsgListIndexes OpMsg used to get parameter.
func (h *Handler) MsgListIndexes(ctx context.Context, msg *wire.OpMsg) (*wire.OpMsg, error) {
	document, err := msg.Document()
	if err != nil {
		return nil, lazyerrors.Error(err)
	}

	command := document.Command()

	var db string
	if db, err = common.GetRequiredParam[string](document, "$db"); err != nil {
		return nil, err
	}

	names, err := h.pgPool.Tables(ctx, db)
	if err != nil {
		return nil, lazyerrors.Error(err)
	}

	var collection string
	if collection, err = common.GetRequiredParam[string](document, command); err != nil {
		return nil, err
	}

	if !contains(names, collection) {
		return nil, errors.New("no collection")
	}

	fmt.Println("&db >", db)
	fmt.Println("&collection >", db)

	sql := `SELECT FROM ` + pgx.Identifier{db, collection}.Sanitize()

	fmt.Println(sql)
	// names, err := h.pgPool.Tables(ctx, db)
	// if err != nil {
	// 	return nil, lazyerrors.Error(err)
	// }

	// fmt.Println(db, names)

	// sql := `SELECT * FROM pg_indexes WHERE tablename = 'actor'`
	// SELECT indexname, indexdef FROM pg_indexes WHERE tablename = 'your_table';
	// sql := `SELECT `

	// sql += `_jsonb FROM ` + pgx.Identifier{param.db, param.collection}.Sanitize()

	// rows, err := h.pgPool.Query(ctx, sql)
	// if err != nil {
	// 	fmt.Println(">>>", err)
	// 	return nil, lazyerrors.Error(err)
	// }
	// defer rows.Close()

	// fmt.Println(collection)

	// for {
	// 	fmt.Println("$")
	// 	if !rows.Next() {
	// 		if err := rows.Err(); err != nil {
	// 			fmt.Println(">", err)
	// 			return nil, lazyerrors.Error(err)
	// 		}
	// 		fmt.Println("***")
	// 		return nil, io.EOF
	// 	}

	// 	var s string
	// 	if err := rows.Scan(&s); err != nil {
	// 		fmt.Println(">>", err)
	// 		return nil, lazyerrors.Error(err)
	// 	}

	// 	fmt.Println(s)
	// }

	//common.UnimplementedNonDefault

	return nil, nil
}

func contains(slice []string, c string) bool {
	for _, s := range slice {
		if s == c {
			return true
		}
	}
	return false
}
