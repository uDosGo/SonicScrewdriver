# Getting Started with Sonic Dev Flow

## Quick Start

```bash
# Clone the repository
git clone https://github.com/fredporter/sonic-screwdriver.git
cd sonic-screwdriver

# Initialize submodules (including Ventoy)
git submodule update --init --recursive

# Explore the dev flow structure
ls -la dev/
```

## Development Flow Overview

The Sonic Dev Flow is designed to keep active work visible while providing clear pathways for integration, validation, and promotion.

### Core Components

```
dev/
├── active/              # Current 2-week execution slice
├── integration/         # Cross-component coordination
├── process/             # Workflow definitions and checklists
├── compost/             # Archived material
└── GETTING-STARTED.md   # This guide
```

## Daily Workflow

### 1. Start Your Day

```bash
# Review active work
cat dev/active/active-index.md

# Check yesterday's progress
cat dev/active/execution-notes.md

# Update execution notes for today
$EDITOR dev/active/execution-notes.md
```

### 2. Work On Tasks

```bash
# Example: Working on Ventoy integration

# Review integration brief
cat dev/integration/VENTOY-INTEGRATION-BRIEF.md

# Check validation checklist
cat dev/process/checklists/ventoy-validation.md

# Run build
./modules/ventoy/build.sh

# Run validation
./modules/ventoy/validate.sh
```

### 3. Update Progress

```bash
# Mark tasks as completed in execution notes
# Use status markers:
# - [ ] not started
# - [-] in progress
# - [x] done
# - [!] blocked

# Update active index if priorities change
$EDITOR dev/active/active-index.md
```

### 4. End Your Day

```bash
# Review what was accomplished
cat dev/active/execution-notes.md

# Update end-of-day status
$EDITOR dev/active/execution-notes.md

# Commit changes
git add dev/active/
git commit -m "Update execution notes for YYYY-MM-DD"
```

## Ventoy Development

### Build and Validate

```bash
# Build Ventoy with patches
./modules/ventoy/build.sh

# Run validation checklist
./modules/ventoy/validate.sh

# Check validation results
cat dev/process/checklists/ventoy-validation.md
```

### Promotion Workflow

```bash
# After successful validation
mkdir -p release/ventoy
cp -v build/ventoy/* release/ventoy/

# Verify promotion
ls -la release/ventoy/

# Generate manifest
ls -la release/ventoy > release/ventoy/MANIFEST.txt
```

## Working with Modules

### Sonic Home

```bash
# Build sonic-home module
cd modules/sonic-home
go build ./cmd/sonic-home

# Run sonic-home commands
./sonic-home version
./sonic-home pack --help
```

### Sonic Express

```bash
# Build sonic-express module
cd modules/sonic-express
go build ./cmd/sonic-express

# Run sonic-express commands
./sonic-express version
./sonic-express pack --help
```

## Integration Work

### Creating Integration Briefs

```bash
# Create a new integration brief
cp dev/process/templates/integration-brief-template.md \
   dev/integration/COMPONENT-DESCRIPTION.md

# Edit the brief
$EDITOR dev/integration/COMPONENT-DESCRIPTION.md

# Follow the template structure:
# - Overview
# - Current State
# - Integration Points
# - Open Questions
# - Action Items
# - References
```

### Cross-Component Coordination

```bash
# Check all integration briefs
grep -r "TODO\|TBD\|WIP" dev/integration/

# Review open questions in briefs
cat dev/integration/*/ | grep -A5 "Open Questions"

# Update integration status
$EDITOR dev/integration/VENTOY-INTEGRATION-BRIEF.md
```

## Process Documentation

### Creating Checklists

```bash
# Create a new checklist
cp dev/process/templates/checklist-template.md \
   dev/process/checklists/component-validation.md

# Edit the checklist
$EDITOR dev/process/checklists/component-validation.md

# Use clear validation steps with commands
```

### Creating Workflows

```bash
# Create a new workflow document
cp dev/process/templates/workflow-template.md \
   dev/process/workflows/component-workflow.md

# Edit the workflow
$EDITOR dev/process/workflows/component-workflow.md

# Define clear phases and steps
```

## Composting Old Work

### When to Compost

- Work is completed and documented
- Material is obsolete or superseded
- Keeping it in active lanes creates noise

### How to Compost

```bash
# Move completed requests
mv dev/active/some-request.md dev/compost/requests/

# Move obsolete notes
mv dev/notes/old-note.md dev/compost/notes/

# Add compost header
$EDITOR dev/compost/requests/some-request.md

# Example compost header:
# ---
# composted: 2026-04-19
# reason: completed in vA1.1.0 release
# superseded-by: docs/new-documentation.md
# ---
```

