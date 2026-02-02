# Getting Started with commit-gen

## Installation & Setup (One-Time)

### Step 1: Start OpenCode Server
Open a terminal and run:
```bash
opencode serve
```

This starts the OpenCode server on localhost:4096. **Keep this running in the background.**

### Step 2: Go to Your Git Repository
```bash
cd /path/to/your/git/repo
```

### Step 3: Install the Git Hook
```bash
/home/avgt/all/kanban/commit-gen/commit-gen install
```

You should see:
```
✓ Git hook installed successfully
Now you can use: git commit -m ""
```

Done! Now you're ready to use commit-gen.

---

## Using commit-gen

### Make Your Changes

```bash
# Edit files, add new features, fix bugs, etc.
# ... your work ...

# Stage your changes
git add .
```

### Generate Commit Message Automatically

```bash
git commit -m ""
```

That's it! The tool will:
1. Detect the empty message
2. Get your staged changes
3. Send them to OpenCode AI
4. Generate a descriptive commit message
5. Fill it in for you

The commit will complete with an AI-generated message!

---

## Advanced Usage

### Manual Generation (Without Committing)

```bash
/home/avgt/all/kanban/commit-gen/commit-gen generate
```

Output:
```
✓ Commit message generated:
  feat(api): add user authentication endpoints
```

### Preview Changes & Generated Message

```bash
/home/avgt/all/kanban/commit-gen/commit-gen preview
```

Shows:
1. The staged changes (git diff)
2. The generated commit message

### Use Different Commit Styles

**Conventional (default):**
```bash
/home/avgt/all/kanban/commit-gen/commit-gen generate --style conventional
# Output: feat(auth): add user authentication
```

**Imperative:**
```bash
/home/avgt/all/kanban/commit-gen/commit-gen generate --style imperative
# Output: Add user authentication to login page
```

**Detailed:**
```bash
/home/avgt/all/kanban/commit-gen/commit-gen generate --style detailed
# Output: feat(auth): add user authentication with OAuth2 support
```

### Check Configuration

```bash
/home/avgt/all/kanban/commit-gen/commit-gen config
```

Shows current settings (OpenCode host/port, style, cache settings, etc.)

### Manage Cache

Check cache status:
```bash
/home/avgt/all/kanban/commit-gen/commit-gen cache status
```

Clear cache (removes all cached sessions):
```bash
/home/avgt/all/kanban/commit-gen/commit-gen cache clear
```

---

## Adding to PATH (Optional)

To use `commit-gen` from anywhere without the full path:

### Bash
Add to `~/.bashrc`:
```bash
export PATH="/home/avgt/all/kanban/commit-gen:$PATH"
```

Then reload:
```bash
source ~/.bashrc
```

### Zsh
Add to `~/.zshrc`:
```bash
export PATH="/home/avgt/all/kanban/commit-gen:$PATH"
```

Then reload:
```bash
source ~/.zshrc
```

After this, you can just use:
```bash
commit-gen install
commit-gen generate
# etc.
```

---

## Complete Workflow Example

```bash
# 1. Start OpenCode in background (first time only)
opencode serve &

# 2. Go to your project
cd ~/projects/my-app

# 3. Install hook (first time only)
commit-gen install

# 4. Make changes
vim src/auth.ts
vim src/api.ts

# 5. Stage changes
git add src/

# 6. Commit with empty message - AI will generate it!
git commit -m ""

# 7. Done! Your commit now has a descriptive message
git log -1
# feat(auth): implement OAuth2 authentication
```

---

## Troubleshooting

### "opencode server is not running at localhost:4096"

**Solution:** Make sure you have `opencode serve` running in another terminal.

```bash
# In a separate terminal:
opencode serve
```

The server should output something like:
```
OpenCode server running at http://localhost:4096
```

### "no staged changes found"

**Solution:** You need to stage your changes first.

```bash
git add .
# or
git add path/to/specific/file
```

### The hook isn't working

**Solution:** Reinstall the hook.

```bash
commit-gen uninstall
commit-gen install
```

Verify it worked:
```bash
cat .git/hooks/prepare-commit-msg
```

Should contain "commit-gen" in the script.

### I want to disable the hook temporarily

**Solution:** Just run git commit normally (not with empty message):

```bash
git commit -m "My manual message"
# This bypasses the hook entirely
```

To permanently disable:
```bash
commit-gen uninstall
```

---

## Configuration (Optional)

Create a file at `~/.config/commit-gen/config.yaml` to customize:

```yaml
opencode:
  host: localhost
  port: 4096
  timeout: 30

generation:
  style: conventional  # Options: conventional, imperative, detailed
  model:
    provider: anthropic
    model_id: claude-3-5-sonnet-20241022

cache:
  enabled: true
  ttl: 24h
```

Or use environment variables:

```bash
export COMMIT_GEN_GENERATION_STYLE=imperative
export COMMIT_GEN_OPENCODE_HOST=localhost
```

---

## Supported Commit Styles

### Conventional Commits
Format: `type(scope): description`
- Types: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `chore`
- Best for: Projects following conventional commit standard
- Example: `feat(auth): add login form validation`

### Imperative
Format: Direct command/action
- Start with a verb (Add, Fix, Update, Remove, etc.)
- Written as if giving a command
- Example: `Add user input validation to login form`

### Detailed
Format: `type(scope): description` with more detail
- Includes scope and detailed description
- Best for: Documenting changes thoroughly
- Example: `feat(auth): add login form validation with email verification`

---

## Tips & Tricks

### Tip 1: Review Before Committing
Use preview to see what will be generated:
```bash
commit-gen preview
```

Then decide if you want to proceed with `git commit -m ""`.

### Tip 2: Different Styles for Different Projects
Set environment variable before running git commit:
```bash
COMMIT_GEN_GENERATION_STYLE=imperative git commit -m ""
```

### Tip 3: Speed Up Commits
Session caching means first commit takes ~2-3 seconds, subsequent commits in same repo are faster (1-2 seconds).

### Tip 4: Dry Run
Preview what will be generated without writing:
```bash
commit-gen generate --dry-run
```

---

## Uninstalling

To remove commit-gen from your repository:

```bash
commit-gen uninstall
```

This removes the git hook. The binary still exists if you want to use it manually later.

To completely remove:
```bash
rm /home/avgt/all/kanban/commit-gen/commit-gen
# Or if in PATH:
which commit-gen && rm $(which commit-gen)
```

---

## Next Steps

1. **Test it out** on a repository with some staged changes
2. **Customize** the style if the default doesn't match your preference
3. **Share** with your team if they use the same commit conventions
4. **Provide feedback** if you find any issues

Enjoy faster, better commit messages! 🚀
