package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// SaveRoles stores roles in cache for the given userID.
func (s *Store) SaveRoles(ctx context.Context, userID uint64, roles []string, ttl time.Duration) error {
	tr := otel.Tracer(RoleCacheTracerName)
	ctx, span := tr.Start(ctx, SpanNameRoleSave, trace.WithAttributes(
		attribute.String(commonkeys.Operation, OperationRoleSave),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
	))
	defer span.End()

	cacheKey := fmt.Sprintf(RoleKeyFormat, userID)
	data, err := json.Marshal(roles)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToSerializeRoles, AttributeRoleCacheKey, cacheKey, commonkeys.Error, err)
		return err
	}

	if ttl <= 0 {
		ttl = RoleExpirationDefault
	}

	if err := s.cache.Set(ctx, cacheKey, string(data), ttl); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToSaveRolesToCache, AttributeRoleCacheKey, cacheKey, AttributeRoleTTL, ttl.String(), commonkeys.Error, err)
		return err
	}

	span.SetStatus(codes.Ok, RolesSavedSuccessfully)
	return nil
}

// GetRoles retrieves cached roles for the given userID.
// Returns nil slice when roles are not cached.
func (s *Store) GetRoles(ctx context.Context, userID uint64) ([]string, error) {
	tr := otel.Tracer(RoleCacheTracerName)
	ctx, span := tr.Start(ctx, SpanNameRoleGet, trace.WithAttributes(
		attribute.String(commonkeys.Operation, OperationRoleGet),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
	))
	defer span.End()

	cacheKey := fmt.Sprintf(RoleKeyFormat, userID)
	data, err := s.cache.Get(ctx, cacheKey)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToGetRolesFromCache, AttributeRoleCacheKey, cacheKey, commonkeys.Error, err)
		return nil, err
	}

	if data == "" {
		span.SetStatus(codes.Ok, RolesRetrievedSuccessfully)
		return nil, nil
	}

	var roles []string
	if err := json.Unmarshal([]byte(data), &roles); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToDeserializeRoles, AttributeRoleCacheKey, cacheKey, commonkeys.Error, err)
		return nil, err
	}

	span.SetStatus(codes.Ok, RolesRetrievedSuccessfully)
	return roles, nil
}

// InvalidateRoles removes cached roles for the given userID.
func (s *Store) InvalidateRoles(ctx context.Context, userID uint64) error {
	tr := otel.Tracer(RoleCacheTracerName)
	ctx, span := tr.Start(ctx, SpanNameRoleDelete, trace.WithAttributes(
		attribute.String(commonkeys.Operation, OperationRoleDelete),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
	))
	defer span.End()

	cacheKey := fmt.Sprintf(RoleKeyFormat, userID)
	if err := s.cache.Del(ctx, cacheKey); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToDeleteRolesFromCache, AttributeRoleCacheKey, cacheKey, commonkeys.Error, err)
		return err
	}

	span.SetStatus(codes.Ok, RolesDeletedSuccessfully)
	return nil
}
