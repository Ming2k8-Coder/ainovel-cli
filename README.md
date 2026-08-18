# ainovel-cli

Fully autonomous AI long-form novel creation engine. Driven by a deterministic engine codebase paired with Large Language Models (LLMs) invoked precisely at semantic decision boundaries: the Engine routes workflows based on factual state to dispatch three autonomous creation agents (**Architect** / **Writer** / **Editor**), while invoking the **Arbiter** on-demand for semantic decisions. 

From a single sentence prompt to a complete long-form novel, the entire process runs end-to-end without requiring human intervention.

<p align="center">
  <img src="scripts/sample.gif" alt="ainovel-cli demo" width="800">
  <img src="scripts/novel.png" alt="ainovel-cli bg" width="800">
</p>

## 🌟 Key Features

- **Deterministic Engine + Multi-Agent Collaboration** — The Engine dispatches three autonomous creation agents (**Architect** / **Writer** / **Editor**) based on a factual decision route. The main loop incurs **zero LLM overhead** for routing, making system behavior fully testable and exhaustive.
- **Auditable Semantic Arbiter** — Decisions such as planner selection, user intervention triage, and deadlock recovery are handled via single-call Arbiter functions. Every decision is logged to disk for deterministic replay. Simple, stable, and zero complex graph orchestration.
- **Step-Level Checkpoint Recovery** — State checkpoints are written to disk after every successful tool execution. Upon crash or network failure, recovery is exact to the step: `plan` -> `draft` -> `check` -> `commit`.
- **Rolling Volume/Arc Planning** — Long-form novels are no longer planned up-front in a hollow manner. Initially, only the skeleton of the first 2 volumes + detailed chapters for Arc 1 are created. Subsequent arcs and volumes are dynamically expanded by the Architect as writing progresses, referencing past summaries and character snapshots.
- **Smart Related-Chapter Recommendations** — During each chapter's drafting, relevant historical chapters are automatically retrieved across 4 dimensions: *Foreshadowing, Character Appearances, State Changes, and Relationships*, paired with previews of upcoming chapters to guarantee continuity across 500+ chapters.
- **Adaptive Context Strategy** — Automatically switches between Full / Sliding Window / Hierarchical Summaries based on total chapter count, supporting massive 500+ chapter narratives.
- **Seven-Dimensional Quality Review** — The Editor audits drafted chapters across 7 dimensions: *Setting Consistency, Character Behavior, Pacing, Narrative Coherence, Foreshadowing, Hooks, and Literary Aesthetic Quality*. The aesthetic dimension is subdivided into *Descriptive Texture, Narrative Technique, Dialogue Differentiation, Diction Quality, and Emotional Impact* (each requiring direct quote evidence from the draft).
- **Real-Time User Steering** — Inject modification feedback into the prompt input box at any time during active writing (without pausing). The system automatically evaluates the blast radius and rewrites impacted chapters.
- **Optional Step-by-Step Review Mode** — Fully automated by default. When fine-grained control is desired, enable `/review on`. Each `/next` command unlocks exactly one new chapter; reworks and crash recoveries do not consume permission tokens.
- **Dual Interface: TUI + Headless** — Observe and steer interactively via a rich Terminal UI (TUI), or run headless on a server, NAS, or CI pipeline.
- **Multi-LLM Support** — Seamlessly switch between OpenRouter, Anthropic (Claude), Google Gemini, OpenAI, DeepSeek, Ollama (local), and more.

---

## 🏗️ Architecture

Core Design Principle: **Deterministic Fact Layer, Autonomous Semantic Layer**.

Enumerated state transitions are executed by deterministic code (`Engine` + `Route`). Clear-boundary judgements consult LLM functions on-demand (`Arbiter`). Open-ended creative drafting is handled by autonomous LLM loops (`Workers`). 

In short: A serial deterministic Engine, three autonomous Workers, a few on-demand Arbiter functions, and a file-system Fact Layer.

```text
┌─────────────────────────────────────────────────────────┐
│              Host / Engine (Deterministic)              │
│  Read Store → Route → Dispatch Worker → Main Loop       │
│  Trigger Arbiter / Triage Steering / Failure → Arbiter  │
└────┬──────────┬──────────┬─────────────┬────────────────┘
     │          │          │             │
 ┌───▼────┐ ┌───▼───┐ ┌────▼────┐   ┌────▼────┐
 │Architect│ │Writer │ │ Editor  │   │ Arbiter │
 │(LLM Loop│ │(LLM Loop│ │(LLM Loop│   │(LLM Func│
 └───┬────┘ └───┬───┘ └────┬────┘   └─────────┘
     └──────────┼──────────┘
                │ Tool Calls (IO + Checkpoint)
┌───────────────▼─────────────────────────────────────────┐
│                       Store                             │
│  Progress / Checkpoint / Outline / Drafts / ...         │
└─────────────────────────────────────────────────────────┘
```

