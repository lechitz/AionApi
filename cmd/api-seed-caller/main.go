package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	defaultHost     = "http://localhost:5001"
	defaultContext  = "/aion"
	defaultAPIRoot  = "/api/v1"
	defaultGraphql  = "/graphql"
	defaultUser     = "user1"
	defaultPass     = "testpassword123"
	defaultLogPath  = "infrastructure/db/seed/api_success.log"
	defaultHTTPTO   = 10 * time.Second
	userAgent       = "api-seed-caller/1.0"
	sourceName      = "api-seed-caller"
	statusCompleted = "completed"
	timezoneUTC     = "UTC"
)

type cfg struct {
	host        string
	contextPath string
	apiRoot     string
	graphqlPath string
	username    string
	password    string
	autoCreate  bool
	cleanBefore bool
	cleanOnly   bool
	userPrefix  string
	userCount   int
	runID       string
	successLog  string
	timeout     time.Duration
}

type gqlRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

type gqlResponse struct {
	Data   map[string]json.RawMessage `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

type loginResponse struct {
	Token  string `json:"token"`
	Result struct {
		Token string `json:"token"`
	} `json:"result"`
}

type createdCategory struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type listedCategory struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type createdTag struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	CategoryID string `json:"categoryId"`
}

type createdRecord struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	TagID       string `json:"tagId"`
}

type categoryPayload struct {
	Key         string
	Name        string
	Description string
	ColorHex    string
	Icon        string
}

type tagPayload struct {
	Key         string
	Name        string
	CategoryKey string
	Description string
}

type recordPayload struct {
	Title           string
	Description     string
	CategoryKey     string
	TagKey          string
	DurationSeconds int
	Value           float64
}

//nolint:gocognit,funlen,nestif // main orchestrates seed workflow; refactor left for future
func main() {
	ctx := context.Background()
	config := loadConfig()
	client := &http.Client{Timeout: config.timeout}

	apiBase := buildURL(config.host, config.contextPath, config.apiRoot)
	loginURL := strings.TrimSuffix(apiBase, "/") + "/auth/login"
	createUserURL := strings.TrimSuffix(apiBase, "/") + "/user/create"
	deleteUserURL := strings.TrimSuffix(apiBase, "/") + "/user"
	graphqlURL := buildURL(config.host, config.contextPath, config.apiRoot, config.graphqlPath)

	usernames := buildUsernames(config)

	_, _ = fmt.Fprintf(
		os.Stdout,
		"▶️  run_id=%s host=%s users=%d strict=%v clean=%v\n",
		config.runID,
		config.host,
		len(usernames),
		!config.autoCreate,
		config.cleanBefore,
	)

	for idx, baseUsername := range usernames {
		userRunID := buildUserRunID(config.runID, baseUsername)
		_, _ = fmt.Fprintf(os.Stdout, "\n👤 [%d/%d] user=%s run_id=%s\n", idx+1, len(usernames), baseUsername, userRunID)

		token, username, err := loginOrCreateUser(ctx, client, loginURL, createUserURL, config, baseUsername, userRunID)
		if err != nil {
			exitWithErrf("%v", err)
		}
		logSuccess(config.successLog, fmt.Sprintf("login success user=%s run_id=%s", username, userRunID))

		headers := map[string]string{}
		if userRunID != "" {
			headers["X-Seed-Run-Id"] = userRunID
		}
		if debugOn() {
			_, _ = fmt.Fprintf(os.Stdout, "debug: token_len=%d run_id=%s user=%s\n", len(token), userRunID, username)
		}

		categoryInputs := []categoryPayload{
			{Key: "saude_fisica", Name: "saude_fisica", Description: "Atividades físicas e condicionamento", ColorHex: "#E94F37", Icon: "🏋️"},
			{Key: "saude_mental", Name: "saude_mental", Description: "Saúde mental e reflexão", ColorHex: "#F8B400", Icon: "🧠"},
			{Key: "pessoal", Name: "pessoal", Description: "Momentos pessoais e lazer", ColorHex: "#FF6F00", Icon: "🎉"},
		}

		tagInputs := []tagPayload{
			{Key: "alongamento", Name: "Stretching", CategoryKey: "saude_fisica", Description: "Stretching and flexibility"},
			{Key: "corrida", Name: "Run", CategoryKey: "saude_fisica", Description: "Running sessions"},
			{Key: "diario", Name: "Diary", CategoryKey: "saude_mental", Description: "Personal diary entries"},
			{Key: "leitura", Name: "Reading", CategoryKey: "saude_mental", Description: "Reading and books"},
			{Key: "filme", Name: "Movie", CategoryKey: "pessoal", Description: "Watching movies"},
			{Key: "viagem", Name: "Travel", CategoryKey: "pessoal", Description: "Travel activities"},
		}

		recordInputs := []recordPayload{
			{
				Title:           "Alongamento matinal",
				Description:     "Rotina curta de alongamento",
				CategoryKey:     "saude_fisica",
				TagKey:          "alongamento",
				DurationSeconds: 900,
				Value:           15,
			},
			{Title: "Corrida leve", Description: "Corrida de 3km no parque", CategoryKey: "saude_fisica", TagKey: "corrida", DurationSeconds: 1800, Value: 3},
			{Title: "Diário", Description: "Notas rápidas do dia", CategoryKey: "saude_mental", TagKey: "diario", DurationSeconds: 600, Value: 1},
			{Title: "Leitura", Description: "Capítulo de livro técnico", CategoryKey: "saude_mental", TagKey: "leitura", DurationSeconds: 1800, Value: 20},
			{Title: "Filme", Description: "Sessão de cinema em casa", CategoryKey: "pessoal", TagKey: "filme", DurationSeconds: 7200, Value: 1},
			{Title: "Planejar viagem", Description: "Revisar roteiros e reservas", CategoryKey: "pessoal", TagKey: "viagem", DurationSeconds: 2700, Value: 1},
		}

		if config.cleanBefore {
			_, _ = fmt.Fprintln(os.Stdout, "🧹 Limpando registros anteriores via GraphQL...")
			if err := softDeleteAllRecords(ctx, client, graphqlURL, token, headers); err != nil {
				exitWithErrf("clean failed: %v", err)
			}
			logSuccess(config.successLog, fmt.Sprintf("clean records success run_id=%s", userRunID))

			if config.cleanOnly {
				_, _ = fmt.Fprintln(os.Stdout, "🧹 Limpando tags...")
				if err := softDeleteAllTags(ctx, client, graphqlURL, token, headers); err != nil {
					exitWithErrf("clean tags failed: %v", err)
				}
				logSuccess(config.successLog, fmt.Sprintf("clean tags success run_id=%s", userRunID))

				_, _ = fmt.Fprintln(os.Stdout, "🧹 Limpando categorias...")
				if err := softDeleteAllCategories(ctx, client, graphqlURL, token, headers); err != nil {
					exitWithErrf("clean categories failed: %v", err)
				}
				logSuccess(config.successLog, fmt.Sprintf("clean categories success run_id=%s", userRunID))

				_, _ = fmt.Fprintln(os.Stdout, "🧹 Soft deleting user...")
				if err := softDeleteUser(ctx, client, deleteUserURL, token, headers); err != nil {
					exitWithErrf("clean user failed: %v", err)
				}
				logSuccess(config.successLog, fmt.Sprintf("clean user success user=%s run_id=%s", username, userRunID))
				_, _ = fmt.Fprintf(os.Stdout, "✅ Clean-only finalizado para usuário %s\n", username)
				continue
			}
		}

		categories := make(map[string]createdCategory)
		for _, cat := range categoryInputs {
			name := appendRunID(cat.Name, userRunID)
			created, err := ensureCategory(ctx, client, graphqlURL, token, headers, cat, name)
			if err != nil {
				exitWithErrf("category failed for %s: %v", cat.Key, err)
			}
			categories[cat.Key] = created
			logSuccess(config.successLog, fmt.Sprintf("category key=%s id=%s name=%s run_id=%s", cat.Key, created.ID, created.Name, userRunID))
			_, _ = fmt.Fprintf(os.Stdout, "✅ category: %s\n", created.Name)
		}

		tags := make(map[string]createdTag)
		for _, tg := range tagInputs {
			category, ok := categories[tg.CategoryKey]
			if !ok {
				exitWithErrf("category key not found for tag %s", tg.CategoryKey)
			}
			name := appendRunID(tg.Name, userRunID)
			created, err := ensureTag(ctx, client, graphqlURL, token, headers, tg, name, category.ID)
			if err != nil {
				exitWithErrf("tag failed for %s: %v", tg.Key, err)
			}
			tags[tg.Key] = created
			logSuccess(
				config.successLog,
				fmt.Sprintf("tag key=%s id=%s name=%s categoryId=%s run_id=%s", tg.Key, created.ID, created.Name, created.CategoryID, userRunID),
			)
			_, _ = fmt.Fprintf(os.Stdout, "✅ tag: %s\n", created.Name)
		}

		for i, rc := range recordInputs {
			_, ok := categories[rc.CategoryKey]
			if !ok {
				exitWithErrf("category key not found for record %s", rc.CategoryKey)
			}

			tag, ok := tags[rc.TagKey]
			if !ok {
				exitWithErrf("tag key not found for record %s", rc.TagKey)
			}

			eventTime := time.Now().UTC().Add(-time.Duration(i) * time.Hour).Format(time.RFC3339)
			title := fmt.Sprintf("%s %d", rc.Title, i+1)
			created, err := createRecord(ctx, client, graphqlURL, token, headers, rc, title, tag.ID, eventTime)
			if err != nil {
				exitWithErrf("createRecord failed for %s: %v", rc.Title, err)
			}
			if id := strings.TrimSpace(created.ID); id == "" || id == "0" {
				exitWithErrf("createRecord returned empty ID for %s", title)
			}
			logSuccess(config.successLog, fmt.Sprintf("record id=%s description=%s tagId=%s run_id=%s", created.ID, created.Description, created.TagID, userRunID))
			_, _ = fmt.Fprintf(os.Stdout, "✅ record: %s\n", created.Description)
		}

		if err := verifyArtifacts(ctx, client, graphqlURL, token, headers, categoryInputs, tagInputs, recordInputs, userRunID); err != nil {
			exitWithErrf("verificação falhou: %v", err)
		}

		logSuccess(config.successLog, fmt.Sprintf("user seed completed user=%s run_id=%s", username, userRunID))
	}

	_, _ = fmt.Fprintf(os.Stdout, "📄 Success log written to %s\n", config.successLog)
}

func loadConfig() cfg {
	return cfg{
		host:        getenv("API_CALLER_HOST", defaultHost),
		contextPath: getenv("API_CALLER_CONTEXT", defaultContext),
		apiRoot:     getenv("API_CALLER_ROOT", defaultAPIRoot),
		graphqlPath: getenv("API_CALLER_GRAPHQL", defaultGraphql),
		username:    getenv("API_CALLER_USER", defaultUser),
		password:    resolvePassword(),
		autoCreate:  parseBool(getenv("API_CALLER_AUTO_CREATE", "")),
		cleanBefore: parseBool(getenv("API_CALLER_CLEAN", "")),
		cleanOnly:   parseBool(getenv("API_CALLER_ONLY_CLEAN", "")),
		userPrefix:  getenv("API_CALLER_USER_PREFIX", "user"),
		userCount:   parseInt(getenv("API_CALLER_COUNT", ""), 1),
		runID:       getenvAllowEmpty("API_CALLER_RUN_ID", ""),
		successLog:  getenv("API_CALLER_SUCCESS_LOG", defaultLogPath),
		timeout:     parseDuration(getenv("API_CALLER_TIMEOUT", "")),
	}
}

func parseDuration(raw string) time.Duration {
	if raw == "" {
		return defaultHTTPTO
	}
	d, err := time.ParseDuration(raw)
	if err != nil {
		return defaultHTTPTO
	}
	return d
}

func parseBool(raw string) bool {
	raw = strings.TrimSpace(strings.ToLower(raw))
	return raw == "1" || raw == "true" || raw == "yes" || raw == "y"
}

func parseInt(raw string, fallback int) int {
	if raw == "" {
		return fallback
	}

	n, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil || n <= 0 {
		return fallback
	}
	return n
}

func getenv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getenvAllowEmpty(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

func resolvePassword() string {
	if val := os.Getenv("API_CALLER_PASSWORD"); val != "" {
		return val
	}
	return getenv("API_CALLER_PASS", defaultPass)
}

func buildUsernames(config cfg) []string {
	if config.userCount <= 1 {
		return []string{config.username}
	}

	users := make([]string, 0, config.userCount)
	for i := 1; i <= config.userCount; i++ {
		users = append(users, fmt.Sprintf("%s%d", config.userPrefix, i))
	}
	return users
}

func buildUserRunID(runID, baseUser string) string {
	if strings.TrimSpace(runID) == "" {
		return ""
	}
	return fmt.Sprintf("%s_%s", runID, baseUser)
}

func buildURL(parts ...string) string {
	if len(parts) == 0 {
		return ""
	}

	out := parts[0]
	for _, p := range parts[1:] {
		if p == "" {
			continue
		}
		if !strings.HasPrefix(p, "/") {
			p = "/" + p
		}
		out = strings.TrimSuffix(out, "/") + p
	}
	return out
}

func appendRunID(base, runID string) string {
	if strings.TrimSpace(runID) == "" {
		return base
	}
	return fmt.Sprintf("%s_%s", base, runID)
}

func loginOrCreateUser(ctx context.Context, client *http.Client, loginURL, createUserURL string, config cfg, baseUsername, runID string) (string, string, error) {
	username := baseUsername
	var email string

	_, _ = fmt.Fprintf(os.Stdout, "🔑 Logging in via %s\n", loginURL)
	token, err := login(ctx, client, loginURL, username, config.password)
	if err == nil {
		return token, username, nil
	}

	if !config.autoCreate {
		return "", "", fmt.Errorf("login failed for %s: %w (garanta credenciais válidas ou exporte API_CALLER_AUTO_CREATE=true)", username, err)
	}

	_, _ = fmt.Fprintf(os.Stdout, "⚠️  login falhou para %s, criando usuário e tentando novamente: %v\n", username, err)
	username = appendRunID(baseUsername, runID)
	email = fmt.Sprintf("%s@aionapi.dev", username)
	if err := createUser(ctx, client, createUserURL, username, email, config.password, runID); err != nil {
		return "", "", fmt.Errorf("user create failed for %s: %w", username, err)
	}
	logSuccess(config.successLog, fmt.Sprintf("user created username=%s email=%s run_id=%s", username, email, runID))

	token, err = login(ctx, client, loginURL, username, config.password)
	if err != nil {
		return "", "", fmt.Errorf("login failed after create for %s: %w", username, err)
	}

	return token, username, nil
}

func login(ctx context.Context, client *http.Client, url, user, pass string) (string, error) {
	payload := map[string]string{
		"username": user,
		"password": pass,
	}
	body, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	var parsed loginResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return "", err
	}

	if parsed.Token == "" {
		parsed.Token = parsed.Result.Token
	}

	if parsed.Token == "" {
		return "", errors.New("empty token in response")
	}

	return parsed.Token, nil
}

func createUser(ctx context.Context, client *http.Client, url, username, email, password, runID string) error {
	payload := map[string]string{
		"name":     fmt.Sprintf("%s_name", username),
		"username": username,
		"email":    email,
		"password": password,
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("X-Seed-Run-Id", runID)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusCreated {
		return nil
	}

	bodyBytes, _ := io.ReadAll(resp.Body)
	return fmt.Errorf("user create status %d body=%s", resp.StatusCode, string(bodyBytes))
}

func ensureCategory(ctx context.Context, client *http.Client, url, token string, headers map[string]string, cat categoryPayload, name string) (createdCategory, error) {
	if existing, err := getCategoryByName(ctx, client, url, token, headers, name); err == nil {
		if id := strings.TrimSpace(existing.ID); id != "" && id != "0" {
			return existing, nil
		}
	}
	created, err := createCategory(ctx, client, url, token, headers, cat, name)
	if err != nil {
		return createdCategory{}, err
	}
	if id := strings.TrimSpace(created.ID); id == "" || id == "0" {
		return createdCategory{}, fmt.Errorf("createCategory retornou ID vazio para %s", name)
	}
	return created, nil
}

func ensureTag(ctx context.Context, client *http.Client, url, token string, headers map[string]string, tg tagPayload, name, categoryID string) (createdTag, error) {
	if existing, err := getTagByName(ctx, client, url, token, headers, name); err == nil {
		if id := strings.TrimSpace(existing.ID); id != "" && id != "0" {
			return existing, nil
		}
	}
	created, err := createTag(ctx, client, url, token, headers, tg, name, categoryID)
	if err != nil {
		return createdTag{}, err
	}
	id := strings.TrimSpace(created.ID)
	cid := strings.TrimSpace(created.CategoryID)
	if id == "" || id == "0" || cid == "" || cid == "0" {
		return createdTag{}, fmt.Errorf("createTag retornou ID vazio para %s", name)
	}
	return created, nil
}

func getCategoryByName(ctx context.Context, client *http.Client, url, token string, headers map[string]string, name string) (createdCategory, error) {
	const query = `
query($name: String!) {
  categoryByName(name: $name) { id name }
}`
	var out createdCategory
	if err := doGraphQL(ctx, client, url, token, headers, query, map[string]interface{}{"name": name}, &out); err != nil {
		return createdCategory{}, err
	}
	return out, nil
}

func getTagByName(ctx context.Context, client *http.Client, url, token string, headers map[string]string, name string) (createdTag, error) {
	const query = `
query($name: String!) {
  tagByName(name: $name) { id name categoryId }
}`
	var out createdTag
	if err := doGraphQL(ctx, client, url, token, headers, query, map[string]interface{}{"name": name}, &out); err != nil {
		return createdTag{}, err
	}
	return out, nil
}

func createCategory(
	ctx context.Context,
	client *http.Client,
	url, token string,
	headers map[string]string,
	cat categoryPayload,
	finalName string,
) (createdCategory, error) {
	const mutation = `
mutation($input: CreateCategoryInput!) {
  createCategory(input: $input) { id name }
}`

	input := map[string]interface{}{
		"name":        finalName,
		"description": cat.Description,
		"colorHex":    cat.ColorHex,
		"icon":        cat.Icon,
	}

	var out createdCategory
	if err := doGraphQL(ctx, client, url, token, headers, mutation, map[string]interface{}{"input": input}, &out); err != nil {
		return createdCategory{}, err
	}
	return out, nil
}

func createTag(ctx context.Context, client *http.Client, url, token string, headers map[string]string, tg tagPayload, finalName, categoryID string) (createdTag, error) {
	const mutation = `
mutation($input: CreateTagInput!) {
  createTag(input: $input) { id name categoryId }
}`

	input := map[string]interface{}{
		"name":        finalName,
		"categoryId":  categoryID,
		"description": tg.Description,
	}

	var out createdTag
	if err := doGraphQL(ctx, client, url, token, headers, mutation, map[string]interface{}{"input": input}, &out); err != nil {
		return createdTag{}, err
	}
	return out, nil
}

func createRecord(
	ctx context.Context,
	client *http.Client,
	url, token string,
	headers map[string]string,
	rc recordPayload,
	title, tagID, eventTime string,
) (createdRecord, error) {
	const mutation = `
mutation($input: CreateRecordInput!) {
  createRecord(input: $input) { id description tagId }
}`

	// Note: title is prepended to description since CreateRecordInput doesn't have a title field
	// categoryId is not needed since tag already belongs to a category
	description := fmt.Sprintf("%s - %s", title, rc.Description)
	input := map[string]interface{}{
		"tagId":           tagID,
		"description":     description,
		"eventTime":       eventTime,
		"recordedAt":      eventTime,
		"durationSeconds": rc.DurationSeconds,
		"value":           rc.Value,
		"source":          sourceName,
		"timezone":        timezoneUTC,
		"status":          statusCompleted,
	}

	var out createdRecord
	if err := doGraphQL(ctx, client, url, token, headers, mutation, map[string]interface{}{"input": input}, &out); err != nil {
		return createdRecord{}, err
	}
	return out, nil
}

func softDeleteAllRecords(ctx context.Context, client *http.Client, url, token string, headers map[string]string) error {
	const mutation = `
mutation {
  softDeleteAllRecords
}`
	var out bool
	return doGraphQL(ctx, client, url, token, headers, mutation, nil, &out)
}

func softDeleteAllCategories(ctx context.Context, client *http.Client, url, token string, headers map[string]string) error {
	cats, err := listCategories(ctx, client, url, token, headers)
	if err != nil {
		return err
	}

	for _, cat := range cats {
		if err := softDeleteCategory(ctx, client, url, token, headers, cat.ID); err != nil {
			return fmt.Errorf("soft delete category %s (%s): %w", cat.Name, cat.ID, err)
		}
	}
	return nil
}

func softDeleteCategory(ctx context.Context, client *http.Client, url, token string, headers map[string]string, id string) error {
	const mutation = `
mutation($input: DeleteCategoryInput!) {
  softDeleteCategory(input: $input)
}`
	var out bool
	input := map[string]interface{}{"id": id}
	return doGraphQL(ctx, client, url, token, headers, mutation, map[string]interface{}{"input": input}, &out)
}

func listCategories(ctx context.Context, client *http.Client, url, token string, headers map[string]string) ([]listedCategory, error) {
	const query = `
query {
  categories { id name }
}`
	var out []listedCategory
	if err := doGraphQL(ctx, client, url, token, headers, query, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func softDeleteAllTags(ctx context.Context, client *http.Client, url, token string, headers map[string]string) error {
	tags, err := listAllTags(ctx, client, url, token, headers)
	if err != nil {
		return err
	}

	for _, tg := range tags {
		if err := softDeleteTag(ctx, client, url, token, headers, tg.ID); err != nil {
			return fmt.Errorf("soft delete tag %s (%s): %w", tg.Name, tg.ID, err)
		}
	}
	return nil
}

func softDeleteTag(ctx context.Context, client *http.Client, url, token string, headers map[string]string, id string) error {
	const mutation = `
mutation($input: DeleteTagInput!) {
  softDeleteTag(input: $input)
}`
	var out bool
	input := map[string]interface{}{"input": map[string]interface{}{"id": id}}
	return doGraphQL(ctx, client, url, token, headers, mutation, input, &out)
}

func verifyArtifacts(
	ctx context.Context,
	client *http.Client,
	url, token string,
	headers map[string]string,
	cats []categoryPayload,
	tags []tagPayload,
	records []recordPayload,
	runID string,
) error {
	// Verify categories
	for _, cat := range cats {
		name := appendRunID(cat.Name, runID)
		got, err := getCategoryByName(ctx, client, url, token, headers, name)
		if err != nil {
			return fmt.Errorf("categoria não encontrada pós-seed: %s: %w", name, err)
		}
		if got.ID == "" {
			return fmt.Errorf("categoria não encontrada pós-seed: %s (empty id)", name)
		}
	}

	// Verify tags
	for _, tg := range tags {
		name := appendRunID(tg.Name, runID)
		got, err := getTagByName(ctx, client, url, token, headers, name)
		if err != nil {
			return fmt.Errorf("tag não encontrada pós-seed: %s: %w", name, err)
		}
		if got.ID == "" {
			return fmt.Errorf("tag não encontrada pós-seed: %s (empty id)", name)
		}
	}

	// Verify records presence by description fragment
	expectedDescriptions := make(map[string]struct{})
	for i, rc := range records {
		title := fmt.Sprintf("%s %d", rc.Title, i+1)
		description := fmt.Sprintf("%s - %s", title, rc.Description)
		expectedDescriptions[description] = struct{}{}
	}

	found, err := listRecords(ctx, client, url, token, headers, len(records)*3)
	if err != nil {
		return fmt.Errorf("erro ao listar records: %w", err)
	}

	for _, rec := range found {
		delete(expectedDescriptions, rec.Description)
	}

	if len(expectedDescriptions) > 0 {
		missing := make([]string, 0, len(expectedDescriptions))
		for t := range expectedDescriptions {
			missing = append(missing, t)
		}
		return fmt.Errorf("records ausentes: %s", strings.Join(missing, ", "))
	}

	return nil
}

func listRecords(ctx context.Context, client *http.Client, url, token string, headers map[string]string, limit int) ([]createdRecord, error) {
	const query = `
query($limit: Int) {
  records(limit: $limit) { id description tagId }
}`
	var out []createdRecord
	if err := doGraphQL(ctx, client, url, token, headers, query, map[string]interface{}{"limit": limit}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func doGraphQL(
	ctx context.Context,
	client *http.Client,
	url, token string,
	headers map[string]string,
	query string,
	variables map[string]interface{},
	target interface{},
) error {
	payload := gqlRequest{
		Query:     query,
		Variables: variables,
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Authorization", "Bearer "+token)
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("graphql status %d", resp.StatusCode)
	}

	var parsed gqlResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return err
	}

	if len(parsed.Errors) > 0 {
		if debugOn() {
			_, _ = fmt.Fprintf(os.Stdout, "debug: graphql errors: %+v\n", parsed.Errors)
		}
		return fmt.Errorf("graphql error: %s", parsed.Errors[0].Message)
	}

	if len(parsed.Data) != 1 {
		return fmt.Errorf("unexpected data payload size: %d", len(parsed.Data))
	}

	for _, raw := range parsed.Data {
		if bytes.Equal(raw, []byte("null")) {
			return errors.New("graphql returned null payload")
		}
		if debugOn() {
			_, _ = fmt.Fprintf(os.Stdout, "debug: raw graphql data: %s\n", string(raw))
		}
		if err := json.Unmarshal(raw, target); err != nil {
			return err
		}
		break
	}

	return nil
}

func listAllTags(ctx context.Context, client *http.Client, url, token string, headers map[string]string) ([]createdTag, error) {
	const query = `
query {
  tags { id name categoryId }
}`
	var out []createdTag
	if err := doGraphQL(ctx, client, url, token, headers, query, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func softDeleteUser(ctx context.Context, client *http.Client, url, token string, headers map[string]string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("user soft delete status %d body=%s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	return nil
}

func logSuccess(path, msg string) {
	if path == "" {
		return
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "warn: cannot create log dir: %v\n", err)
		return
	}

	f, err := os.OpenFile(filepath.Clean(path), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "warn: cannot open log file: %v\n", err)
		return
	}
	defer func() {
		_ = f.Close()
	}()

	timestamp := time.Now().UTC().Format(time.RFC3339)
	_, _ = fmt.Fprintf(f, "%s %s\n", timestamp, msg)
}

func exitWithErrf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, "❌ "+format+"\n", args...)
	os.Exit(1)
}

func debugOn() bool {
	return parseBool(os.Getenv("API_CALLER_DEBUG"))
}
