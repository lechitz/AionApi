# Database Seed Data

This directory contains seed data for populating the database with test/development data.

## 📋 Overview

The seed structure follows a hierarchical relationship:
```
Users → Categories → Tags → Records
```

## 🗂️ Seed Files

### 1. `user.sql`
Creates test users in the system.

**Users:**
- User ID 1: Primary test user (focus of category/tag/record seeds)
- User ID 2-5: Additional test users

### 2. `category.sql`
Creates tag categories for organizing tags.

**User 1 Categories (8 total):**
1. **Learning** - Books, courses, studies (#F8B400, book icon)
2. **Fitness** - Exercise, sports, physical activities (#E94F37, dumbbell icon)
3. **Mindfulness** - Meditation, breathing, mental health (#9C27B0, spa icon)
4. **Career** - Work projects, coding, professional development (#1976D2, briefcase icon)
5. **Social** - Friends, family, social connections (#FF6F00, users icon)
6. **Creative** - Art, music, writing, hobbies (#00ACC1, palette icon)
7. **Nutrition** - Meal planning, cooking, healthy eating (#388E3C, utensils icon)
8. **Rest** - Sleep, relaxation, recovery (#5E35B1, moon icon)

### 3. `tags.sql`
Creates specific tags within each category (2 tags per category = 16 total).

**User 1 Tags (16 total):**
- **Learning:** Reading, Online Course
- **Fitness:** Running, Weight Training
- **Mindfulness:** Meditation, Breathing
- **Career:** Coding, Meeting
- **Social:** Family Time, Friends
- **Creative:** Writing, Music
- **Nutrition:** Meal Prep, Hydration
- **Rest:** Sleep, Stretching

### 4. `records.sql`
Creates habit tracking records (128 total for User 1).

**Distribution:**
- **Total records:** 128
- **Days covered:** 16 (2025-01-01 to 2025-01-16)
- **Records per day:** 8 (one from each category)
- **User focused:** User ID 1

**Record Types:**
- Mixed duration activities (5 minutes to 4 hours)
- Various value types (distance, pages, minutes, liters)
- Realistic timestamps throughout each day
- All marked as "completed" status

## 🔄 Seed Order

Seeds must be applied in this order due to foreign key constraints:

```bash
1. user.sql          # Creates users first
2. category.sql      # Creates categories (depends on users)
3. tags.sql          # Creates tags (depends on categories)
4. records.sql       # Creates records (depends on tags)
```

## 🚀 How to Apply Seeds

### Option 1: Using Make Command
```bash
make seed
```

### Option 2: Manual psql
```bash
psql -U aion_user -d aion_db -f infrastructure/db/seed/user.sql
psql -U aion_user -d aion_db -f infrastructure/db/seed/category.sql
psql -U aion_user -d aion_db -f infrastructure/db/seed/tags.sql
psql -U aion_user -d aion_db -f infrastructure/db/seed/records.sql
```

### Option 3: Docker Exec
```bash
docker exec -i aion-postgres psql -U aion_user -d aion_db < infrastructure/db/seed/user.sql
docker exec -i aion-postgres psql -U aion_user -d aion_db < infrastructure/db/seed/category.sql
docker exec -i aion-postgres psql -U aion_user -d aion_db < infrastructure/db/seed/tags.sql
docker exec -i aion-postgres psql -U aion_user -d aion_db < infrastructure/db/seed/records.sql
```

## 🧪 Testing Queries

After seeding, test with these GraphQL queries:

### Get All Categories for User 1
```graphql
query {
  categories {
    id
    name
    description
    colorHex
    icon
  }
}
```

### Get Tags by Category
```graphql
query {
  tagsByCategoryId(categoryId: "1") {
    id
    name
    description
  }
}
```

### Get Records for a Specific Day
```graphql
query {
  recordsByDay(date: "2025-01-01") {
    id
    title
    eventTime
    categoryId
    tagId
    value
    durationSeconds
  }
}
```

### Get Records by Tag
```graphql
query {
  recordsByTag(tagId: "1", limit: 20) {
    id
    title
    eventTime
    value
  }
}
```

### Get Records in Date Range
```graphql
query {
  recordsBetween(
    startDate: "2025-01-01T00:00:00Z"
    endDate: "2025-01-16T23:59:59Z"
    limit: 128
  ) {
    id
    title
    eventTime
    categoryId
    tagId
  }
}
```

## 📊 Data Statistics (User 1)

| Entity | Count | Distribution |
|--------|-------|-------------|
| Categories | 8 | Evenly distributed across life aspects |
| Tags | 16 | 2 per category |
| Records | 128 | 8 per day × 16 days |

### Records by Category (approximate)
- Learning: 16 records
- Fitness: 16 records
- Mindfulness: 16 records
- Career: 16 records
- Social: 16 records
- Creative: 16 records
- Nutrition: 16 records
- Rest: 16 records

## 🔧 Regenerating Records

The `records.sql` file is generated using `generate_records.py`:

```bash
cd infrastructure/db/seed
python3 generate_records.py > records.sql
```

This script creates realistic, varied records with:
- Random but sensible durations
- Appropriate values for each activity type
- Spread throughout the day
- Realistic descriptions

## 🎯 Design Principles

1. **Realistic Data**: Activities mirror real-world habit tracking
2. **Variety**: Mix of different activity types, durations, and times
3. **Consistency**: Regular daily tracking pattern
4. **Completeness**: All categories and tags are represented
5. **Testability**: Enough data to test all query types and filters

## 📝 Notes

- All timestamps use 'America/Sao_Paulo' timezone
- All records marked as 'completed' status
- Source varies between 'mobile_app' and 'web'
- Values represent different metrics (minutes, km, liters, etc.)
- User 1 is the primary test user with full dataset
- Other users (2-5) have minimal data for multi-user testing

## 🔄 Updating Seeds

When adding new seeds:
1. Maintain foreign key order (users → categories → tags → records)
2. Use consistent timestamp format: `YYYY-MM-DD HH:MM:SS`
3. Include descriptive comments
4. Test seeds in isolated database first
5. Update this README with new counts/structure

---

**Last Updated:** 2025-01-14  
**Total Records (User 1):** 128  
**Date Range:** 2025-01-01 to 2025-01-16