## Best Practices

### Active Work

1. **Keep it small**: No more than 10 active requests
2. **Update daily**: Execution notes reflect real progress
3. **Be explicit**: Document blockers with `[!]` marker
4. **Align with roadmap**: Active index mirrors milestones

### Integration

1. **Clear ownership**: Assign owners to integration points
2. **Document questions**: List open questions explicitly
3. **Track action items**: Use checkboxes for clear progress
4. **Reference upstream**: Link to authoritative sources

### Process

1. **Automate validation**: Prefer scripts over manual checks
2. **Document everything**: All steps should be repeatable
3. **Use templates**: Consistent structure across documents
4. **Review regularly**: Update processes based on lessons learned

### Composting

1. **Be aggressive**: Move completed work promptly
2. **Add context**: Include compost date and reason
3. **Organize**: Keep compost structured by type
4. **Don't hoard**: If it's not useful, consider deleting

## Common Commands

### Development

```bash
# Build all modules
make build

# Run tests
make test

# Check git status
git status

# View recent commits
git log --oneline -10
```

### Ventoy Specific

```bash
# Ventoy submodule status
git submodule status Ventoy

# Update Ventoy submodule
git submodule update --remote Ventoy

# Ventoy build with verbose output
bash -x ./modules/ventoy/build.sh
```

### Documentation

```bash
# Search for TODO items
grep -r "TODO" dev/

# Find all checklists
find dev/process/checklists -name "*.md"

# List all integration briefs
ls -la dev/integration/
```

## Troubleshooting

### Build Issues

```bash
# Check build dependencies
./modules/ventoy/build.sh --check-deps

# Verbose build output
bash -x ./modules/ventoy/build.sh

# Check environment variables
env | grep VENTOY
```

### Git Issues

```bash
# Fix submodule issues
git submodule sync
git submodule update --init --recursive

# Check submodule configuration
cat .gitmodules

# Update submodule to latest
git submodule update --remote Ventoy
```

### Validation Issues

```bash
# Run specific validation steps
./modules/ventoy/validate.sh --build
./modules/ventoy/validate.sh --functional

# Check validation checklist
cat dev/process/checklists/ventoy-validation.md

# Debug validation failures
bash -x ./modules/ventoy/validate.sh
```

## Advanced Topics

### Customizing Workflow

```bash
# Create custom workflow
cp dev/process/templates/workflow-template.md \
   dev/process/workflows/my-workflow.md

# Edit workflow phases
$EDITOR dev/process/workflows/my-workflow.md

# Add to active execution
echo "- [ ] Follow my-workflow.md" >> dev/active/execution-notes.md
```

### Cross-Component Testing

```bash
# Test Ventoy + Sonic Home integration
./modules/ventoy/validate.sh
cd modules/sonic-home && go test ./...

# Test full promotion flow
./modules/ventoy/build.sh
./modules/ventoy/validate.sh
cp -v build/ventoy/* release/ventoy/
```

### Release Preparation

```bash
# Check all validation checklists
find dev/process/checklists -name "*.md" -exec grep -l "\- \[ \]" {} \;

# Review integration briefs
cat dev/integration/*/

# Update completed summary
$EDITOR dev/active/completed-summary.md
```

## Resources

### Documentation

- **Complete Compendium**: `dev/COMPLETE-DEV-FLOW-COMPENDIUM.md` — Master document with all briefs
- **Dev Flow**: `dev/README.md`
- **Roadmap**: `docs/roadmap.md`
- **Promotion**: `docs/promotion.md`
- **Ventoy Integration**: `dev/integration/VENTOY-INTEGRATION-BRIEF.md`

### Templates

- **Integration Brief**: `dev/process/templates/integration-brief-template.md`
- **Checklist**: `dev/process/templates/checklist-template.md`
- **Workflow**: `dev/process/templates/workflow-template.md`

### Tools

- **Build**: `modules/ventoy/build.sh`
- **Validate**: `modules/ventoy/validate.sh`
- **Promote**: (future) `modules/ventoy/promote.sh`

## Next Steps

1. **Explore the structure**: `ls -la dev/`
2. **Review active work**: `cat dev/active/active-index.md`
3. **Check Ventoy integration**: `cat dev/integration/VENTOY-INTEGRATION-BRIEF.md`
4. **Run a build**: `./modules/ventoy/build.sh`
5. **Update execution notes**: Start tracking your work!