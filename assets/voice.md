## Writing Standards

These are quality guidelines; do not treat them as rigid mechanical checkboxes. A chapter must first feel natural and coherent, and secondarily satisfy these quality checks.

- Establish conflict, suspense, desire, or abnormality quickly at the opening; avoid abstract retrospectives.
- Advance the plot using action, dialogue, and sensory details; minimize exposition and summaries.
- Character dialogue must reflect distinct identity, subtext, and underlying intent—never preach or lecture.
- Show emotions through physical reactions and character choices rather than static labels.
- Relationship changes require triggering events; do not jump from total strangers to absolute trust in a single chapter.
- Reveal secrets incrementally; never prematurely explain major mysteries not required by the outline.
- Ending hooks can take the form of crisis, difficult choice, emotional lingering, relationship shifts, or unfulfilled goals—not every chapter requires an exaggerated cliffhanger.
- **De-AI Tone**: Avoid all patterns listed in `reference_pack.references.anti_ai_tone` during drafting (spanning structure, vocabulary, description, dialogue, and pacing). Mechanically countable fatigued words and formulaic phrases are bound by thresholds in `working_memory.user_rules.structured`, strictly validated upon commit.
- **Sentence Structure Diversity**: `episodic_memory.style_stats` (if present) represents automated statistical metrics of your previously drafted text—a mirror of your own writing habits. Actively reduce high-frequency patterns recorded here. Common sources of repetitive phrasing include correction structures ("not X, but Y"), single timing units ("in a few breaths"), and consecutive simile patterns. Rotate ending styles (short punchy closures, dialogue lingering, scene imagery, suspense questions) relative to recent chapters; avoid opening every chapter with time indicators like "at night / at dawn / upon waking".
- **No Recapitulation**: Summaries, foreshadowing, and state entries in `episodic_memory` are reference memos of already written content for continuity—NOT raw material to be written into the current chapter. Do not re-explain prior chapter events unless required by new narrative perspective. Re-stating past events like a recap is strictly forbidden (cross-chapter verbatim repetition will be flagged in `style_stats.repeated_sentences`).
