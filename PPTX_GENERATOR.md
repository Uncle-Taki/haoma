# PowerPoint Presentation Generator

This tool converts Marp markdown presentations to PowerPoint (PPTX) format.

## Features

- Converts Marp markdown slides to PowerPoint presentations
- Supports slide titles, bullet points, and nested content
- Handles Mermaid diagram placeholders
- Preserves code blocks with special formatting
- Extracts and preserves speaker notes

## Usage

### Building the Generator

```bash
make build-pptx-generator
```

Or manually:
```bash
go build -o pptx-generator ./cmd/pptx-generator/main.go
```

### Generating a Presentation

```bash
./pptx-generator -input data/yggdrasil-slides.md -output yggdrasil-presentation.pptx
```

Or use the Makefile target:
```bash
make generate-pptx INPUT=data/yggdrasil-slides.md OUTPUT=output.pptx
```

### Command Line Options

- `-input`: Path to the input Marp markdown file (required)
- `-output`: Path to the output PowerPoint file (optional, defaults to input filename with .pptx extension)

## Marp Markdown Format

The generator expects Marp-formatted markdown with the following structure:

```markdown
---
marp: true
title: Presentation Title
---
# Slide Title
## Subtitle
- Bullet point 1
- Bullet point 2
  - Nested bullet

Notes:
- Speaker note 1
- Speaker note 2
---
# Another Slide
...
```

### Supported Elements

1. **Frontmatter** (YAML) - Skipped during conversion
2. **Slide Separators** - `---` creates a new slide
3. **Titles** - `# Title` creates slide title
4. **Subtitles** - `## Subtitle` adds subtitle text
5. **Bullet Points** - `- Item` and `  - Nested item`
6. **Code Blocks** - ` ```language ... ``` ` with special formatting
7. **Mermaid Diagrams** - ` ```mermaid ... ``` ` shown as placeholders
8. **Notes** - `Notes:` section converts to PowerPoint speaker notes

## Example

An example slide deck is included in `data/yggdrasil-slides.md` showcasing:
- Technical architecture slides
- Mermaid diagrams (sequence, flowchart, ERD)
- Code snippets
- Nested bullet points
- Speaker notes

Generate it with:
```bash
./pptx-generator -input data/yggdrasil-slides.md -output yggdrasil.pptx
```

## Limitations

- Mermaid diagrams are shown as text placeholders (not rendered)
- Images are not yet supported
- Tables are not yet supported  
- Limited text formatting (bold, italic, etc.)
- Fixed slide layout (title + content)

## Technical Details

The generator creates PowerPoint presentations by:
1. Parsing Marp markdown into slide structures
2. Building PPTX XML files according to Office Open XML format
3. Creating a valid ZIP archive with proper relationships

No external PowerPoint library license is required as the generator creates the XML structure directly.
