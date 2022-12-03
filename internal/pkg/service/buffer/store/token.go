package store

import (
	"context"
	"sort"

	etcd "go.etcd.io/etcd/client/v3"

	"github.com/keboola/keboola-as-code/internal/pkg/service/buffer/store/key"
	"github.com/keboola/keboola-as-code/internal/pkg/service/buffer/store/model"
	serviceError "github.com/keboola/keboola-as-code/internal/pkg/service/common/errors"
	"github.com/keboola/keboola-as-code/internal/pkg/service/common/etcdop/iterator"
	"github.com/keboola/keboola-as-code/internal/pkg/service/common/etcdop/op"
	"github.com/keboola/keboola-as-code/internal/pkg/telemetry"
)

func (s *Store) ListTokens(ctx context.Context, receiverKey key.ReceiverKey) (out []model.TokenForExport, err error) {
	_, span := s.tracer.Start(ctx, "keboola.go.buffer.configstore.ListTokens")
	defer telemetry.EndSpan(span, &err)

	tokens, err := s.getReceiverTokensOp(ctx, receiverKey).Do(ctx, s.client).All()
	if err != nil {
		return nil, err
	}
	sort.SliceStable(tokens, func(i, j int) bool {
		return string(tokens[i].KV.Key) < string(tokens[j].KV.Key)
	})

	return tokens.Values(), nil
}

func (s *Store) UpdateTokens(ctx context.Context, tokens []model.TokenForExport) (err error) {
	_, span := s.tracer.Start(ctx, "keboola.go.buffer.configstore.UpdateTokens")
	defer telemetry.EndSpan(span, &err)

	ops := make([]op.Op, 0, len(tokens))
	for _, token := range tokens {
		ops = append(ops, s.updateTokenOp(ctx, token))
	}

	_, err = op.MergeToTxn(ops...).Do(ctx, s.client)

	return err
}

func (s *Store) getReceiverTokensOp(_ context.Context, receiverKey key.ReceiverKey) iterator.Definition[model.TokenForExport] {
	return s.schema.
		Secrets().
		Tokens().
		InReceiver(receiverKey).
		GetAll()
}

func (s *Store) updateTokenOp(_ context.Context, token model.TokenForExport) op.NoResultOp {
	return s.schema.
		Secrets().
		Tokens().
		InExport(token.ExportKey).
		Put(token)
}

func (s *Store) createTokenOp(_ context.Context, token model.TokenForExport) op.BoolOp {
	return s.schema.
		Secrets().
		Tokens().
		InExport(token.ExportKey).
		PutIfNotExists(token).
		WithProcessor(func(_ context.Context, _ etcd.OpResponse, ok bool, err error) (bool, error) {
			if !ok && err == nil {
				return false, serviceError.NewResourceAlreadyExistsError("token", token.ExportKey.String(), "export")
			}
			return ok, err
		})
}

func (s *Store) getTokenOp(_ context.Context, exportKey key.ExportKey) op.ForType[*op.KeyValueT[model.TokenForExport]] {
	return s.schema.
		Secrets().
		Tokens().
		InExport(exportKey).
		Get().
		WithProcessor(func(_ context.Context, _ etcd.OpResponse, kv *op.KeyValueT[model.TokenForExport], err error) (*op.KeyValueT[model.TokenForExport], error) {
			if kv == nil && err == nil {
				return nil, serviceError.NewResourceNotFoundError("token", exportKey.String())
			}
			return kv, err
		})
}

func (s *Store) deleteExportTokenOp(_ context.Context, exportKey key.ExportKey) op.BoolOp {
	return s.schema.
		Secrets().
		Tokens().
		InExport(exportKey).
		Delete().
		WithProcessor(func(_ context.Context, _ etcd.OpResponse, result bool, err error) (bool, error) {
			if !result && err == nil {
				return false, serviceError.NewResourceNotFoundError("token", exportKey.String())
			}
			return result, err
		})
}

func (s *Store) deleteReceiverTokensOp(_ context.Context, receiverKey key.ReceiverKey) op.CountOp {
	return s.schema.
		Secrets().
		Tokens().
		InReceiver(receiverKey).
		DeleteAll()
}
