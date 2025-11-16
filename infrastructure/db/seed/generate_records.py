#!/usr/bin/env python3
"""
Generate 128 realistic record seeds for user_id=1
Distributed over 16 days (8 records per day)
Uses the actual categories and tags from the database
"""

from datetime import datetime, timedelta
import random

# Configuration
USER_ID = 1
START_DATE = datetime(2025, 1, 1)
RECORDS_PER_DAY = 8
TOTAL_DAYS = 16

# Map category/tag to activity templates (title_template, desc_template, min_dur, max_dur, min_val, max_val)
# Duration in minutes, value depends on activity type
ACTIVITIES = {
    ('saude_fisica', 'Stretching'): [
        ('Morning Stretching', 'Full body stretch routine', 10, 20, 10, 20),
        ('Evening Stretching', 'Flexibility and recovery', 15, 25, 15, 25),
    ],
    ('saude_fisica', 'Push-up'): [
        ('Push-up Challenge', '{} reps completed', 5, 15, 20, 100),
        ('Push-up Sets', '3 sets of push-ups', 10, 20, 30, 60),
    ],
    ('saude_fisica', 'ABS'): [
        ('ABS Workout', 'Core strengthening routine', 10, 20, 50, 150),
        ('Plank Challenge', 'Hold and variations', 8, 15, 60, 180),
    ],
    ('saude_fisica', 'Run'): [
        ('Morning Run', '{}km outdoor run', 20, 45, 3, 10),
        ('Interval Training', 'HIIT running workout', 25, 40, 4, 8),
    ],
    ('saude_fisica', 'Pull-up'): [
        ('Pull-up Session', '{} reps completed', 5, 15, 5, 30),
        ('Pull-up Training', 'Back and biceps workout', 10, 20, 10, 40),
    ],
    ('saude_fisica', 'Sunbathe'): [
        ('Sunbathing', 'Vitamin D exposure', 15, 30, 15, 30),
        ('Outdoor Relaxation', 'Sun and fresh air', 20, 40, 20, 40),
    ],
    ('saude_fisica', 'Gym'): [
        ('Gym Session', 'Full body workout', 60, 120, 60, 120),
        ('Strength Training', 'Weight lifting routine', 45, 90, 45, 90),
    ],
    ('saude_fisica', 'Walking'): [
        ('Walk', '{}km walk', 20, 60, 2, 8),
        ('Evening Walk', 'Relaxing stroll', 30, 50, 3, 6),
    ],
    ('meditacao', 'Meditation'): [
        ('Morning Meditation', 'Mindfulness practice', 10, 30, 10, 30),
        ('Guided Meditation', 'App-guided session', 15, 25, 15, 25),
    ],
    ('saude_mental', 'Planning'): [
        ('Daily Planning', 'Organized tasks and goals', 15, 30, 15, 30),
        ('Week Planning', 'Strategic planning session', 30, 60, 30, 60),
    ],
    ('saude_mental', 'Diary'): [
        ('Journal Writing', 'Daily reflections', 15, 30, 15, 30),
        ('Gratitude Journal', 'Positive thoughts', 10, 20, 10, 20),
    ],
    ('saude_mental', 'Reading'): [
        ('Reading: {}', 'Read {} pages', 30, 90, 20, 80),
        ('Book Session', 'Deep reading', 45, 120, 30, 100),
    ],
    ('saude_mental', 'TheNews'): [
        ('News Reading', 'Caught up with news', 15, 30, 15, 30),
        ('Current Events', 'News analysis', 20, 40, 20, 40),
    ],
    ('saude_mental', 'Emails'): [
        ('Email Processing', 'Inbox management', 20, 45, 20, 45),
        ('Email Cleanup', 'Organized communications', 15, 30, 15, 30),
    ],
    ('estudo_trabalho', 'Dev'): [
        ('Development: {}', 'Feature implementation', 60, 240, 60, 240),
        ('Code Sprint', 'Focused coding session', 90, 180, 90, 180),
    ],
    ('estudo_trabalho', 'College'): [
        ('College Class', 'Attended lecture', 60, 120, 60, 120),
        ('College Study', 'Assignment work', 90, 180, 90, 180),
    ],
    ('estudo_trabalho', 'Golang'): [
        ('Golang Practice', 'Coding exercises', 45, 120, 45, 120),
        ('Golang Study', 'Learning advanced concepts', 60, 90, 60, 90),
    ],
    ('estudo_trabalho', 'Notion'): [
        ('Notion Organization', 'Database setup', 30, 60, 30, 60),
        ('Notion Planning', 'Project management', 20, 45, 20, 45),
    ],
    ('estudo_trabalho', 'GPT'): [
        ('GPT Session', 'AI-assisted learning', 30, 90, 30, 90),
        ('GPT Practice', 'Prompt engineering', 45, 75, 45, 75),
    ],
    ('estudo_trabalho', 'Full Cycle'): [
        ('Full Cycle Course', 'Module completion', 90, 180, 90, 180),
        ('Full Cycle Practice', 'Hands-on exercises', 60, 120, 60, 120),
    ],
    ('estudo_trabalho', 'FreeCodeCamp'): [
        ('FreeCodeCamp', 'Challenges completed', 60, 120, 60, 120),
        ('FreeCodeCamp Project', 'Certification project', 90, 180, 90, 180),
    ],
    ('estudo_trabalho', 'Coursera'): [
        ('Coursera Lecture', 'Video lessons', 45, 90, 45, 90),
        ('Coursera Assignment', 'Quiz and exercises', 60, 120, 60, 120),
    ],
    ('estudo_trabalho', 'Course'): [
        ('Online Course', 'Learning session', 60, 120, 60, 120),
        ('Course Practice', 'Practical exercises', 45, 90, 45, 90),
    ],
    ('estudo_trabalho', 'Aion'): [
        ('Aion Development', 'Project work', 90, 240, 90, 240),
        ('Aion Planning', 'Architecture design', 60, 120, 60, 120),
    ],
    ('estudo_trabalho', 'Work'): [
        ('Work Session', 'Professional tasks', 120, 480, 120, 480),
        ('Work Meeting', 'Team collaboration', 30, 90, 30, 90),
    ],
    ('estudo_trabalho', 'RD'): [
        ('R&D Session', 'Research and experiments', 90, 180, 90, 180),
        ('Technical Research', 'Exploring new tech', 60, 120, 60, 120),
    ],
    ('estudo_trabalho', 'Podcast'): [
        ('Podcast: {}', 'Educational listening', 30, 90, 30, 90),
        ('Tech Podcast', 'Industry insights', 45, 75, 45, 75),
    ],
    ('estudo_trabalho', 'AudioBook'): [
        ('AudioBook: {}', 'Learning on the go', 45, 120, 45, 120),
        ('AudioBook Session', 'Professional development', 60, 90, 60, 90),
    ],
    ('estudo_trabalho', 'Finance'): [
        ('Finance Study', 'Investment learning', 45, 90, 45, 90),
        ('Finance Planning', 'Budget review', 30, 60, 30, 60),
    ],
    ('estudo_trabalho', 'Interview'): [
        ('Interview Prep', 'Practice questions', 60, 120, 60, 120),
        ('Mock Interview', 'Practice session', 45, 90, 45, 90),
    ],
    ('idiomas', 'English'): [
        ('English Practice', 'Conversation and grammar', 30, 60, 30, 60),
        ('English Lesson', 'Structured learning', 45, 90, 45, 90),
    ],
    ('idiomas', 'Spanish'): [
        ('Spanish Practice', 'Vocabulary and conversation', 30, 60, 30, 60),
        ('Spanish Lesson', 'Grammar exercises', 45, 75, 45, 75),
    ],
    ('idiomas', 'German'): [
        ('German Study', 'Language practice', 30, 60, 30, 60),
        ('German Lesson', 'Structured learning', 45, 90, 45, 90),
    ],
    ('idiomas', 'French'): [
        ('French Practice', 'Conversation practice', 30, 60, 30, 60),
        ('French Lesson', 'Grammar and vocab', 45, 75, 45, 75),
    ],
    ('idiomas', 'Chinese'): [
        ('Chinese Study', 'Characters and tones', 30, 60, 30, 60),
        ('Chinese Lesson', 'Structured practice', 45, 90, 45, 90),
    ],
    ('pessoal', 'My'): [
        ('Personal Time', 'Self care and reflection', 30, 90, 30, 90),
        ('Me Time', 'Relaxation and hobbies', 45, 120, 45, 120),
    ],
    ('pessoal', 'Travel'): [
        ('Travel Planning', 'Trip preparation', 45, 90, 45, 90),
        ('Travel Research', 'Destination exploration', 30, 60, 30, 60),
    ],
    ('pessoal', 'Movie'): [
        ('Movie: {}', 'Watched film', 90, 180, 90, 180),
        ('Movie Night', 'Cinema experience', 120, 150, 120, 150),
    ],
    ('pessoal', 'Series'): [
        ('Series: {}', 'Watched episodes', 45, 120, 45, 120),
        ('Binge Watching', 'TV series marathon', 90, 240, 90, 240),
    ],
    ('pessoal', 'Game'): [
        ('Gaming Session', 'Video games', 60, 180, 60, 180),
        ('Game: {}', 'Playing session', 90, 240, 90, 240),
    ],
    ('pessoal', 'Beach'): [
        ('Beach Day', 'Sun and ocean', 120, 300, 120, 300),
        ('Beach Walk', 'Seaside stroll', 45, 90, 45, 90),
    ],
    ('pessoal', 'Hanging Out'): [
        ('Hanging Out', 'Social time with friends', 90, 240, 90, 240),
        ('Social Gathering', 'Meetup with friends', 120, 180, 120, 180),
    ],
    ('trabalho_de_casa', 'Housework'): [
        ('Housework', 'Cleaning and organizing', 45, 120, 45, 120),
        ('House Cleaning', 'Deep clean session', 60, 180, 60, 180),
    ],
    ('trabalho_de_casa', 'Supermarket'): [
        ('Grocery Shopping', 'Weekly shopping', 45, 90, 45, 90),
        ('Supermarket Run', 'Quick shopping trip', 30, 60, 30, 60),
    ],
    ('trabalho_de_casa', 'Doctor'): [
        ('Doctor Appointment', 'Medical checkup', 60, 120, 60, 120),
        ('Health Consultation', 'Medical visit', 45, 90, 45, 90),
    ],
    ('outros', 'OFF'): [
        ('Day OFF', 'Rest and recovery', None, None, 480, 720),
        ('休息时间', 'Relaxation day', None, None, 360, 600),
    ],
    ('outros', 'Travelling'): [
        ('Travel Day', 'Journey time', 120, 480, 120, 480),
        ('In Transit', 'Moving between locations', 90, 360, 90, 360),
    ],
}