- **Engine** — Reads facts from Store each turn, dispatches Workers per Route table. Makes execution decisions without literary judgment. Crash recovery = read store and resume; session-less.
- **Arbiter** — On-demand semantic judge (planner selection, user intervention triage, failure/deadlock recovery). Factual input, structured decision output. Auditable & replayable.
- **Workers** — Architect / Writer / Editor run independent LLM context loops, collaborating via Store artifacts.
- **Tools** — Single-file atomic IO with idempotent replay. Chapter commits use persistent Saga + Checkpoints, returning factual JSON only without instruction clutter.

### Agent Roles & Responsibilities

| Role | Responsibilities | Tools Used |
|---|---|---|
| **Arbiter** | Semantic judgment: Planner selection, user intervention triage, deadlock recovery. | None (Single LLM call returning structured decision). |
| **Architect** | Generates title, book intro, premise, outline, character dossiers, world rules. | `novel_context`, `save_book`, `save_foundation` |
| **Writer** | Autonomously conceptualizes, drafts, self-checks, and commits a chapter. | `novel_context`, `read_chapter`, `plan_chapter`, `draft_chapter`, `check_consistency`, `commit_chapter` |
| **Editor** | Reads raw draft, auditing at two levels: Narrative Structure & Aesthetic Quality. | `novel_context`, `read_chapter`, `save_review`, `save_arc_summary`, `save_volume_summary` |

---

## 🔄 Writing Workflow

```text
User Request → Arbiter Selects Planner → Architect Plans Skeleton + Arc 1 → Writer Drafts Chapters → Editor Arc Review
                   (Disk Logged)                                                    ↑                      │
                                                                                    ├── Rewrite/Polish ◄───┘
                                                                                    │
                                                                              Architect Expands Next Arc/Vol
                                                                              (Ref Summaries + Character Snapshots)
```

Each "Who to dispatch next" step is derived by the Engine's `Route` table based on Store facts (tested via exhaustive state-combination tests), consuming **zero LLM tokens**.

### Writer's Sequential Chapter Execution Flow:
1. `novel_context` — Load context (previous summaries, foreshadowing, character states, style rules, related chapter recommendations).
2. `read_chapter` — Read back prior text to match tone, rhythm, and pacing.
3. `plan_chapter` — Plan chapter goals, central conflict, and emotional arc.
4. `draft_chapter` — Draft full chapter prose.
5. `check_consistency` — Audit drafted text against factual state (must run after `draft`).
6. `commit_chapter` — Commit final draft, persisting facts (`arc_end` / `next_chapter` / feedback pool). Engine derives next step via Route table.

---

## 🚦 State Transition Rules

System states are split into two layers:
- **Phase** — High-level lifecycle stage (Initialization, Premise, Outline, Writing, Completed).
- **Flow** — Active workflow within the Writing phase (Normal Writing, Reviewing, Rewriting, Polishing, Steering).

### Phase Transitions (Forward-Only)
`init` -> `premise` -> `outline` -> `writing` -> `complete`

### Flow Transitions
`writing` <-> `reviewing` / `rewriting` / `polishing` / `steering`

---

## 📦 Assets & Customization Structure

For localizing system prompts and custom lore (e.g., world-building, writing rules, voice guidelines):

1. **System Prompts (`assets/prompts/`)**:
   - `architect-long.md` & `architect-short.md`: Architect prompts for world building and outline expansion.
   - `writer.md`: Writer prompt for prose drafting.
   - `editor.md`: Editor prompt for 7-dimensional quality auditing.
   - `arbiter-*.md`: Arbiter prompt templates for semantic decision-making.
2. **Style & Voice Guidelines (`assets/styles/` & `assets/voice.md`)**:
   - Writing style instructions, allowed/prohibited vocabulary, and tone anchors.

---

## ⚡ Quick Start

```bash
# Run locally with Go
go run ./cmd/ainovel

# Build executable binary
go build -o ainovel ./cmd/ainovel
./ainovel
```
