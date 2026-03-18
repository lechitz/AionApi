# GraphQL Playground

Interactive GraphQL workspace with preloaded operations from the current schema.
In published docs, Demo Mode is enabled by default so readers can run examples without backend access.

This page is a convenience surface. If it drifts, `docs/graphql/schema.graphql` and `contracts/graphql/` remain canonical.

## Playground

<div style="border:1px solid #cfd8e3;border-radius:12px;overflow:hidden;height:82vh;background:#fff;">
  <iframe
    src="../tools/graphql-playground/index.html"
    title="AionApi GraphQL Playground"
    style="width:100%;height:100%;border:0;"
    loading="lazy">
  </iframe>
</div>

## Available Operations

<div class="ops-grid">
  <section class="ops-card">
    <h3>Queries</h3>
    <ul>
      <li><code>categories</code></li>
      <li><code>categoryById</code></li>
      <li><code>categoryByName</code></li>
      <li><code>chatHistory</code></li>
      <li><code>chatContext</code></li>
      <li><code>chatDataPack</code></li>
      <li><code>recordById</code></li>
      <li><code>records</code></li>
      <li><code>recordsLatest</code></li>
      <li><code>recordProjectionById</code></li>
      <li><code>recordProjections</code></li>
      <li><code>recordProjectionsLatest</code></li>
      <li><code>recordsByTag</code></li>
      <li><code>recordsByCategory</code></li>
      <li><code>recordsByDay</code></li>
      <li><code>recordsUntil</code></li>
      <li><code>recordsBetween</code></li>
      <li><code>searchRecords</code></li>
      <li><code>recordStats</code></li>
      <li><code>dashboardSnapshot</code></li>
      <li><code>insightFeed</code></li>
      <li><code>analyticsSeries</code></li>
      <li><code>metricDefinitions</code></li>
      <li><code>dashboardViews</code></li>
      <li><code>dashboardView</code></li>
      <li><code>dashboardWidgetCatalog</code></li>
      <li><code>suggestMetricDefinitions</code></li>
      <li><code>tagByName</code></li>
      <li><code>tagById</code></li>
      <li><code>tags</code></li>
      <li><code>tagsByCategoryId</code></li>
      <li><code>userStats</code></li>
    </ul>
  </section>

  <section class="ops-card">
    <h3>Mutations</h3>
    <ul>
      <li><code>createCategory</code></li>
      <li><code>updateCategory</code></li>
      <li><code>softDeleteCategory</code></li>
      <li><code>createRecord</code></li>
      <li><code>updateRecord</code></li>
      <li><code>softDeleteRecord</code></li>
      <li><code>softDeleteAllRecords</code></li>
      <li><code>createTag</code></li>
      <li><code>updateTag</code></li>
      <li><code>softDeleteTag</code></li>
      <li><code>upsertMetricDefinition</code></li>
      <li><code>upsertGoalTemplate</code></li>
      <li><code>deleteGoalTemplate</code></li>
      <li><code>createDashboardView</code></li>
      <li><code>setDefaultDashboardView</code></li>
      <li><code>upsertDashboardWidget</code></li>
      <li><code>reorderDashboardWidgets</code></li>
      <li><code>deleteDashboardWidget</code></li>
      <li><code>createMetricAndWidget</code></li>
    </ul>
  </section>
</div>

## Notes

- Demo Mode returns mock responses for each operation (default enabled).
- Disable Demo Mode to execute real requests against your own endpoint (local/staging) and optional token.
- Add `Authorization: Bearer <token>` in the token field for protected operations.
- Keep templates in sync with schema changes in the same PR.
- Shared operation documents under `contracts/graphql/` are the preferred consumer contract reference for copy/paste examples.