# Sample data for templates
BOOKS = ['Atomic Habits', 'Deep Work', 'Sapiens', 'Clean Code', 'The Lean Startup']
DEV_FEATURES = ['GraphQL API', 'Authentication', 'Database optimization', 'Frontend components']
PODCASTS = ['Syntax.fm', 'Go Time', 'Software Engineering Daily', 'The Changelog']
AUDIOBOOKS = ['Designing Data-Intensive Applications', 'The Pragmatic Programmer', 'Domain-Driven Design']
MOVIES = ['Inception', 'The Matrix', 'Interstellar', 'The Shawshank Redemption']
SERIES = ['Breaking Bad', 'The Office', 'Game of Thrones', 'Stranger Things']
GAMES = ['Zelda', 'Elden Ring', 'The Witcher', 'Red Dead Redemption']

def generate_record(day, seq, category, tag, activity_template):
    template, desc_template, min_dur, max_dur, min_val, max_val = activity_template

    # Generate time for this record (spread throughout the day)
    hour = 6 + (seq * 2) + random.randint(0, 1)
    minute = random.randint(0, 59)
    event_time = START_DATE + timedelta(days=day, hours=hour, minutes=minute)

    # Duration and value
    if min_dur is not None:
        duration = random.randint(min_dur, max_dur) * 60  # Convert to seconds
        value = float(duration // 60)  # Minutes
        duration_str = f"{duration}"
    else:
        duration = None
        duration_str = "NULL"
        value = round(random.uniform(min_val, max_val), 1)

    # Recorded at (usually shortly after event)
    recorded_at = event_time + timedelta(minutes=random.randint(0, 15))

    # Customize title/description based on category/tag
    if '{}' in template:
        if 'Reading' in tag:
            title = template.format(random.choice(BOOKS))
            desc = desc_template.format(random.randint(20, 80))
        elif 'Development' in template or 'Dev' in tag:
            title = template.format(random.choice(DEV_FEATURES))
            desc = desc_template
        elif 'Podcast' in tag:
            title = template.format(random.choice(PODCASTS))
            desc = desc_template
        elif 'AudioBook' in tag:
            title = template.format(random.choice(AUDIOBOOKS))
            desc = desc_template
        elif 'Movie' in tag:
            title = template.format(random.choice(MOVIES))
            desc = desc_template
        elif 'Series' in tag:
            title = template.format(random.choice(SERIES))
            desc = desc_template
        elif 'Game' in tag:
            title = template.format(random.choice(GAMES))
            desc = desc_template
        elif 'km' in desc_template:
            km = round(random.uniform(3, 10), 1)
            title = template
            desc = desc_template.format(km)
        else:
            num = random.randint(20, 100)
            title = template.format(num)
            desc = desc_template
    else:
        title = template
        desc = desc_template

    # Escape single quotes
    title = title.replace("'", "''")
    desc = desc.replace("'", "''")

    return f"""  (1, '{title}', '{desc}',
   (SELECT category_id FROM aion_api.tag_categories WHERE name='{category}' AND user_id=1 LIMIT 1),
   (SELECT tag_id FROM aion_api.tags WHERE name='{tag}' AND user_id=1 LIMIT 1),
   '{event_time.strftime("%Y-%m-%d %H:%M:%S")}', '{recorded_at.strftime("%Y-%m-%d %H:%M:%S")}',
   {duration_str}, {value}, 'seed', 'America/Sao_Paulo', 'completed',
   '{recorded_at.strftime("%Y-%m-%d %H:%M:%S")}', '{recorded_at.strftime("%Y-%m-%d %H:%M:%S")}')"""

# Generate all records
print("-- " + "=" * 76)
print("-- Records Seed Data - User ID 1 (Test User)")
print("-- " + "=" * 76)
print("-- 128 records distributed over 16 days (8 records per day)")
print(f"-- Date range: {START_DATE.strftime('%Y-%m-%d')} to {(START_DATE + timedelta(days=TOTAL_DAYS-1)).strftime('%Y-%m-%d')}")
print("-- " + "=" * 76)
print()

# Get all category/tag combinations
combos = list(ACTIVITIES.keys())
random.shuffle(combos)

for day in range(TOTAL_DAYS):
    date_str = (START_DATE + timedelta(days=day)).strftime('%Y-%m-%d')
    print(f"-- DAY {day + 1}: {date_str}")
    print("INSERT INTO aion_api.records (user_id, title, description, category_id, tag_id, event_time, recorded_at, duration_seconds, value, source, timezone, status, created_at, updated_at) VALUES")

    records = []
    for seq in range(RECORDS_PER_DAY):
        # Cycle through categories/tags to ensure variety
        combo = combos[seq % len(combos)]
        activity_templates = ACTIVITIES[combo]
        activity = random.choice(activity_templates)

        record = generate_record(day, seq, combo[0], combo[1], activity)
        records.append(record)

    print(",\n".join(records) + ";\n")

print("-- " + "=" * 76)
print(f"-- Total: {TOTAL_DAYS * RECORDS_PER_DAY} records created")
print("-- " + "=" * 76)

