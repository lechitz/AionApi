package tags

import (
	"strconv"

	"github.com/lechitz/AionApi/internal/adapters/primary/graph/model"
	"github.com/lechitz/AionApi/internal/core/domain"
)

// toDomainCreate converte o input GraphQL em entidade de domínio.
// Observações:
// - ID de Category vem como string (GraphQL ID) e vira uint64 no domínio.
// - Name é obrigatório no schema (String!), então aqui não precisa validar vazio.
func toDomainCreate(in model.CreateTagInput) (domain.Tag, error) {
	cid, err := strconv.ParseUint(in.CategoryID, 10, 64)
	if err != nil {
		return domain.Tag{}, err
	}
	t := domain.Tag{
		Name:       in.Name,
		CategoryID: cid,
	}
	if in.Description != nil {
		t.Description = *in.Description
	}
	return t, nil
}

// toModelOut mapeia a entidade de domínio para o modelo GraphQL (saída).
func toModelOut(t domain.Tag) *model.Tag {
	return &model.Tag{
		ID:          strconv.FormatUint(t.ID, 10),
		Name:        t.Name,
		CategoryID:  strconv.FormatUint(t.CategoryID, 10),
		Description: strPtr(t.Description),
	}
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
