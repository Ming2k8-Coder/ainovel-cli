# Assets Content Map

Before adding a "paragraph / reference document / rule" to the system, check the table below to determine its location and wiring method.

| Directory | What it Contains | Consumed By | Wiring Method |
|---|---|---|---|
| `prompts/` | Worker system prompts (`writer` / `editor` / `architect`×2), Arbiter decision prompts, and one-off task prompts (`import` / `simulation` / `revision`) | `agents/build.go`, `internal/arbiter`, imp / sim / revision runner | `load.go` `Prompts` fields. Note: `simulation_guidance` is injected during `load.go` loading and is not visible inside the `.md` files directly. |
| `references/` | Genre-agnostic writing reference materials. Not included in system prompts; sliced by `novel_context` per role/chapter and injected into `reference_pack` | `writer` / `editor` / `architect` | **3-point wiring**: Add field to `tools.References` + read in `load.go` `loadReferences` + inject in `novel_context.go` (`writerReferences` / `architectReferences`). Merely placing files in directory will NOT auto-load them. |
| `references/genres/<style>/` | Genre-specific knowledge (`style-references` / `arc-templates`) | Same as above, loaded when `style != default` | `load.go` `loadReferences` |
| `rules/` | Deprecated legacy built-in rules directory. Mechanical baselines migrated to code; user rules loaded from natural language snapshots at `~/.ainovel/rules/*.md` / `./.ainovel/rules/*.md` | `userrules.Service` normalized to `meta/user_rules.json`; injected by `novel_context`; checked by `commit_chapter` | Built-in baselines in `internal/rules/snapshot.go` (`SystemDefaults()`); user `.md` files require zero formatting or YAML, normalized via natural language. |
| `styles/<style>.md` | Genre-specific writing style instructions | Concatenated into **writer** system prompt (`agents/build.go`) | Filename equals `config.style` value. Conceptually paired with `references/genres/<style>/`: styles are writing instructions, references are knowledge materials. |

## New Content Location Checklist (5 Questions)

1. Must this process be **guaranteed** programmatically? → Do NOT write a prompt; write code constraints (`StopAfterTools` / Tool Guards / Flow Router).
2. Is this an Arbiter decision criterion? → Table-driven routing goes to `internal/flow/router.go`; semantic judgment goes to `prompts/arbiter-*.md`.
3. Is this a genre writing instruction? → `styles/<style>.md`.
4. Is this reference knowledge material? → `references/`.
5. Is this a user-defined custom writing rule? → Place in `./.ainovel/rules/*.md`.
